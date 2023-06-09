package handler

import (
	"encoding/json"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"math"
	"os"
	"time"

	jrconstant "job_posting_retreiver/constant"

	"github.com/gocarina/gocsv"
)

func FetchBuiltInJobs(category_id string) error {
	var builtin *model.BuiltIn
	jbService := repository.NewJobBoardService(builtin)
	err := jbService.RequestJobs(1, category_id)
	if err != nil {
		return err
	}

	total_pages := int(math.Ceil(float64(jbService.JBBuiltIn.JobCount) / float64(jrconstant.PerPage)))
	var joblistings []model.JobListing
	for page := 1; page <= total_pages; page++ {
		err := jbService.RequestJobs(page, category_id)
		if err != nil {
			return err
		}

		for _, company := range jbService.JBBuiltIn.Companies {
			for _, job := range company.Jobs {
				var joblisting model.JobListing
				joblisting.Company = company.Company.Title
				joblisting.JobLink = job.JobLink
				joblisting.JobTitle = job.JobTitle
				joblisting.Location = job.Location
				joblisting.Remote = job.Remote
				joblistings = append(joblistings, joblisting)
			}
		}
	}
	jsonArr, err := json.Marshal(joblistings)
	if err != nil {
		return err
	}
	t := time.Now()
	ts := t.Format("20060102150405")
	if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".json", jsonArr, 0666); err != nil {
		return err
	}

	csvArr, err := gocsv.MarshalString(&joblistings)
	if err != nil {
		return err
	}
	if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".csv", []byte(csvArr), 0666); err != nil {
		return err
	}
	if err := os.WriteFile("./data/builtinjobs_"+category_id+".csv", []byte(csvArr), 0666); err != nil {
		return err
	}
	return nil
}
