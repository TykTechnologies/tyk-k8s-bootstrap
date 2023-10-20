package readiness

import (
	"context"
	"errors"
	"fmt"
	"github.com/TykTechnologies/tyk-k8s-bootstrap/data"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CheckIfRequiredDeploymentsAreReady() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	var attemptCount int
	for {
		attemptCount++
		if attemptCount > 180 {
			return errors.New("attempted readiness check too many times")
		}
		pods, err := clientset.CoreV1().Pods(data.AppConfig.TykPodNamespace).
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("There are %d other pods in the cluster\n", len(pods.Items)-1)

		var requiredPods []v1.Pod
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "dashboard") ||
				strings.Contains(pod.Name, "redis") {
				requiredPods = append(requiredPods, pod)
			}
		}

		notReadyPods := make(map[string]struct{})
		for _, pod := range requiredPods {
			podStatus := pod.Status
			for container := range pod.Spec.Containers {
				if !podStatus.ContainerStatuses[container].Ready {
					notReadyPods[pod.Name] = struct{}{}
				}
			}
		}

		if len(notReadyPods) == 0 {
			return nil
		}

		fmt.Printf("The following pods have containers that are NOT ready: ")
		for pod, _ := range notReadyPods {
			fmt.Println(pod)
		}

		time.Sleep(2 * time.Second)
	}
}
