package data

const (
	environmentVarPrefix = "TYK_K8SBOOTSTRAP"

	OperatorSecretNameEnvVar         = "OPERATOR_SECRET_NAME"
	EnterprisePortalSecretNameEnvVar = "ENTERPRISE_PORTAL_SECRET_NAME"
	TykPodNamespaceEnvVar            = "TYK_POD_NAMESPACE"
	TykDashboardLicenseEnvVarName    = "TYK_DB_LICENSEKEY"

	TykBootstrapLabel                = "tyk.tyk.io/k8s-bootstrap"
	TykBootstrapPreDeleteLabel       = "tyk-k8s-bootstrap-pre-delete"
	TykBootstrapDashboardDeployLabel = "tyk-dashboard"
	TykBootstrapDashboardSvcLabel    = "tyk-dashboard"
	TykBootstrapReleaseLabel         = "release"
)
