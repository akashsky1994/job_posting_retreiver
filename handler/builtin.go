package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/dal"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BuiltInHandler struct {
	repo   repository.BuiltInService
	dao    dal.DataAccessObject
	config *config.Config
}

func NewBuiltInHandler(config *config.Config) *BuiltInHandler {
	var builtin *model.BuiltInOutput
	return &BuiltInHandler{
		repo:   *repository.NewBuiltInService(builtin),
		dao:    *dal.NewDataAccessService(config.DB),
		config: config,
	}
}

func (bh *BuiltInHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	category_id := chi.URLParam(req, "category_id")
	err := bh.FetchJobs(category_id)
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		bh.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
	// http.Redirect(res, req, filepath, http.StatusOK)
	// res.WriteHeader(http.StatusOK)
	// res.Header().Set("Content-Type", "application/octet-stream")
	// res.Write(fileBytes)
	// return
}

func (handler *BuiltInHandler) FetchJobs(category_id string) error {
	err := handler.repo.RequestJobs(1, category_id)
	if err != nil {
		return err
	}

	total_pages := handler.repo.JBBuiltIn.PageCount
	var joblistings []model.JobListing
	for page := 1; page <= total_pages; page++ {
		err := handler.repo.RequestJobs(page, category_id)
		if err != nil {
			return err
		}

		for _, company := range handler.repo.JBBuiltIn.Companies {
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
	}

	err = handler.dao.SaveJobs(joblistings)
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
