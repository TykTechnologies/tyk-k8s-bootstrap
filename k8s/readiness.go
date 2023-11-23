package k8s

import (
	"context"
	"errors"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CheckIfRequiredDeploymentsAreReady checks if the required Deployments are ready to be bootstrapped
// or not. At the moment, it checks for Redis and Tyk Dashboard pods to be ready.
func (c *Client) CheckIfRequiredDeploymentsAreReady() error {
	time.Sleep(5 * time.Second)
	var attemptCount int

	for {
		attemptCount++
		if attemptCount > 180 {
			return errors.New("attempted readiness check too many times")
		}

		pods, err := c.clientSet.
			CoreV1().
			Pods(c.appArgs.K8s.ReleaseNamespace).
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		c.l.Infof("There are %d other pods in the cluster", len(pods.Items)-1)

		var requiredPods []v1.Pod

		for i := range pods.Items {
			pod := pods.Items[i]
			if strings.Contains(pod.Name, "dashboard") ||
				strings.Contains(pod.Name, "redis") {
				requiredPods = append(requiredPods, pod)
			}
		}

		notReadyPods := make(map[string]struct{})

		for i := range requiredPods {
			pod := requiredPods[i]
			podStatus := pod.Status

			for container := range pod.Spec.Containers {
				if !podStatus.ContainerStatuses[container].Ready {
					notReadyPods[pod.Name] = struct{}{}
				}
			}
		}

		if len(notReadyPods) == 0 {
			c.l.Info("All Pods are ready")
			return nil
		}

		for podName := range notReadyPods {
			c.l.Infof("Pod: %v has containers that are NOT ready", podName)
		}

		time.Sleep(2 * time.Second)
	}
}
