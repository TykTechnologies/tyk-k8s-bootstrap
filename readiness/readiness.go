package readiness

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tyk/tyk/bootstrap/data"

	v1 "k8s.io/api/core/v1"
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

		pods, err := clientset.CoreV1().Pods(data.BootstrapConf.K8s.ReleaseNamespace).
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		fmt.Printf("There are %d other pods in the cluster\n", len(pods.Items)-1)

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
			return nil
		}

		fmt.Printf("The following pods have containers that are NOT ready: ")

		for podName, _ := range notReadyPods {
			fmt.Println(podName)
		}

		time.Sleep(2 * time.Second)
	}
}
