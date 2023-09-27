package handler

import (
	"encoding/json"
	"job_posting_retreiver/config"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/dal"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"job_posting_retreiver/utils"
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

type BuiltInHandler struct {
	repo      repository.BuiltInService
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
	name      string
}

func NewBuiltInHandler(config *config.Config) *BuiltInHandler {
	var builtin *model.BuiltInRecord
	return &BuiltInHandler{
		name:      "builtin",
		repo:      *repository.NewBuiltInService(builtin),
		dao:       *dal.NewDataAccessService(config.DB),
		config:    config,
		data_path: constant.BUILTIN_DATA_PATH,
	}
}

func (handler *BuiltInHandler) AggregateJobs(res http.ResponseWriter, req *http.Request) {
	// category_id := chi.URLParam(req, "category_id")
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

func (handler *BuiltInHandler) FetchJobs() error {
	total_pages := 1
	currPage := 0
	for _, category_id := range constant.CATEGORY_IDS {
		for currPage != total_pages {
			response, err := handler.repo.RequestJobs(currPage, category_id)
			if err != nil {
				return err
			}

			var records model.BuiltInRecord
			err = json.Unmarshal(response, &records)
			if err != nil {
				return errors.DataProcessingError.Wrap(err, "BuiltInHandler Error", logrus.ErrorLevel)
			}
			total_pages = records.PageCount

			file_path, err := utils.WriteRawDataToJSONFile(response, handler.data_path)
			if err != nil {
				return err
			}
			err = handler.dao.CreateFileLog(file_path, handler.name)
			if err != nil {
				return err
			}
			currPage += 1
		}
	}

	return nil
}

func (handler *BuiltInHandler) ProcessJobs() error {
	files, err := handler.dao.ListPendingFiles(handler.name)
	if err != nil {
		return err
	}

	for _, file := range files {
		handler.config.Logger.Info("Reading file:", file.FilePath)

		content, err := os.ReadFile(file.FilePath)
		if err != nil {
			return errors.DataProcessingError.Wrap(err, "Error Reading file", logrus.ErrorLevel)
		}
		var records model.BuiltInRecord
		var joblistings []model.JobListing
		err = json.Unmarshal(content, &records)
		if err != nil {
			return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from builtin response", logrus.ErrorLevel)
		}
		for _, company := range records.Companies {
			db_company, err := handler.dao.GetCompany(company.Company.Title)
			if err != nil {
				handler.config.Logger.Error(err)
				raven.CaptureErrorAndWait(err, nil)
				continue
			}
			for _, job := range company.Jobs {
				remote := true
				if job.Remote != "REMOTE_ENABLED" {
					remote = false
				}
				// Create Job Object
				joblisting := model.JobListing{
					JobLink:   utils.CleanURL(job.JobLink),
					JobTitle:  job.JobTitle,
					OrgName:   company.Company.Title,
					Locations: []string{job.Location},
					Remote:    remote,
					Company:   db_company,
					Source:    handler.name,
					FileLog:   file,
				}
				joblistings = append(joblistings, joblisting)
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
