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

func (dao *DataAccessObject) ListJobs(query Query) (*Query, error) {
	var joblistings []*model.JobListing
	err := dao.conn.Model(&model.JobListing{}).Order("id desc").Preload("Company").Scopes(paginate(&joblistings, &query, dao.conn), dao.exclude_user_jobs(query.UserID), dao.exclude_older_jobs()).Select("id", "job_link", "job_title", "source", "locations", "company_id").Find(&joblistings).Error
	if err != nil {
		return nil, errors.Unexpected.Wrap(err, "Something went wrong while fetch data from db", log.ErrorLevel)
	}
	query.Rows = joblistings
	return &query, err
}

func (dao *DataAccessObject) SaveRegions(countries []model.Country) error {
	err := dao.conn.Create(&countries).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding jobs to DB", log.ErrorLevel)
	}
	return nil
}

func (dao *DataAccessObject) exclude_user_jobs(user_id int) func(db *gorm.DB) *gorm.DB {
	subquery := dao.conn.Table("user_jobs").Select("job_id").Where("user_id = ?", user_id)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id not in (?)", subquery)
	}
}

func (dao *DataAccessObject) exclude_older_jobs() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("updated_at > (NOW() - INTERVAL '1' MONTH)")
	}
}
