package utils

import (
	"job_posting_retreiver/errors"
	"os"
	"time"

	"path/filepath"

	"github.com/sirupsen/logrus"
)

func WriteRawDataToJSONFile(content []byte, parent_path string) (string, error) {
	t := time.Now()
	ts := t.Format("20060102150405")
	file_path := filepath.Join(parent_path, "raw_files", ts+".json")
	if err := os.MkdirAll(filepath.Dir(file_path), 0777); err != nil {
		return "", errors.Unexpected.Wrap(err, "Error Creating directory", logrus.ErrorLevel)
	}
	if err := os.WriteFile(file_path, content, 0777); err != nil {
		return "", errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	}
	return file_path, nil
}
