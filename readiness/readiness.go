package readiness

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"time"
)

func CheckIfDeploymentsAreReady() bool {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods(os.Getenv("TYK_POD_NAMESPACE")).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	time.Sleep(5 * time.Second)

	var attemptCount int
	for {
		attemptCount++
		if attemptCount > 180 {
			return false
		}
		pods, err = clientset.CoreV1().Pods(os.Getenv("TYK_POD_NAMESPACE")).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
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
			fmt.Println("NOT READY")
		} else {
			fmt.Println("TYK PRO IS READY TO BE BOOTSTRAPPED")
			return true
		}
		time.Sleep(2 * time.Second)
	}
}
