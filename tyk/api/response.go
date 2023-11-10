package api

type DashboardAPIResp struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    string `json:"Meta"`
}

type OrgAPIResp struct {
	Organisations []map[string]interface{} `json:"organisations"`
	Pages         int                      `json:"pages"`
}

type CreateUserResp struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    struct {
		ID string `json:"id"`
	} `json:"Meta"`
}
