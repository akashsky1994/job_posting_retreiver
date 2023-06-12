package repository

import (
	"encoding/json"
	"io/ioutil"
	jrconstant "job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type JobBoardService struct {
	JBBuiltIn *model.BuiltInOutput
}

func NewJobBoardService(builtin *model.BuiltInOutput) *JobBoardService {
	return &JobBoardService{JBBuiltIn: builtin}
}

func (jbservice *JobBoardService) RequestJobs(page int, category_id string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", jrconstant.BuiltInURI, nil)
	if err != nil {
		return errors.Unexpected.Wrap(err, "Something went wrong while creating new request", log.ErrorLevel)
	}
	q := req.URL.Query()

	for key, val := range jrconstant.BuiltInURIQueries {
		q.Add(key, val)
	}
	q.Set("page", strconv.Itoa(page))
	q.Set("categories", category_id)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return errors.ExternalAPIError.Wrap(err, "Error getting response from built in api", log.ErrorLevel)
	}

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Unexpected.Wrap(err, "Error Loading data from BuiltIn API", log.ErrorLevel)
	}
	json.Unmarshal(resbody, &jbservice.JBBuiltIn)
	return nil
}
