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

	kubeconfig string
	restConfig *rest.Config
	clientSet  *kubernetes.Clientset
}

// K8sClient returns a new Client to interact with Kubernetes, based on the provided kubeconfig.
// If the kubeconfig is an empty string, it uses in-cluster configuration.
func K8sClient(kubeconfig string) (*Client, error) {
	cl := &Client{kubeconfig: kubeconfig}

	var err error

	if cl.kubeconfig == "" {
		fmt.Println("using in-cluster configuration")
		cl.restConfig, err = rest.InClusterConfig()
	} else {
		fmt.Printf("using configuration from '%s'\n", kubeconfig)
		cl.restConfig, err = clientcmd.BuildConfigFromFlags("", cl.kubeconfig)
	}

	if err != nil {
		fmt.Printf("failed to generate config, err: %v", err)
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(cl.restConfig)
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
