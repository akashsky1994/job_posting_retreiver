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

func (bh *BuiltInHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	category_id := chi.URLParam(req, "category_id")
	err := bh.FetchJobs(category_id)
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		bh.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
		return
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
}

func (handler *BuiltInHandler) FetchJobs(category_id string) error {
	total_pages := 1
	currPage := 0
	var joblistings []model.JobListing
	for currPage != total_pages {
		var records model.BuiltInRecord
		response, err := handler.repo.RequestJobs(currPage, category_id)
		if err != nil {
			return err
		}
		file_path, err := utils.WriteRawDataToJSONFile(response, handler.data_path)
		if err != nil {
			return err
		}
		err = handler.dao.SaveFile(file_path, "BuiltIn")
		if err != nil {
			return err
		}
		err = json.Unmarshal(response, &records)
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
					JobLink:  job.JobLink,
					JobTitle: job.JobTitle,
					OrgName:  company.Company.Title,
					Location: job.Location,
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
		currPage += 1
	}

	err := handler.dao.SaveJobs(joblistings)
	if err != nil {
		return err
	}

	// // Saving to JSON for history management
	// jsonArr, err := json.Marshal(joblistings)
	// if err != nil {
	// 	return errors.DataProcessingError.Wrap(err, "Error Loading data into JSON", logrus.ErrorLevel)
	// }
	// t := time.Now()
	// ts := t.Format("20060102150405")
	// if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".json", jsonArr, 0666); err != nil {
	// 	return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	// }

	// // Saving to JSON for legacy maintenance
	// csvArr, err := gocsv.MarshalString(&joblistings)
	// if err != nil {
	// 	return errors.DataProcessingError.Wrap(err, "Error Converting byte data into csv", logrus.ErrorLevel)
	// }
	// if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".csv", []byte(csvArr), 0666); err != nil {
	// 	return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	// }
	// if err := os.WriteFile("./data/builtinjobs_"+category_id+".csv", []byte(csvArr), 0666); err != nil {
	// 	return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	// }
	return nil
}
