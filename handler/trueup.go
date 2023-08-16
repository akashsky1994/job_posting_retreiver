package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"job_posting_retreiver/config"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"job_posting_retreiver/utils"
	"net/http"
	"strings"
	"time"

	"job_posting_retreiver/dal"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type TrueupHandler struct {
	repo      *repository.AlgoliaConfig
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
	name      string
}

func NewTrueupHandler(config *config.Config) *TrueupHandler {
	last_date := time.Now().AddDate(0, 0, -7)
	timestamp := last_date.Unix()
	params := append(constant.TRUEUP_QUERY_PARAMS, opt.NumericFilter("updated_at_timestamp>="+fmt.Sprint(timestamp)))
	return &TrueupHandler{
		name: "trueup",
		repo: repository.NewAlgoliaService(
			"trueup",
			constant.ALGOLIA_TRUEUP_INDEX,
			constant.ALGOLIA_TRUEUP_API_KEY,
			constant.ALGOLIA_TRUEUP_APP_ID,
			params,
			constant.TRUEUP_FACET_FILTERS,
		),
		dao:       *dal.NewDataAccessService(config.DB),
		config:    config,
		data_path: constant.TRUEUP_DATA_PATH,
	}
}

func (handler *TrueupHandler) AggregateJobs(res http.ResponseWriter, req *http.Request) {
	err := handler.FetchJobs()
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
		return
	}
	err = handler.ProcessJobs()
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
		return
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
}

func (handler *TrueupHandler) FetchJobs() error {
	params, err := handler.repo.GetFilterParams()
	if err != nil {
		return err
	}
	for _, param := range params {
		total_pages := 1
		currPage := 0
		for currPage < total_pages {
			results, err := handler.repo.RequestJobs(currPage, []interface{}{param})
			if err != nil {
				return err
			}
			payload, err := json.Marshal(results.Hits)
			if err != nil {
				return err
			}
			file_path, err := utils.WriteRawDataToJSONFile(payload, handler.data_path)
			if err != nil {
				return err
			}
			err = handler.dao.CreateFileLog(file_path, handler.name)
			if err != nil {
				return err
			}
			total_pages = results.NbPages

			currPage += 1
		}
	}
	return nil
}

func (handler *TrueupHandler) ProcessJobs() error {
	var ALLOWED_REGIONS []string
	if data, found := handler.config.Cache.Get("allowed_regions"); found {
		ALLOWED_REGIONS = data.([]string)
	}
	files, err := handler.dao.ListPendingFiles(handler.name)
	if err != nil {
		return err
	}
	for _, file := range files {
		handler.config.Logger.Info("Reading file:", file.FilePath)

		content, err := ioutil.ReadFile(file.FilePath)
		if err != nil {
			return errors.DataProcessingError.Wrap(err, "Error Reading file", logrus.ErrorLevel)
		}
		var records []model.TrueUpRecord
		var joblistings []model.JobListing
		err = json.Unmarshal(content, &records)
		if err != nil {
			return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from algolia response: trueup", log.ErrorLevel)
		}

		for _, job := range records {
			is_allowed := false
			for _, loc := range strings.Split(job.Location, ",") {
				if utils.StringInSlice(strings.TrimSpace(loc), ALLOWED_REGIONS) {
					is_allowed = true
				}
			}
			if is_allowed {
				db_company, err := handler.dao.GetCompany(job.Company)
				if err != nil {
					handler.config.Logger.Warn(err)
					raven.CaptureErrorAndWait(err, nil)
					continue
				}
				remote := false
				if job.Remote == 1 {
					remote = true
				}
				joblistings = append(joblistings, model.JobListing{
					JobLink:   utils.CleanURL(job.JobLink),
					JobTitle:  job.JobTitle,
					OrgName:   job.Company,
					Locations: []string{job.Location},
					Remote:    remote,
					Company:   db_company,
					Source:    handler.name,
					FileLog:   file,
				})
			}
		}
		err = handler.dao.SaveJobs(joblistings)
		if err != nil {
			return err
		}
		file.Status = "Completed"
		err = handler.dao.UpdateFileLog(file)
		if err != nil {
			return err
		}
	}
	return nil
}
