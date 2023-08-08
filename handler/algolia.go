package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/dal"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/repository"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type AlgoliaHandler struct {
	repo      repository.AlgoliaConfig
	dao       dal.DataAccessObject
	config    *config.Config
	data_path string
}

func (handler *AlgoliaHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	err := handler.FetchJobs()
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
}

func (handler *AlgoliaHandler) FetchJobs() error {
	return errors.NotImplementedError.New("Can not call basic handler directly", log.ErrorLevel)
}
