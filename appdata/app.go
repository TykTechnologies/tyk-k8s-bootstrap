package appdata

type AppArg struct {
	DashboardHost        string
	DashboardPort        int
	DashBoardLicense     string
	TykAdminSecret       string
	CurrentOrgName       string
	TykAdminPassword     string
	Cname                string
	TykAdminFirstName    string
	TykAdminLastName     string
	TykAdminEmailAddress string
	UserAuth             string
	OrgId                string
	CatalogId            string
	DashboardUrl         string
}

var Config = AppArg{
	DashboardHost:        "",
	DashboardPort:        3000,
	DashBoardLicense:     "",
	TykAdminSecret:       "12345",
	CurrentOrgName:       "TYKTYK",
	Cname:                "tykCName",
	TykAdminPassword:     "123456",
	TykAdminFirstName:    "andrei",
	TykAdminEmailAddress: "andrei@tyk.io",
	TykAdminLastName:     "psk",
	UserAuth:             "",
	OrgId:                "",
	CatalogId:            "",
	DashboardUrl:         "",
}
