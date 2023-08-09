package constant

import "strconv"

var BuiltInURI string = "https://api.builtin.com/services/job-retrieval/legacy-collapsed-jobs"

var PerPage int = 100

var BuiltInURIQueries = map[string]string{
	"categories":        "149",
	"subcategories":     "",
	"experiences":       "1,1-3,3-5,3,5",
	"industry":          "",
	"company_sizes":     "",
	"regions":           "",
	"locations":         "",
	"remote":            "",
	"working_option":    "",
	"page":              "1",
	"per_page":          strconv.Itoa(PerPage),
	"search":            "",
	"sortStrategy":      "recency",
	"job_locations":     "",
	"company_locations": "",
	"jobs_board":        "true",
	"hybridEnabled":     "true",
	"national":          "true",
}

var BUILTIN_DATA_PATH = "./data/builtin"
