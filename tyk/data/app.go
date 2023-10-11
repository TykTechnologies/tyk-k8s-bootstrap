package data

import (
	"fmt"
	"os"
	"strconv"
	"tyk/tyk/bootstrap/constants"
)

type AppArguments struct {
	debugging bool `json:"debugging"`

	DashboardPort                 int    `json:"dashboardPort"`
	DashBoardLicense              string `json:"dashBoardLicense"`
	TykAdminSecret                string `json:"tykAdminSecret"`
	CurrentOrgName                string `json:"currentOrgName"`
	TykAdminPassword              string `json:"tykAdminPassword"`
	Cname                         string `json:"cname"`
	TykAdminFirstName             string `json:"tykAdminFirstName"`
	TykAdminLastName              string `json:"tykAdminLastName"`
	TykAdminEmailAddress          string `json:"tykAdminEmailAddress"`
	UserAuth                      string `json:"userAuth"`
	OrgId                         string `json:"orgId"`
	CatalogId                     string `json:"catalogId"`
	DashboardUrl                  string `json:"dashboardUrl"`
	DashboardProto                string `json:"dashboardProto"`
	TykPodNamespace               string `json:"tykPodNamespace"`
	DashboardSvc                  string `json:"dashboardSvc"`
	DashboardInsecureSkipVerify   bool   `json:"dashboardInsecureSkipVerify"`
	IsDashboardEnabled            bool   `json:"isDashboardEnabled"`
	OperatorSecretEnabled         bool   `json:"operatorSecretEnabled"`
	OperatorSecretName            string `json:"operatorSecretName"`
	EnterprisePortalSecretEnabled bool   `json:"enterprisePortalSecretEnabled"`
	EnterprisePortalSecretName    string `json:"enterprisePortalSecretName"`
	GatewayAddress                string `json:"gatewayAddress"`
	BootstrapPortal               bool   `json:"bootstrapPortal"`
	DashboardDeploymentName       string `json:"dashboardDeploymentName"`
}

func InitAppDataPreDelete() *AppArguments {
	args := AppArguments{}
	args.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)
	args.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)
	args.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)

	return &args
}

func InitAppDataPostInstall() (*AppArguments, error) {
	var appArgs = AppArguments{
		DashboardPort:        3000,
		TykAdminSecret:       "12345",
		CurrentOrgName:       "TYKTYK",
		Cname:                "tykCName",
		TykAdminPassword:     "123456",
		TykAdminFirstName:    "firstName",
		TykAdminEmailAddress: "tyk@tyk.io",
		TykAdminLastName:     "lastName",
	}

	debugging := os.Getenv(constants.DebuggingEnvVar)
	if debugging != "" {
		appArgs.debugging, _ = strconv.ParseBool(debugging)
	}

	appArgs.TykAdminFirstName = os.Getenv(constants.TykAdminFirstNameEnvVar)
	appArgs.TykAdminLastName = os.Getenv(constants.TykAdminLastNameEnvVar)
	appArgs.TykAdminEmailAddress = os.Getenv(constants.TykAdminEmailEnvVar)
	appArgs.TykAdminPassword = os.Getenv(constants.TykAdminPasswordEnvVar)
	appArgs.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)
	appArgs.DashboardProto = os.Getenv(constants.TykDashboardProtoEnvVar)
	appArgs.DashboardSvc = os.Getenv(constants.TykDashboardSvcEnvVar)
	dbPort, err := strconv.ParseInt(os.Getenv(constants.TykDbListenport), 10, 64)
	if err != nil {
		return nil, err
	}
	appArgs.DashboardPort = int(dbPort)
	appArgs.DashBoardLicense = os.Getenv(constants.TykDbLicensekeyEnvVar)
	appArgs.TykAdminSecret = os.Getenv(constants.TykAdminSecretEnvVar)
	appArgs.CurrentOrgName = os.Getenv(constants.TykOrgNameEnvVar)
	appArgs.Cname = os.Getenv(constants.TykOrgCnameEnvVar)
	appArgs.DashboardUrl = appArgs.dashboardURL()

	dashEnabledRaw := os.Getenv(constants.DashboardEnabledEnvVar)
	if dashEnabledRaw != "" {
		appArgs.IsDashboardEnabled, err = strconv.ParseBool(os.Getenv(constants.DashboardEnabledEnvVar))
		if err != nil {
			return nil, err
		}
	}

	operatorSecretEnabledRaw := os.Getenv(constants.OperatorSecretEnabledEnvVar)
	if operatorSecretEnabledRaw != "" {
		appArgs.OperatorSecretEnabled, err = strconv.ParseBool(operatorSecretEnabledRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)

	enterprisePortalSecretEnabledRaw := os.Getenv(constants.EnterprisePortalSecretEnabledEnvVar)
	if enterprisePortalSecretEnabledRaw != "" {
		appArgs.EnterprisePortalSecretEnabled, err = strconv.ParseBool(enterprisePortalSecretEnabledRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)

	appArgs.GatewayAddress = os.Getenv(constants.GatewayAddressEnvVar)
	bootstrapPortalBoolRaw := os.Getenv(constants.BootstrapPortalEnvVar)
	if bootstrapPortalBoolRaw != "" {
		appArgs.BootstrapPortal, err = strconv.ParseBool(bootstrapPortalBoolRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.DashboardDeploymentName = os.Getenv(constants.TykDashboardDeployEnvVar)

	dashboardInsecureSkipVerifyRaw := os.Getenv(constants.TykDashboardInsecureSkipVerify)
	if dashboardInsecureSkipVerifyRaw != "" {
		appArgs.DashboardInsecureSkipVerify, err = strconv.ParseBool(dashboardInsecureSkipVerifyRaw)
		if err != nil {
			return nil, err
		}
	}

	return &appArgs, nil
}

func (a *AppArguments) dashboardURL() string {
	if a.debugging {
		return fmt.Sprintf("%s://%s:%d", a.DashboardProto, a.DashboardSvc, a.DashboardPort)
	}

	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
		a.DashboardProto,
		a.DashboardSvc,
		a.TykPodNamespace,
		a.DashboardPort)
}

type PortalFields struct {
	JumboCTATitle       string `json:"JumboCTATitle"`
	SubHeading          string `json:"SubHeading"`
	JumboCTALink        string `json:"JumboCTALink"`
	JumboCTALinkTitle   string `json:"JumboCTALinkTitle"`
	PanelOneContent     string `json:"PanelOneContent"`
	PanelOneLink        string `json:"PanelOneLink"`
	PanelOneLinkTitle   string `json:"PanelOneLinkTitle"`
	PanelOneTitle       string `json:"PanelOneTitle"`
	PanelThereeContent  string `json:"PanelThereeContent"`
	PanelThreeContent   string `json:"PanelThreeContent"`
	PanelThreeLink      string `json:"PanelThreeLink"`
	PanelThreeLinkTitle string `json:"PanelThreeLinkTitle"`
	PanelThreeTitle     string `json:"PanelThreeTitle"`
	PanelTwoContent     string `json:"PanelTwoContent"`
	PanelTwoLink        string `json:"PanelTwoLink"`
	PanelTwoLinkTitle   string `json:"PanelTwoLinkTitle"`
	PanelTwoTitle       string `json:"PanelTwoTitle"`
}
