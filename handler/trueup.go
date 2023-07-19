package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"net/http"

	"job_posting_retreiver/dal"

	log "github.com/sirupsen/logrus"
)

type TrueupHandler struct {
	*AlgoliaHandler
}

func NewTrueupHandler(config *config.Config) *TrueupHandler {
	var record *model.TrueUpRecord
	return &TrueupHandler{
		&AlgoliaHandler{
			repo:   *repository.NewTrueUpService(record),
			dao:    *dal.NewDataAccessService(config.DB),
			config: config,
		},
	}
}

func (handler *TrueupHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	err := handler.FetchJobs()
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
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
			var records []model.TrueUpRecord
			var joblistings []model.JobListing
			results, err := handler.repo.RequestJobs(currPage, []interface{}{param})
			if err != nil {
				return err
			}
			total_pages = results.NbPages
			handler.config.Logger.Info(results.NbHits, results.NbPages, results.HitsPerPage, results.Page)
			err = results.UnmarshalHits(&records)
			if err != nil {
				return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from algolia response: simplify", log.ErrorLevel)
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
					JobLink:  job.JobLink,
					JobTitle: job.JobTitle,
					OrgName:  job.Company,
					Location: job.Location,
					Remote:   remote,
					Company:  db_company,
					Source:   "trueup",
				})
			}
			err = handler.dao.SaveJobs(joblistings)
			if err != nil {
				return err
			}
			currPage += 1
		}
	}
	return nil
}
