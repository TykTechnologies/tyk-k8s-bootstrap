package data

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type BootstrapConf struct {
	// Tyk Configurations
	BootstrapDashboard     bool   `json:"bootstrapDashboard,omitempty"`
	BootstrapClassicPortal bool   `json:"bootstrapClassicPortal,omitempty"`
	TykDashboardLicense    string `json:"tykDashboardLicense,omitempty"`
	TykAdminFirstName      string `json:"tykAdminFirstName,omitempty"`
	TykAdminLastName       string `json:"tykAdminLastName,omitempty"`
	TykAdminEmailAddress   string `json:"tykAdminEmailAddress,omitempty"`
	TykAdminPassword       string `json:"tykAdminPassword,omitempty"`
	TykAdminSecret         string `json:"tykAdminSecret,omitempty"`
	TykUserAuth            string `json:"tykUserAuth,omitempty"`
	TykOrgId               string `json:"tykOrgId,omitempty"`
	TykOrgName             string `json:"tykOrgName,omitempty"`
	TykPortalCname         string `json:"tykPortalCname,omitempty"`

	// Tyk K8s
	DashboardSvcPort              int32  `json:"dashboardSvcPort,omitempty"`
	DashboardSvcAddr              string `json:"dashboardSvcAddr,omitempty"`
	DashboardSvcProto             string `json:"dashboardSvcProto,omitempty"`
	DashboardSvcName              string `json:"dashboardSvcName,omitempty"`
	DashboardDeploymentName       string `json:"dashboardDeploymentName,omitempty"`
	OperatorSecretEnabled         bool   `json:"operatorSecretEnabled,omitempty"`
	OperatorSecretName            string `json:"operatorSecretName,omitempty"`
	EnterprisePortalSecretEnabled bool   `json:"enterprisePortalSecretEnabled,omitempty"`
	EnterprisePortalSecretName    string `json:"enterprisePortalSecretName,omitempty"`

	// Release of the Tyk Dashboard
	ReleaseName      string `json:"releaseName,omitempty"`
	ReleaseNamespace string `json:"releaseNamespace,omitempty"`

	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

func InitAppDataPreDelete() *BootstrapConf {
	args := BootstrapConf{}
	args.OperatorSecretName = os.Getenv(OperatorSecretNameEnvVar)
	args.EnterprisePortalSecretName = os.Getenv(EnterprisePortalSecretNameEnvVar)
	args.ReleaseNamespace = os.Getenv(TykPodNamespaceEnvVar)

	return &args
}

func InitAppDataPostInstall() (*BootstrapConf, error) {
	conf := &BootstrapConf{}
	err := envconfig.Process(environmentVarPrefix, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
