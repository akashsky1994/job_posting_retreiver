package repository

import (
	"fmt"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/model"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	log "github.com/sirupsen/logrus"
)

type AlgoliaConfig struct {
	Source       string
	Index        string
	API_KEY      string
	APP_ID       string
	Params       []interface{}
	Record       interface{}
	FacetFilters [][]string
}

// index string, api_key string, app_id string, params []interface{}
func NewSimplifyService(record *model.SimplifyRecord) *AlgoliaConfig {
	return &AlgoliaConfig{
		Source:       "Simplify",
		Index:        constant.SIMPLIFY_INDEX,
		API_KEY:      constant.ALGOLIA_SIMPLIFY_API_KEY,
		APP_ID:       constant.ALGOLIA_SIMPLIFY_APP_ID,
		Params:       constant.SIMPLIFY_QUERY_PARAMS,
		Record:       record,
		FacetFilters: constant.SIMPLIFY_FACET_FILTERS,
	}
}

// index string, api_key string, app_id string, params []interface{}
func NewTrueUpService(record *model.TrueUpRecord) *AlgoliaConfig {
	last_date := time.Now().AddDate(0, 0, -1)
	timestamp := last_date.Unix()
	params := append(constant.TRUEUP_QUERY_PARAMS, opt.NumericFilter("updated_at_timestamp>="+fmt.Sprint(timestamp)))
	return &AlgoliaConfig{
		Source:       "Trueup",
		Index:        constant.TRUEUP_INDEX,
		API_KEY:      constant.ALGOLIA_TRUEUP_API_KEY,
		APP_ID:       constant.ALGOLIA_TRUEUP_APP_ID,
		Params:       params,
		Record:       record,
		FacetFilters: constant.TRUEUP_FACET_FILTERS,
	}
}

func (jbservice *AlgoliaConfig) RequestJobs(page int, additional_params []interface{}) (search.QueryRes, error) {
	client := search.NewClient(jbservice.APP_ID, jbservice.API_KEY)
	index := client.InitIndex(jbservice.Index)

	params := append(jbservice.Params, additional_params...)
	params = append(params, opt.Page(page))
	results, err := index.Search("", params...)
	if err != nil {
		return results, errors.Unexpected.Wrap(err, "Error Loading data from Algolia "+jbservice.Source+" Index", log.ErrorLevel)
	}

	return results, nil
}

func (jbservice *AlgoliaConfig) GetFilterParams() ([]interface{}, error) {
	var FilterParams []interface{}
	all_possible_combinations := dfs(jbservice.FacetFilters, 0, []string{})
	for _, filter := range all_possible_combinations {
		FilterParams = append(FilterParams, opt.FacetFilterAnd(filter))
	}
	return FilterParams, nil
}

func dfs(filters [][]string, idx int, combinations []string) [][]string {
	var res [][]string
	for i := 0; i < len(filters[idx]); i++ {
		new_combinations := append(combinations, filters[idx][i])

		if idx == len(filters)-1 {
			res = append(res, new_combinations)
		} else {
			res = append(res, dfs(filters, idx+1, new_combinations)...)
		}
	}
	return res
}
