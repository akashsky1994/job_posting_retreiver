package dal

import (
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"

	log "github.com/sirupsen/logrus"
)

func (dao *DataAccessObject) CreateFileLog(file_path string, source string) error {
	filelog := model.FileLog{
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

func (dao *DataAccessObject) ListPendingFiles(source string) ([]model.FileLog, error) {
	filelogquery := model.FileLog{
		Source: source,
		Status: "Pending",
	}
	var files []model.FileLog
	err := dao.conn.Where(&filelogquery).Find(&files).Error
	if err != nil {
		return files, errors.DataProcessingError.Wrap(err, "Error Adding Files to DB", log.ErrorLevel)
	}
	return files, nil
}

func (dao *DataAccessObject) UpdateFileLog(filelog model.FileLog) error {
	err := dao.conn.Save(&filelog).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding Files to DB", log.ErrorLevel)
	}
	return nil
}
