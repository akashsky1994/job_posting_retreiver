package model

type BuiltIn struct {
	JobCount  int `json:"job_all_count"`
	Companies []struct {
		Company BuiltInCompany `json:"company"`
		Jobs    []BuiltInJob   `json:"jobs"`
	} `json:"company_jobs"`
}

type BuiltInCompany struct {
	Title string `json:"title"`
}

type BuiltInJob struct {
	JobLink  string `json:"how_to_apply"`
	JobTitle string `json:"title"`
	Location string `json:"location"`
	Remote   string `json:"remote"`
}
