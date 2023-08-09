package dal

import (
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"

	log "github.com/sirupsen/logrus"
)

func (dao *DataAccessObject) GetCompany(title string) (model.Company, error) {
	// Get Company ID if exists in DB else create new entry
	company := model.Company{
		Name: title,
	}
	err := dao.conn.FirstOrCreate(&company, model.Company{Name: title}).Error
	if err != nil {
		return company, errors.DataProcessingError.Wrap(err, "Error getting Company from DB", log.ErrorLevel)
	}
	return company, nil
}
