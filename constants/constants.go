package constants

const (
	OperatorSecretEnabledEnvVar        = "OPERATOR_SECRET_ENABLED"
	DeveloperPortalSecretEnabledEnvVar = "DEVELOPER_PORTAL_SECRET_ENABLED"
	BootstrapPortalEnvVar              = "BOOTSTRAP_PORTAL"
	TykDashboardDeployEnvVar           = "TYK_DASHBOARD_DEPLOY"
	OperatorSecretNameEnvVar           = "OPERATOR_SECRET_NAME"
	DeveloperPortalSecretNameEnvVar    = "DEVELOPER_PORTAL_SECRET_NAME"
	TykAdminFirstNameEnvVar            = "TYK_ADMIN_FIRST_NAME"
	TykAdminLastNameEnvVar             = "TYK_ADMIN_LAST_NAME"
	TykAdminEmailEnvVar                = "TYK_ADMIN_EMAIL"
	TykAdminPasswordEnvVar             = "TYK_ADMIN_PASSWORD"
	TykPodNamespaceEnvVar              = "TYK_POD_NAMESPACE"
	TykDashboardProtoEnvVar            = "TYK_DASHBOARD_PROTO"
	TykDashboardInsecureSkipVerify     = "TYK_DASHBOARD_INSECURE_SKIP_VERIFY"
	TykDashboardLicenseEnvVarName      = "TYK_DB_LICENSEKEY"
	TykDbLicensekeyEnvVar              = "TYK_DB_LICENSEKEY"
	TykAdminSecretEnvVar               = "TYK_ADMIN_SECRET"
	DashboardEnabledEnvVar             = "DASHBOARD_ENABLED"
	TykOrgNameEnvVar                   = "TYK_ORG_NAME"
	TykOrgCnameEnvVar                  = "TYK_ORG_CNAME"

	TykBootstrapLabel                = "tyk.tyk.io/k8s-bootstrap"
	TykBootstrapPreDeleteLabel       = "tyk-k8s-bootstrap-pre-delete"
	TykBootstrapDashboardDeployLabel = "tyk-dashboard"
	TykBootstrapDashboardSvcLabel    = "tyk-dashboard"
)
