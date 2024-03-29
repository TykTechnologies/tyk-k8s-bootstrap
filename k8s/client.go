package k8s

import (
	"tyk/tyk/bootstrap/pkg/config"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	appArgs   *config.Config
	clientSet *kubernetes.Clientset
	l         *logrus.Entry
}

// NewClient returns a new Client to interact with Kubernetes. It first tries to instantiate in-cluster client;
// otherwise, returns client via reading default kubeconfig located /$HOME/.kube/config
func NewClient(conf *config.Config, l *logrus.Entry) (*Client, error) {
	cl := &Client{appArgs: conf, l: l}

	var err error
	var config *rest.Config

	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags(
			"",
			clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename(),
		)
		if err != nil {
			return nil, err
		}
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		cl.l.Errorf("failed to generate client, err: %v", err)
		return nil, err
	}

	cl.clientSet = cs

	return cl, nil
}
