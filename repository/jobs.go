package repository

import (
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type JobBoardService struct {
	Conn *gorm.DB
}

func NewJobBoardService(conn *gorm.DB) *JobBoardService {
	return &JobBoardService{Conn: conn}
}

func (jbservice *JobBoardService) ListJobs(limit int) ([]model.JobListing, error) {
	var joblistings []model.JobListing
	var err error
	if limit != 0 {
		err = jbservice.Conn.Preload("Company").Limit(limit).Find(&joblistings).Error
	} else {
		err = jbservice.Conn.Preload("Company").Find(&joblistings).Error
	}

	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Something went wrong while fetch data from db", log.ErrorLevel)
	}
	return joblistings, err

}
