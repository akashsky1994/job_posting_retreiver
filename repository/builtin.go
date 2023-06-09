package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	jrconstant "job_posting_retreiver/constant"
	"job_posting_retreiver/model"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type JobBoardService struct {
	JBBuiltIn *model.BuiltIn
}

func NewJobBoardService(builtin *model.BuiltIn) *JobBoardService {
	return &JobBoardService{JBBuiltIn: builtin}
}

func (jbservice *JobBoardService) RequestJobs(page int, category_id string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", jrconstant.BuiltInURI, nil)
	if err != nil {
		log.Fatal(err)
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
		fmt.Println(err)
	}

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(resbody, &jbservice.JBBuiltIn)
	// fmt.Println(string(resbody))
	// return string(resbody)
	return err
}
