package handler

import (
	"encoding/json"
	"io/ioutil"
	"job_posting_retreiver/config"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"job_posting_retreiver/utils"
	"log"
	"net/http"
	"path/filepath"

	"job_posting_retreiver/dal"

	"github.com/sirupsen/logrus"
)

type SimplifyHandler struct {
	repo      *repository.TypesenseConfig
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
}

func NewSimplifyHandler(config *config.Config) *SimplifyHandler {
	return &SimplifyHandler{
		repo: repository.NewTypeSenseService(
			"simplify",
			constant.TYPESENSE_SIMPLIFY_COLLECTION,
			constant.TYPESENSE_SIMPLIFY_API_KEY,
			constant.TYPESENSE_SIMPLIFY_URI,
			constant.TYPESENSE_SIMPLIFY_SEARCH_PARAMS,
		),
		dao:       *dal.NewDataAccessService(config.DB),
		config:    config,
		data_path: constant.SIMPLIFY_DATA_PATH,
	}
}

func (handler *SimplifyHandler) JobSourceHandler(res http.ResponseWriter, req *http.Request) {
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

func (handler *SimplifyHandler) FetchJobs() error {
	total_pages := 1
	currPage := 0
	for currPage != total_pages {
		results, err := handler.repo.RequestJobs(currPage)
		if err != nil {
			return err
		}
		total_pages = (*results.Found) / results.RequestParams.PerPage
		payload, err := json.Marshal(results.Hits)
		if err != nil {
			return err
		}
		file_path, err := utils.WriteRawDataToJSONFile(payload, handler.data_path)
		if err != nil {
			return err
		}
		err = handler.dao.SaveFile(file_path, "simplify")
		if err != nil {
			return err
		}
		currPage += 1
	}

	return nil
}

func (handler *SimplifyHandler) ProcessJobs() error {
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
		var records []model.SimplifyRecord
		var joblistings []model.JobListing
		err = json.Unmarshal(content, &records)
		if err != nil {
			return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from algolia response: trueup", logrus.ErrorLevel)
		}
		for _, job := range records {
			db_company, err := handler.dao.GetCompany(job.Company)
			if err != nil {
				handler.config.Logger.Warn(err)
			}
			joblistings = append(joblistings, model.JobListing{
				JobLink:  utils.CleanURL(job.JobLink),
				JobTitle: job.JobTitle,
				OrgName:  job.Company,
				Location: job.Location,
				Remote:   job.Remote,
				Company:  db_company,
				Source:   "simplify",
			})
		}
		err = handler.dao.SaveJobs(joblistings)
		if err != nil {
			return err
		}
	}
	return nil
}
