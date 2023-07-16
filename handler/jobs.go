package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/repository"
	"net/http"
	"strconv"
)

type JobHandler struct {
	repo   repository.JobBoardService
	config *config.Config
}

func NewJobHandler(config *config.Config) *JobHandler {
	return &JobHandler{
		repo:   *repository.NewJobBoardService(config.DB),
		config: config,
	}
}

func (jhandler *JobHandler) ListJobs(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	limit := 0
	if val, err := strconv.Atoi(query.Get("limit")); err == nil {
		limit = val
	}
	jobs, err := jhandler.repo.ListJobs(limit)
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		jhandler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	RespondwithJSON(res, http.StatusOK, jobs)
}

func (jh *JobHandler) AddJobs(res http.ResponseWriter, req *http.Request) {
	message := map[string]string{"message": "Not Implemented"}
	RespondwithJSON(res, http.StatusNotImplemented, message)
}
