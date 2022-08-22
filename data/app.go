package data

import "fmt"

type AppArguments struct {
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

var AppConfig = AppArguments{
	DashboardHost:        "",
	DashboardPort:        3000,
	DashBoardLicense:     "",
	TykAdminSecret:       "12345",
	CurrentOrgName:       "TYKTYK",
	Cname:                "tykCName",
	TykAdminPassword:     "123456",
	TykAdminFirstName:    "andrei",
	TykAdminEmailAddress: "andrei@tyk.io",
	TykAdminLastName:     "ierdna",
	UserAuth:             "",
	OrgId:                "",
	CatalogId:            "",
	DashboardUrl:         "",
}

func InitAppData() {
	dashAddress := "dashboard-svc-tyk-pro.tyk.svc.cluster.local"
	dashUrl := "http://" + dashAddress + fmt.Sprintf(":%v", AppConfig.DashboardPort)
	AppConfig.DashboardUrl = dashUrl
}
