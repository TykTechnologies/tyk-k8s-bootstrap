package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
	"tyk/tyk/bootstrap/tyk/data"
)

type Client struct {
	AppArgs *data.AppArguments

	clientSet *kubernetes.Clientset
}

// NewClient returns a new Client to interact with Kubernetes. It first tries instantiate in-cluster client; otherwise,
// returns client via reading default kubeconfig located /$HOME/.kube/config
func NewClient() (*Client, error) {
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

	return cl, nil
}

func (c *Client) RestartDashboardDeployment() error {
	deploymentsClient := c.clientSet.
		AppsV1().
		Deployments(c.AppArgs.TykPodNamespace)

	timeStamp := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`,
		time.Now().Format("20060102150405"))

	_, err := deploymentsClient.Patch(context.TODO(), c.AppArgs.DashboardDeploymentName,
		types.StrategicMergePatchType, []byte(timeStamp), metav1.PatchOptions{})

	return err
}
