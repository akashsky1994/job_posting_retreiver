package utils

import (
	"job_posting_retreiver/errors"
	"os"

	"path/filepath"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func WriteRawDataToJSONFile(content []byte, parent_path string) (string, error) {
	file_name := uuid.New().String() + ".json"
	file_path := filepath.Join(parent_path, "raw_files", file_name)
	if err := os.MkdirAll(filepath.Dir(file_path), 0777); err != nil {
		return "", errors.Unexpected.Wrap(err, "Error Creating directory", logrus.ErrorLevel)
	}
	if err := os.WriteFile(file_path, content, 0777); err != nil {
		return "", errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	}
	return file_path, nil
}
