package readiness

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tyk/tyk/bootstrap/data"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CheckIfDeploymentsAreReady() error {
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
		var containerReady int
		var totalContainers int
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, "bootstrap") {
				continue
			}
			podStatus := pod.Status
			for container := range pod.Spec.Containers {
				if podStatus.ContainerStatuses[container].Ready {
					containerReady++
				}
				totalContainers++
			}
		}
		fmt.Printf("Ready containers: %v\n", containerReady)
		fmt.Printf("Total containers: %v\n", totalContainers)

		if containerReady != totalContainers {
			fmt.Println("TYK PRO IS NOT READY")
		} else {
			fmt.Println("TYK PRO IS READY TO BE BOOTSTRAPPED")
			return nil
		}
		time.Sleep(2 * time.Second)
	}
}
