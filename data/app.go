package data

import (
	"fmt"
	"os"
	"strconv"
	"tyk/tyk/bootstrap/constants"
)

type AppArguments struct {
	DashboardHost                 string
	DashboardPort                 int
	DashBoardLicense              string
	TykAdminSecret                string
	CurrentOrgName                string
	TykAdminPassword              string
	Cname                         string
	TykAdminFirstName             string
	TykAdminLastName              string
	TykAdminEmailAddress          string
	UserAuth                      string
	OrgId                         string
	CatalogId                     string
	DashboardUrl                  string
	DashboardProto                string
	TykPodNamespace               string
	DashboardSvc                  string
	IsDashboardEnabled            bool
	OperatorSecretEnabled         bool
	OperatorSecretName            string
	EnterprisePortalSecretEnabled bool
	EnterprisePortalSecretName    string
	GatewayAdress                 string
	BootstrapPortal               bool
	DashboardDeploymentName       string
}

var AppConfig = AppArguments{
	IsDashboardEnabled:            false,
	OperatorSecretEnabled:         false,
	EnterprisePortalSecretEnabled: false,
	BootstrapPortal:               false,
	DashboardProto:                "",
	DashboardHost:                 "",
	DashboardPort:                 3000,
	DashBoardLicense:              "",
	TykAdminSecret:                "12345",
	CurrentOrgName:                "TYKTYK",
	Cname:                         "tykCName",
	TykAdminPassword:              "123456",
	TykAdminFirstName:             "firstName",
	TykAdminEmailAddress:          "tyk@tyk.io",
	TykAdminLastName:              "lastName",
	UserAuth:                      "",
	OrgId:                         "",
	CatalogId:                     "",
	DashboardUrl:                  "",
	TykPodNamespace:               "",
	DashboardSvc:                  "",
	OperatorSecretName:            "",
	EnterprisePortalSecretName:    "",
	GatewayAdress:                 "",
	DashboardDeploymentName:       "",
}

func InitAppDataPreDelete() error {
	AppConfig.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)
	AppConfig.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)
	AppConfig.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)
	return nil
}

func InitAppDataPostInstall() error {
	AppConfig.TykAdminFirstName = os.Getenv(constants.TykAdminFirstNameEnvVar)
	AppConfig.TykAdminLastName = os.Getenv(constants.TykAdminLastNameEnvVar)
	AppConfig.TykAdminEmailAddress = os.Getenv(constants.TykAdminEmailEnvVar)
	AppConfig.TykAdminPassword = os.Getenv(constants.TykAdminPasswordEnvVar)
	AppConfig.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)
	AppConfig.DashboardProto = os.Getenv(constants.TykDashboardProtoEnvVar)
	AppConfig.DashboardSvc = os.Getenv(constants.TykDashboardSvcEnvVar)
	dbPort, err := strconv.ParseInt(os.Getenv(constants.TykDbListenport), 10, 64)
	if err != nil {
		return err
	}
	AppConfig.DashboardPort = int(dbPort)
	AppConfig.DashBoardLicense = os.Getenv(constants.TykDbLicensekeyEnvVar)
	AppConfig.TykAdminSecret = os.Getenv(constants.TykAdminSecretEnvVar)
	AppConfig.CurrentOrgName = os.Getenv(constants.TykOrgNameEnvVar)
	AppConfig.Cname = os.Getenv(constants.TykOrgCnameEnvVar)
	AppConfig.DashboardUrl = GetDashboardUrl()
	dashEnabledRaw := os.Getenv(constants.DashboardEnabledEnvVar)
	if dashEnabledRaw != "" {
		AppConfig.IsDashboardEnabled, err = strconv.ParseBool(os.Getenv(constants.DashboardEnabledEnvVar))
		if err != nil {
			return err
		}
	}

	operatorSecretEnabledRaw := os.Getenv(constants.OperatorSecretEnabledEnvVar)
	if operatorSecretEnabledRaw != "" {
		AppConfig.OperatorSecretEnabled, err = strconv.ParseBool(operatorSecretEnabledRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)

	enterprisePortalSecretEnabledRaw := os.Getenv(constants.EnterprisePortalSecretEnabledEnvVar)
	if enterprisePortalSecretEnabledRaw != "" {
		AppConfig.EnterprisePortalSecretEnabled, err = strconv.ParseBool(enterprisePortalSecretEnabledRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)

	AppConfig.GatewayAdress = os.Getenv(constants.GatewayAddressEnvVar)
	bootstrapPortalBoolRaw := os.Getenv(constants.BootstrapPortalEnvVar)
	if bootstrapPortalBoolRaw != "" {
		AppConfig.BootstrapPortal, err = strconv.ParseBool(bootstrapPortalBoolRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.DashboardDeploymentName = os.Getenv(constants.TykDashboardDeployEnvVar)

	return nil
}

func GetDashboardUrl() string {
	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
		AppConfig.DashboardProto,
		AppConfig.DashboardSvc,
		AppConfig.TykPodNamespace,
		AppConfig.DashboardPort)
}
