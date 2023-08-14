package handler

import (
	"encoding/json"
	"io/ioutil"
	"job_posting_retreiver/config"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/dal"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"job_posting_retreiver/utils"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type BuiltInHandler struct {
	repo      repository.BuiltInService
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
}

func NewBuiltInHandler(config *config.Config) *BuiltInHandler {
	var builtin *model.BuiltInRecord
	return &BuiltInHandler{
		repo:      *repository.NewBuiltInService(builtin),
		dao:       *dal.NewDataAccessService(config.DB),
		config:    config,
		data_path: constant.BUILTIN_DATA_PATH,
	}
}

func (handler *BuiltInHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	category_id := chi.URLParam(req, "category_id")
	err := handler.FetchJobs(category_id)
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

func (handler *BuiltInHandler) FetchJobs(category_id string) error {
	total_pages := 1
	currPage := 0
	for currPage != total_pages {
		response, err := handler.repo.RequestJobs(currPage, category_id)
		if err != nil {
			return err
		}

		var records model.BuiltInRecord
		err = json.Unmarshal(response, &records)
		total_pages = records.PageCount

		file_path, err := utils.WriteRawDataToJSONFile(response, handler.data_path)
		if err != nil {
			return err
		}
		err = handler.dao.SaveFile(file_path, "BuiltIn")
		if err != nil {
			return err
		}
		currPage += 1
	}
	return nil
}

func (handler *BuiltInHandler) ProcessJobs() error {
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
		var records model.BuiltInRecord
		var joblistings []model.JobListing
		err = json.Unmarshal(content, &records)
		if err != nil {
			return err
		}
		for _, company := range records.Companies {
			db_company, err := handler.dao.GetCompany(company.Company.Title)
			if err != nil {
				handler.config.Logger.Warn(err)
			}
			for _, job := range company.Jobs {
				remote := true
				if job.Remote != "REMOTE_ENABLED" {
					remote = false
				}
				// Create Job Object
				joblisting := model.JobListing{
					JobLink:  utils.CleanURL(job.JobLink),
					JobTitle: job.JobTitle,
					OrgName:  company.Company.Title,
					Location: []string{job.Location},
					Remote:   remote,
					Company: model.Company{
						ID:   db_company.ID,
						Name: company.Company.Title,
					},
					Source: "builtin",
				}
				joblistings = append(joblistings, joblisting)
			}
		}
		err = handler.dao.SaveJobs(joblistings)
		if err != nil {
			return err
		}
	}
	return nil
}
