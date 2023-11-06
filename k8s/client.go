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

//func (c *Client) RestartDashboardDeployment() error {
//	config, err := rest.InClusterConfig()
//	if err != nil {
//		return err
//	}
//
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return err
//	}
//
//	if c.AppArgs.DashboardDeploymentName == "" {
//		ls := metav1.LabelSelector{MatchLabels: map[string]string{
//			data.TykBootstrapLabel: data.TykBootstrapDashboardDeployLabel,
//		}}
//
//		if c.AppArgs.ReleaseName != "" {
//			ls.MatchLabels[data.TykBootstrapReleaseLabel] = c.AppArgs.ReleaseName
//		}
//
//		deployments, err := clientset.
//			AppsV1().
//			Deployments(c.AppArgs.ReleaseNamespace).
//			List(
//				context.TODO(),
//				metav1.ListOptions{
//					LabelSelector: labels.Set(ls.MatchLabels).String(),
//				},
//			)
//		if err != nil {
//			return errors.New(fmt.Sprintf("failed to list Tyk Dashboard Deployment, err: %v", err))
//		}
//
//		for _, deployment := range deployments.Items {
//			c.AppArgs.DashboardDeploymentName = deployment.ObjectMeta.Name
//		}
//	}
//
//	timeStamp := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`,
//		time.Now().Format("20060102150405"))
//
//	_, err = clientset.
//		AppsV1().
//		Deployments(c.AppArgs.ReleaseName).
//		Patch(
//			context.TODO(),
//			c.AppArgs.DashboardDeploymentName,
//			types.StrategicMergePatchType,
//			[]byte(timeStamp),
//			metav1.PatchOptions{},
//		)
//
//	return err
//}
//
//// discoverDashboardSvc lists Service objects with TykBootstrapReleaseLabel label that has
//// TykBootstrapDashboardSvcLabel value and gets this Service's metadata name, and port and
//// updates DashboardSvcName and DashboardSvcPort fields.
//func (c *Client) discoverDashboardSvc() error {
//	ls := metav1.LabelSelector{MatchLabels: map[string]string{
//		data.TykBootstrapLabel: data.TykBootstrapDashboardSvcLabel,
//	}}
//	if c.AppArgs.ReleaseName != "" {
//		ls.MatchLabels[data.TykBootstrapReleaseLabel] = c.AppArgs.ReleaseName
//	}
//
//	l := labels.Set(ls.MatchLabels).String()
//
//	services, err := c.clientSet.
//		CoreV1().
//		Services(c.AppArgs.ReleaseNamespace).
//		List(context.TODO(), metav1.ListOptions{LabelSelector: l})
//	if err != nil {
//		return err
//	}
//
//	if len(services.Items) == 0 {
//		return fmt.Errorf("failed to find services with label %v\n", l)
//	}
//
//	if len(services.Items) > 1 {
//		fmt.Printf("[WARNING] Found multiple services with label %v\n", l)
//	}
//
//	service := services.Items[0]
//	if len(service.Spec.Ports) == 0 {
//		return fmt.Errorf("svc/%v/%v has no open ports\n", service.Name, service.Namespace)
//	}
//	if len(service.Spec.Ports) > 1 {
//		fmt.Printf("[WARNING] Found multiple open ports in svc/%v/%v\n", service.Name, service.Namespace)
//	}
//
//	c.AppArgs.DashboardSvcPort = service.Spec.Ports[0].Port
//	c.AppArgs.DashboardSvcName = service.Name
//
//	return nil
//}
