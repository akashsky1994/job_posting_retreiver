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
	"path/filepath"
	"time"

	"job_posting_retreiver/dal"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	log "github.com/sirupsen/logrus"
)

type TrueupHandler struct {
	repo      *repository.AlgoliaConfig
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
}

func NewTrueupHandler(config *config.Config) *TrueupHandler {
	last_date := time.Now().AddDate(0, 0, -1)
	timestamp := last_date.Unix()
	params := append(constant.TRUEUP_QUERY_PARAMS, opt.NumericFilter("updated_at_timestamp>="+fmt.Sprint(timestamp)))
	return &TrueupHandler{
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

func (handler *TrueupHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
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
			err = handler.dao.SaveFile(file_path, "Trueup")
			if err != nil {
				return err
			}
			total_pages = results.NbPages
			// handler.config.Logger.Info(results.NbHits, results.NbPages, results.HitsPerPage, results.Page)

			currPage += 1
		}
	}
	return nil
}

func (handler *TrueupHandler) ProcessJobs() error {
	dir_path := filepath.Join(handler.data_path, "raw_files")
	files, err := ioutil.ReadDir(dir_path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		file_path := filepath.Join(dir_path, file.Name())
		handler.config.Logger.Info("Reading file:", file_path)

		content, err := ioutil.ReadFile(file_path)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		var records []model.TrueUpRecord
		var joblistings []model.JobListing
		err = json.Unmarshal(content, &records)
		if err != nil {
			return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from algolia response: trueup", log.ErrorLevel)
		}

		for _, job := range records {
			db_company, err := handler.dao.GetCompany(job.Company)
			if err != nil {
				handler.config.Logger.Warn(err)
			}
			remote := false
			if job.Remote == 1 {
				remote = true
			}
			joblistings = append(joblistings, model.JobListing{
				JobLink:  utils.CleanURL(job.JobLink),
				JobTitle: job.JobTitle,
				OrgName:  job.Company,
				Location: []string{job.Location},
				Remote:   remote,
				Company:  db_company,
				Source:   "trueup",
			})
		}
		err = handler.dao.SaveJobs(joblistings)
		if err != nil {
			return err
		}
	}
	return nil
}
