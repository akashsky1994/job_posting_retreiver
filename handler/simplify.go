package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"net/http"

	"job_posting_retreiver/dal"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SimplifyHandler struct {
	*AlgoliaHandler
}

func NewSimplifyHandler(config *config.Config) *SimplifyHandler {
	var record *model.SimplifyRecord
	return &SimplifyHandler{
		&AlgoliaHandler{
			repo:   *repository.NewSimplifyService(record),
			dao:    *dal.NewDataAccessService(config.DB),
			config: config,
		},
	}

}

func (handler *SimplifyHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	err := handler.FetchJobs()
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
}

func (handler *SimplifyHandler) FetchJobs() error {
	params, err := handler.repo.GetFilterParams()
	if err != nil {
		return err
	}
	for _, param := range params {
		total_pages := 1
		currPage := 0
		for currPage != total_pages {
			var records []model.SimplifyRecord
			var joblistings []model.JobListing
			results, err := handler.repo.RequestJobs(currPage, []interface{}{param})
			if err != nil {
				return err
			}

			handler.config.Logger.Info(results.NbHits, results.NbPages, results.HitsPerPage)
			total_pages = results.NbPages
			err = results.UnmarshalHits(&records)
			if err != nil {
				return errors.Unexpected.Wrap(err, "Error Unmarshaling hits from algolia response: simplify", log.ErrorLevel)
			}

			for _, job := range records {
				db_company, err := handler.dao.GetCompany(job.Company)
				if err != nil {
					handler.config.Logger.Warn(err)
				}
				joblistings = append(joblistings, model.JobListing{
					JobLink:  job.JobLink,
					JobTitle: job.JobTitle,
					OrgName:  job.Company,
					Location: strings.Join(job.Location, "||"),
					Remote:   job.Remote,
					Company:  db_company,
					Source:   "simplify",
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
