package data

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const prefix = "TYK_K8SBOOTSTRAP"

type Config struct {
	// InsecureSkipVerify enables InsecureSkipVerify options in http request sent to Tyk - might be useful
	// for Tyk Dashboard with self-signed certs.
	InsecureSkipVerify bool
	// BootstrapDashboard controls bootstrapping Tyk Dashboard or not.
	BootstrapDashboard bool
	// BootstrapPortal controls bootstrapping Tyk Classic Portal or not.
	BootstrapPortal bool

	// OperatorKubernetesSecretName corresponds to the Kubernetes secret name that will be created for Tyk Operator.
	// Set it to an empty to string to disable bootstrapping Kubernetes secret for Tyk Operator.
	// By default, tyk-operator-conf
	OperatorKubernetesSecretName string
	// DevPortalKubernetesSecretName corresponds to the Kubernetes secret name that will be created for
	// Tyk Developer Enterprise Portal. Set it to an empty to string to disable bootstrapping Kubernetes
	// secret for Tyk Developer Enterprise Portal.
	// By default, tyk-dev-portal-conf
	DevPortalKubernetesSecretName string
	// K8s consists of configurations for Kubernetes services of Tyk.
	K8s K8sConf
	// Tyk consists of configurations for Tyk components such as Tyk Dashboard Admin information
	// or Tyk Portal configurations.
	Tyk TykConf
}

type K8sConf struct {
	// DashboardSvcUrl corresponds to the URL of Tyk Dashboard.
	DashboardSvcUrl string
	// DashboardSvcProto corresponds to Tyk Dashboard Service Protocol (either http or https).
	// By default, it is http.
	DashboardSvcProto string
	// ReleaseNamespace corresponds to the namespace where Tyk is deployed via Helm Chart.
	ReleaseNamespace string
	// DashboardDeploymentName corresponds to the name of the Tyk Dashboard Deployment, which is being used
	// to restart Dashboard pod after bootstrapping. By default, it discovers Dashboard Deployment name.
	// If the environment variable is populated, the discovery will not be triggered.
	DashboardDeploymentName string
}

type TykAdmin struct {
	// Secret corresponds to the secret that will be used in Admin APIs.
	Secret string
	// FirstName corresponds to the first name of the admin being created.
	FirstName string
	// LastName corresponds to the last name of the admin being created.
	LastName string
	// EmailAddress corresponds to the email address of the admin being created.
	EmailAddress string
	// Password corresponds to the password of the admin being created.
	Password string

	// Auth corresponds to Tyk Dashboard API Access Credentials of the admin user, and it will be used
	// in Authorization header of the HTTP requests that will be sent to Tyk for bootstrapping.
	// Also, if bootstrapping Tyk Operator Secret is enabled Auth corresponds to TykAuth field in the
	// Kubernetes secret of Tyk Operator.
	Auth string `ignored:"true"`
}

type TykOrg struct {
	// Name corresponds to the name for your organization that is going to be bootstrapped in Tyk
	Name string
	// Cname corresponds to the Organisation CNAME which is going to bind the Portal to.
	Cname string

	// ID corresponds to the organisation ID that is being created.
	ID string `ignored:"true"`
}

type TykConf struct {
	// Admin consists of configurations for Tyk Dashboard Admin.
	Admin TykAdmin
	// Org consists of configurations for the organisation that is going to be created in Tyk Dashboard.
	Org TykOrg

	// DashboardLicense corresponds to the license key of Tyk Dashboard.
	DashboardLicense string
}

var BootstrapConf = Config{}

func InitBootstrapConf() error {
	return envconfig.Process(prefix, &BootstrapConf)
}

func InitPostInstall() error {
	err := InitBootstrapConf()
	if err != nil {
		return err
	}

	if BootstrapConf.BootstrapDashboard {
		dashURL, err := discoverDashboardSvc()
		if err != nil {
			return err
		}

		BootstrapConf.K8s.DashboardSvcUrl = dashURL
	}

	return nil
}

// discoverDashboardSvc lists Service objects with constants.TykBootstrapReleaseLabel label that has
// constants.TykBootstrapDashboardSvcLabel value and returns a service URL for Tyk Dashboard.
func discoverDashboardSvc() (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	ls := metav1.LabelSelector{MatchLabels: map[string]string{
		TykBootstrapLabel: TykBootstrapDashboardSvcLabel,
	}}

	l := labels.Set(ls.MatchLabels).String()

	services, err := c.
		CoreV1().
		Services(BootstrapConf.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: l})
	if err != nil {
		return "", err
	}

	if len(services.Items) == 0 {
		return "", fmt.Errorf("failed to find services with label %v\n", l)
	}

	if len(services.Items) > 1 {
		fmt.Printf("[WARNING] Found multiple services with label %v\n", l)
	}

	service := services.Items[0]
	if len(service.Spec.Ports) == 0 {
		return "", fmt.Errorf("svc/%v/%v has no open ports\n", service.Name, service.Namespace)
	}

	if len(service.Spec.Ports) > 1 {
		fmt.Printf("[WARNING] Found multiple open ports in svc/%v/%v\n", service.Name, service.Namespace)
	}

	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
		BootstrapConf.K8s.DashboardSvcProto,
		service.Name,
		BootstrapConf.K8s.ReleaseNamespace,
		service.Spec.Ports[0].Port,
	), nil
}
