package repository

import (
	"io/ioutil"
	jrconstant "job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type BuiltInService struct {
	uri     string
	queries map[string]string
	record  *model.BuiltInRecord
}

func NewBuiltInService(builtin *model.BuiltInRecord) *BuiltInService {
	return &BuiltInService{
		uri:     jrconstant.BuiltInURI,
		record:  builtin,
		queries: jrconstant.BuiltInURIQueries,
	}
}

func (jbservice *BuiltInService) RequestJobs(page int, category_id string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", jbservice.uri, nil)
	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Something went wrong while creating new request", log.ErrorLevel)
	}
	q := req.URL.Query()

	for key, val := range jbservice.queries {
		q.Add(key, val)
	}
	q.Set("page", strconv.Itoa(page))
	q.Set("categories", category_id)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.ExternalAPIError.Wrap(err, "Error getting response from built in api", log.ErrorLevel)
	}

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Error Loading data from BuiltIn API", log.ErrorLevel)
	}
	return resbody, nil
}
