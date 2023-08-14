package repository

import (
	"job_posting_retreiver/errors"

	log "github.com/sirupsen/logrus"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

type TypesenseConfig struct {
	Source       string
	Index        string
	API_KEY      string
	URI          string
	Params       *api.SearchCollectionParams
	FacetFilters [][]string
}

// index string, api_key string, app_id string, params []interface{}
func NewTypeSenseService(source string, index string, api_key string, uri string, params *api.SearchCollectionParams) *TypesenseConfig {
	return &TypesenseConfig{
		Source:  source,
		Index:   index,
		API_KEY: api_key,
		URI:     uri,
		Params:  params,
	}
}

func (jbservice *TypesenseConfig) RequestJobs(page int) (*api.SearchResult, error) {
	client := typesense.NewClient(
		typesense.WithServer(jbservice.URI),
		typesense.WithAPIKey(jbservice.API_KEY),
	)
	jbservice.Params.Page = &page
	response, err := client.Collection(jbservice.Index).Documents().Search(jbservice.Params)
	if err != nil {
		return nil, errors.ExternalAPIError.Wrap(err, "Error Loading data from Typesense "+jbservice.Source+" Index", log.ErrorLevel)
	}
	return response, nil
}
