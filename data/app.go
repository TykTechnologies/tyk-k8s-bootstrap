package data

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
)

type AppArguments struct {
	debugging bool `json:"debugging"`

	DashboardHost                 string
	DashboardPort                 int32
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
	DashboardInsecureSkipVerify   bool
	IsDashboardEnabled            bool
	OperatorSecretEnabled         bool
	OperatorSecretName            string
	EnterprisePortalSecretEnabled bool
	EnterprisePortalSecretName    string
	BootstrapPortal               bool
	DashboardDeploymentName       string
	ReleaseName                   string
}

func InitAppDataPreDelete() *AppArguments {
	args := AppArguments{}
	args.OperatorSecretName = os.Getenv(OperatorSecretNameEnvVar)
	args.EnterprisePortalSecretName = os.Getenv(EnterprisePortalSecretNameEnvVar)
	args.TykPodNamespace = os.Getenv(TykPodNamespaceEnvVar)

	return &args
}

func InitAppDataPostInstall() (*AppArguments, error) {
	appArgs := &AppArguments{}

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

	appArgs.DashBoardLicense = os.Getenv(TykDbLicensekeyEnvVar)
	appArgs.TykAdminSecret = os.Getenv(TykAdminSecretEnvVar)
	appArgs.CurrentOrgName = os.Getenv(TykOrgNameEnvVar)
	appArgs.Cname = os.Getenv(TykOrgCnameEnvVar)
	appArgs.ReleaseName = os.Getenv(ReleaseNameEnvVar)

	var err error

	dashEnabledRaw := os.Getenv(DashboardEnabledEnvVar)
	if dashEnabledRaw != "" {
		appArgs.IsDashboardEnabled, err = strconv.ParseBool(os.Getenv(DashboardEnabledEnvVar))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v, err: %v", DashboardEnabledEnvVar, err)
		}
	}

	if appArgs.IsDashboardEnabled {
		if err := discoverDashboardSvc(appArgs); err != nil {
			return nil, err
		}

		appArgs.DashboardUrl = fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
			appArgs.DashboardProto,
			appArgs.DashboardSvc,
			appArgs.TykPodNamespace,
			appArgs.DashboardPort,
		)
	}

	operatorSecretEnabledRaw := os.Getenv(OperatorSecretEnabledEnvVar)
	if operatorSecretEnabledRaw != "" {
		appArgs.OperatorSecretEnabled, err = strconv.ParseBool(operatorSecretEnabledRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v, err: %v", OperatorSecretEnabledEnvVar, err)
		}
	}
	appArgs.OperatorSecretName = os.Getenv(OperatorSecretNameEnvVar)

	enterprisePortalSecretEnabledRaw := os.Getenv(EnterprisePortalSecretEnabledEnvVar)
	if enterprisePortalSecretEnabledRaw != "" {
		appArgs.EnterprisePortalSecretEnabled, err = strconv.ParseBool(enterprisePortalSecretEnabledRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v, err: %v", EnterprisePortalSecretEnabledEnvVar, err)
		}
	}
	appArgs.EnterprisePortalSecretName = os.Getenv(EnterprisePortalSecretNameEnvVar)

	bootstrapPortalBoolRaw := os.Getenv(BootstrapPortalEnvVar)
	if bootstrapPortalBoolRaw != "" {
		appArgs.BootstrapPortal, err = strconv.ParseBool(bootstrapPortalBoolRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v, err: %v", BootstrapPortalEnvVar, err)
		}
	}
	appArgs.DashboardDeploymentName = os.Getenv(TykDashboardDeployEnvVar)

	dashboardInsecureSkipVerifyRaw := os.Getenv(TykDashboardInsecureSkipVerify)
	if dashboardInsecureSkipVerifyRaw != "" {
		appArgs.DashboardInsecureSkipVerify, err = strconv.ParseBool(dashboardInsecureSkipVerifyRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v, err: %v", TykDashboardInsecureSkipVerify, err)
		}
	}

	return appArgs, nil
}

// discoverDashboardSvc lists Service objects with TykBootstrapReleaseLabel label that has
// TykBootstrapDashboardSvcLabel value and gets this Service's metadata name, and port and
// updates DashboardSvc and DashboardPort fields.
func discoverDashboardSvc(args *AppArguments) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ls := metav1.LabelSelector{MatchLabels: map[string]string{
		TykBootstrapLabel: TykBootstrapDashboardSvcLabel,
	}}
	if args.ReleaseName != "" {
		ls.MatchLabels[TykBootstrapReleaseLabel] = args.ReleaseName
	}

	l := labels.Set(ls.MatchLabels).String()

	services, err := c.
		CoreV1().
		Services(args.TykPodNamespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: l})
	if err != nil {
		return err
	}

	if len(services.Items) == 0 {
		return fmt.Errorf("failed to find services with label %v\n", l)
	}

	if len(services.Items) > 1 {
		fmt.Printf("[WARNING] Found multiple services with label %v\n", l)
	}

	service := services.Items[0]
	if len(service.Spec.Ports) == 0 {
		return fmt.Errorf("svc/%v/%v has no open ports\n", service.Name, service.Namespace)
	}
	if len(service.Spec.Ports) > 1 {
		fmt.Printf("[WARNING] Found multiple open ports in svc/%v/%v\n", service.Name, service.Namespace)
	}

	args.DashboardPort = service.Spec.Ports[0].Port
	args.DashboardSvc = service.Name

	return nil
}
