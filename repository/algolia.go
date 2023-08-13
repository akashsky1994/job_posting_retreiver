package repository

import (
	"job_posting_retreiver/errors"

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
	FacetFilters [][]string
}

// index string, api_key string, app_id string, params []interface{}
func NewAlgoliaService(source string, index string, api_key string, app_id string, params []interface{}, facetfilters [][]string) *AlgoliaConfig {
	return &AlgoliaConfig{
		Source:       source,
		Index:        index,
		API_KEY:      api_key,
		APP_ID:       app_id,
		Params:       params,
		FacetFilters: facetfilters,
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
