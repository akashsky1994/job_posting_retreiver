package dal

import (
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"

	log "github.com/sirupsen/logrus"
)

func (dao *DataAccessObject) SaveFile(file_path string, source string) error {
	filelog := model.FileLogs{
		Source:   source,
		FilePath: file_path,
		Status:   "Pending",
	}
	err := dao.conn.Create(&filelog).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding Files to DB", log.ErrorLevel)
	}
	return nil
}
