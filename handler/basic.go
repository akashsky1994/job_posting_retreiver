package handler

import (
	"encoding/json"
	"job_posting_retreiver/errors"
	"net/http"

	"github.com/getsentry/raven-go"
)

func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	// fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func HandleError(w http.ResponseWriter, err error, errType errors.ErrorType) {
	raven.CaptureErrorAndWait(err, map[string]string{"type": "api_call"})
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
