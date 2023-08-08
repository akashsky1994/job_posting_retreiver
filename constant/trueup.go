package constant

import "github.com/algolia/algoliasearch-client-go/v3/algolia/opt"

var (
	TRUEUP_INDEX           string = "job_search"
	ALGOLIA_TRUEUP_API_KEY string = "dcbd5d6e6bd6f841353e83b5b9557dee"
	ALGOLIA_TRUEUP_APP_ID  string = "3PR3Y01TFY"
	TRUEUP_QUERY_PARAMS           = []interface{}{
		opt.HitsPerPage(100),
	}
	TRUEUP_FACET_FILTERS_MAP = map[string][]string{
		"level":               {"2", "3", "4"},
		"job_categories_lvl0": {"Engineering (Software)", "Data & Analytics", "Data Science & Machine Learning", "Research"},
	}
	TRUEUP_FACET_FILTERS = [][]string{
		{"job_categories_lvl1:Engineering (Software) > Generalist", "job_categories_lvl1:Engineering (Software) > Data", "job_categories_lvl1:Engineering (Software) > QA", "job_categories_lvl1:Engineering (Software) > Machine Learning", "job_categories_lvl1:Engineering (Software) > Security", "job_categories_lvl1:Engineering (Software) > Site Reliability", "job_categories_lvl1:Engineering (Software) > DevOps", "job_categories_lvl1:Engineering (Software) > Networking", "job_categories_lvl0:Engineering (Software)", "job_categories_lvl0:Data & Analytics", "job_categories_lvl0:Data Science & Machine Learning", "job_categories_lvl0:Research"},
		{"level:2", "level:3", "level:4"},
	}
)

var TRUEUP_DATA_PATH = "./data/trueup"
