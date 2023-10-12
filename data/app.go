package data

import (
	"fmt"
	"os"
	"strconv"
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
	args.OperatorSecretName = os.Getenv(OperatorSecretNameEnvVar)
	args.EnterprisePortalSecretName = os.Getenv(EnterprisePortalSecretNameEnvVar)
	args.TykPodNamespace = os.Getenv(TykPodNamespaceEnvVar)

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

	debugging := os.Getenv(DebuggingEnvVar)
	if debugging != "" {
		appArgs.debugging, _ = strconv.ParseBool(debugging)
	}

	appArgs.TykAdminFirstName = os.Getenv(TykAdminFirstNameEnvVar)
	appArgs.TykAdminLastName = os.Getenv(TykAdminLastNameEnvVar)
	appArgs.TykAdminEmailAddress = os.Getenv(TykAdminEmailEnvVar)
	appArgs.TykAdminPassword = os.Getenv(TykAdminPasswordEnvVar)
	appArgs.TykPodNamespace = os.Getenv(TykPodNamespaceEnvVar)
	appArgs.DashboardProto = os.Getenv(TykDashboardProtoEnvVar)
	appArgs.DashboardSvc = os.Getenv(TykDashboardSvcEnvVar)
	dbPort, err := strconv.ParseInt(os.Getenv(TykDbListenport), 10, 64)
	if err != nil {
		return nil, err
	}
	appArgs.DashboardPort = int(dbPort)
	appArgs.DashBoardLicense = os.Getenv(TykDbLicensekeyEnvVar)
	appArgs.TykAdminSecret = os.Getenv(TykAdminSecretEnvVar)
	appArgs.CurrentOrgName = os.Getenv(TykOrgNameEnvVar)
	appArgs.Cname = os.Getenv(TykOrgCnameEnvVar)
	appArgs.DashboardUrl = appArgs.dashboardURL()

	dashEnabledRaw := os.Getenv(DashboardEnabledEnvVar)
	if dashEnabledRaw != "" {
		appArgs.IsDashboardEnabled, err = strconv.ParseBool(os.Getenv(DashboardEnabledEnvVar))
		if err != nil {
			return nil, err
		}
	}

	operatorSecretEnabledRaw := os.Getenv(OperatorSecretEnabledEnvVar)
	if operatorSecretEnabledRaw != "" {
		appArgs.OperatorSecretEnabled, err = strconv.ParseBool(operatorSecretEnabledRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.OperatorSecretName = os.Getenv(OperatorSecretNameEnvVar)

	enterprisePortalSecretEnabledRaw := os.Getenv(EnterprisePortalSecretEnabledEnvVar)
	if enterprisePortalSecretEnabledRaw != "" {
		appArgs.EnterprisePortalSecretEnabled, err = strconv.ParseBool(enterprisePortalSecretEnabledRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.EnterprisePortalSecretName = os.Getenv(EnterprisePortalSecretNameEnvVar)

	appArgs.GatewayAddress = os.Getenv(GatewayAddressEnvVar)
	bootstrapPortalBoolRaw := os.Getenv(BootstrapPortalEnvVar)
	if bootstrapPortalBoolRaw != "" {
		appArgs.BootstrapPortal, err = strconv.ParseBool(bootstrapPortalBoolRaw)
		if err != nil {
			return nil, err
		}
	}
	appArgs.DashboardDeploymentName = os.Getenv(TykDashboardDeployEnvVar)

	dashboardInsecureSkipVerifyRaw := os.Getenv(TykDashboardInsecureSkipVerify)
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
