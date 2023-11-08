package k8s

import (
	"context"
	"fmt"
	"time"
	"tyk/tyk/bootstrap/pkg/constants"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

// RestartDashboard restarts Tyk Dashboard Deployment. Restarting Tyk Dashboard is needed to apply
// CName.
func (c *Client) RestartDashboard() error {
	if c.appArgs.K8s.DashboardDeploymentName == "" {
		ls := metav1.LabelSelector{MatchLabels: map[string]string{
			constants.TykBootstrapLabel: constants.TykBootstrapDashboardDeployLabel,
		}}

		deployments, err := c.clientSet.
			AppsV1().
			Deployments(c.appArgs.K8s.ReleaseNamespace).
			List(
				context.TODO(),
				metav1.ListOptions{
					LabelSelector: labels.Set(ls.MatchLabels).String(),
				},
			)
		if err != nil {
			return fmt.Errorf("failed to list Tyk Dashboard Deployment, err: %v", err)
		}

		for i := range deployments.Items {
			c.appArgs.K8s.DashboardDeploymentName = deployments.Items[i].ObjectMeta.Name
		}
	}

	timeStamp := fmt.Sprintf(
		`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`,
		time.Now().Format("20060102150405"),
	)

	_, err := c.clientSet.
		AppsV1().
		Deployments(c.appArgs.K8s.ReleaseNamespace).
		Patch(
			context.TODO(),
			c.appArgs.K8s.DashboardDeploymentName,
			types.StrategicMergePatchType,
			[]byte(timeStamp),
			metav1.PatchOptions{},
		)

	return err
}

// discoverDashboardSvc lists Service objects with constants.TykBootstrapLabel label that has
// constants.TykBootstrapDashboardSvcLabel value and returns a service URL for Tyk Dashboard.
func (c *Client) discoverDashboardSvc() (string, error) {
	ls := metav1.LabelSelector{MatchLabels: map[string]string{
		constants.TykBootstrapLabel: constants.TykBootstrapDashboardSvcLabel,
	}}

	l := labels.Set(ls.MatchLabels).String()

	services, err := c.clientSet.
		CoreV1().
		Services(c.appArgs.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: l})
	if err != nil {
		return "", err
	}

	if len(services.Items) == 0 {
		return "", fmt.Errorf("failed to find services with label %v\n", l)
	}

	if len(services.Items) > 1 {
		c.l.Warnf("Found multiple services with label %v", l)
	}

	service := services.Items[0]
	if len(service.Spec.Ports) == 0 {
		return "", fmt.Errorf("svc/%v/%v has no open ports\n", service.Name, service.Namespace)
	}

	if len(service.Spec.Ports) > 1 {
		c.l.Warnf("Found multiple open ports in svc/%v/%v", service.Name, service.Namespace)
	}

	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
		c.appArgs.K8s.DashboardSvcProto,
		service.Name,
		c.appArgs.K8s.ReleaseNamespace,
		service.Spec.Ports[0].Port,
	), nil
}
