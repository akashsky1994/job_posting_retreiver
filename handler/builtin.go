package handler

import (
	"encoding/json"
	"job_posting_retreiver/config"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"job_posting_retreiver/repository"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type BuiltInHandler struct {
	repo   repository.BuiltInService
	config *config.Config
}

func NewBuiltInHandler(config *config.Config) *BuiltInHandler {
	var builtin *model.BuiltInOutput
	return &BuiltInHandler{
		repo:   *repository.NewBuiltInService(builtin),
		config: config,
	}
}

func (bh *BuiltInHandler) FetchJobsHandler(res http.ResponseWriter, req *http.Request) {
	category_id := chi.URLParam(req, "category_id")
	err := bh.FetchJobs(category_id)
	if err != nil {
		errType, severity := errors.GetTypeAndLogLevel(err)
		bh.config.Logger.Log(severity, err)
		HandleError(res, err, errType)
	}
	message := map[string]string{"message": "Fetching Successful"}
	RespondwithJSON(res, http.StatusOK, message)
	// http.Redirect(res, req, filepath, http.StatusOK)
	// res.WriteHeader(http.StatusOK)
	// res.Header().Set("Content-Type", "application/octet-stream")
	// res.Write(fileBytes)
	// return
}

func (bh *BuiltInHandler) FetchJobs(category_id string) error {
	err := bh.repo.RequestJobs(1, category_id)
	if err != nil {
		return err
	}

	total_pages := bh.repo.JBBuiltIn.PageCount
	var joblistings []model.JobListing
	for page := 1; page <= total_pages; page++ {
		err := bh.repo.RequestJobs(page, category_id)
		if err != nil {
			return err
		}

		for _, company := range bh.repo.JBBuiltIn.Companies {
			for _, job := range company.Jobs {

				// Get Company ID if exists in DB else create new entry
				db_company := model.Company{
					Name: company.Company.Title,
				}
				err := bh.config.DB.FirstOrCreate(&db_company, model.Company{Name: company.Company.Title}).Error
				if err != nil {
					bh.config.Logger.Warn(err)
				}

				// Create Job Object
				joblisting := model.JobListing{
					JobLink:  job.JobLink,
					JobTitle: job.JobTitle,
					OrgName:  company.Company.Title,
					Location: job.Location,
					Remote:   job.Remote,
					Company: model.Company{
						ID:   db_company.ID,
						Name: company.Company.Title,
					},
					Source: "builtin",
				}
				joblistings = append(joblistings, joblisting)
			}
		}
	}

	// Saving to DB
	err = bh.config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job_link"}},
		DoUpdates: clause.AssignmentColumns([]string{"job_title", "location", "company_id", "updated_at", "source"}),
	}).Create(&joblistings).Error
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding jobs to DB", logrus.ErrorLevel)
	}

	// Saving to JSON for history management
	jsonArr, err := json.Marshal(joblistings)
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Loading data into JSON", logrus.ErrorLevel)
	}
	t := time.Now()
	ts := t.Format("20060102150405")
	if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".json", jsonArr, 0666); err != nil {
		return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	}

	// Saving to JSON for legacy maintenance
	csvArr, err := gocsv.MarshalString(&joblistings)
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Converting byte data into csv", logrus.ErrorLevel)
	}
	if err := os.WriteFile("./data/builtinjobs_"+category_id+"_"+ts+".csv", []byte(csvArr), 0666); err != nil {
		return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	}
	if err := os.WriteFile("./data/builtinjobs_"+category_id+".csv", []byte(csvArr), 0666); err != nil {
		return errors.Unexpected.Wrap(err, "Error Writing data into output file", logrus.ErrorLevel)
	}
	return nil
}

func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	// fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func HandleError(w http.ResponseWriter, err error, errType errors.ErrorType) {
	var status int
	switch errType {
	case errors.NotFound:
		status = http.StatusNotFound
	case errors.Unexpected:
		status = http.StatusInternalServerError
	case errors.ExternalAPIError:
		status = http.StatusInternalServerError
	case errors.DataProcessingError:
		status = http.StatusInternalServerError
	}
	RespondwithJSON(w, status, map[string]string{"message": err.Error()})
}
