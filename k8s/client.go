package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"tyk/tyk/bootstrap/pkg/config"
)

type Client struct {
	appArgs   *config.Config
	clientSet *kubernetes.Clientset
}

// NewClient returns a new Client to interact with Kubernetes. It first tries to instantiate in-cluster client;
// otherwise, returns client via reading default kubeconfig located /$HOME/.kube/config
func NewClient(conf *config.Config) (*Client, error) {
	cl := &Client{}

	var err error
	var config *rest.Config

	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename())
		if err != nil {
			return nil, err
		}
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("failed to generate client, err: %v", err)
		return nil, err
	}

	cl.clientSet = cs

	if conf.BootstrapDashboard {
		dashURL, err := cl.discoverDashboardSvc()
		if err != nil {
			return nil, err
		}

		conf.K8s.DashboardSvcUrl = dashURL
	}

	cl.appArgs = conf

	return cl, nil
}
