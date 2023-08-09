package dal

import (
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataAccessObject struct {
	conn *gorm.DB
}

func NewDataAccessService(conn *gorm.DB) *DataAccessObject {
	return &DataAccessObject{conn: conn}
}

func (dao *DataAccessObject) SaveJobs(joblistings []model.JobListing) error {
	// Saving to DB
	err := dao.conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job_link"}},
		DoUpdates: clause.AssignmentColumns([]string{"job_title", "location", "company_id", "updated_at", "source"}),
	}).Create(&joblistings).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding jobs to DB", log.ErrorLevel)
	}
	return nil
}

func (dao *DataAccessObject) ListJobs(limit int) ([]model.JobListing, error) {
	var joblistings []model.JobListing
	var err error
	if limit != 0 {
		err = dao.conn.Preload("Company").Limit(limit).Find(&joblistings).Error
	} else {
		err = dao.conn.Preload("Company").Find(&joblistings).Error
	}

	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Something went wrong while fetch data from db", log.ErrorLevel)
	}
	return joblistings, err

}
