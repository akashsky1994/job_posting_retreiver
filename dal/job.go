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
		DoUpdates: clause.AssignmentColumns([]string{"job_title", "locations", "company_id", "updated_at", "source"}),
	}).Create(&joblistings).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding jobs to DB", log.ErrorLevel)
	}
	return nil
}

func (dao *DataAccessObject) ListJobs(pagination Pagination) (*Pagination, error) {
	var joblistings []*model.JobListing
	err := dao.conn.Model(&model.JobListing{}).Order("id desc").Preload("Company").Scopes(paginate(&joblistings, &pagination, dao.conn)).Select("id", "job_link", "job_title", "source", "locations", "company_id").Find(&joblistings).Error
	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Something went wrong while fetch data from db", log.ErrorLevel)
	}
	pagination.Rows = joblistings
	return &pagination, err
}

func (dao *DataAccessObject) SaveRegions(countries []model.Country) error {
	err := dao.conn.Create(&countries).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding jobs to DB", log.ErrorLevel)
	}
	return nil
}
