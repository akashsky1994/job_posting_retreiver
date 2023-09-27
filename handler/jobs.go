package handler

import (
	"job_posting_retreiver/config"
	"job_posting_retreiver/dal"
	"job_posting_retreiver/errors"
	"net/http"
	"strconv"
)

type JobSourceHandler interface {
	FetchJobs() error
	ProcessJobs() error
	AggregateJobs(res http.ResponseWriter, req *http.Request)
}

func NewJobSourceHandler(source string, config *config.Config) JobSourceHandler {
	switch source {
	case "builtin":
		return NewBuiltInHandler(config)
	case "simplify":
		return NewSimplifyHandler(config)
	case "trueup":
		return NewTrueupHandler(config)
	default:
		config.Logger.Panicln("Job Source Not Implemented")
	}
	return nil
}

type JobHandler struct {
	dao    dal.DataAccessObject
	config *config.Config
}

func NewJobHandler(config *config.Config) *JobHandler {
	return &JobHandler{
		dao:    *dal.NewDataAccessService(config.DB),
		config: config,
	}
}

func (handler *JobHandler) ListJobs(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	var pagination dal.Query
	if val, err := strconv.Atoi(query.Get("page")); err == nil {
		pagination.Page = val
	}
	if val, err := strconv.Atoi(query.Get("per_page")); err == nil {
		pagination.PerPage = val
	}
	if val, err := strconv.Atoi(query.Get("user_id")); err == nil {
		pagination.UserID = val
	}
	jobs, err := handler.dao.ListJobs(pagination)
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		handler.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	RespondwithJSON(res, http.StatusOK, jobs)
}

func (handler *JobHandler) AddJobs(res http.ResponseWriter, req *http.Request) {
	message := map[string]string{"message": "Not Implemented"}
	RespondwithJSON(res, http.StatusNotImplemented, message)
}
