package constant

import "github.com/algolia/algoliasearch-client-go/v3/algolia/opt"

var (
	SIMPLIFY_INDEX           string = "JOBS_updated_date_desc"
	ALGOLIA_SIMPLIFY_API_KEY string = "068125d0565c0f7230bd7becf65c46f1"
	ALGOLIA_SIMPLIFY_APP_ID  string = "4N95P1L3C8"
	SIMPLIFY_QUERY_PARAMS           = []interface{}{
		opt.HitsPerPage(100),
		opt.InsideBoundingBox([][4]float64{
			{
				71.5388001, -66.885417, 18.7763, -180,
			},
			{
				71.5388001, 180, 18.7763, 170.5957,
			},
		}),
	}
	SIMPLIFY_FACET_FILTERS_MAP = map[string][]string{
		"experience_level": {"Entry Level/New Grad", "Junior", "Mid Level", "Senior"},
		"functions":        {"Software Engineering", "DevOps & Infrastructure", "AI & Machine Learning", "Lab & Research", "Data & Analytics"},
	}
	SIMPLIFY_FACET_FILTERS = [][]string{
		{"experience_level:Entry Level/New Grad", "experience_level:Junior", "experience_level:Mid Level", "experience_level:Senior"},
		{"functions:Software Engineering", "functions:DevOps & Infrastructure", "functions:AI & Machine Learning", "functions:Data & Analytics", "functions:Lab & Research"},
	}
)
