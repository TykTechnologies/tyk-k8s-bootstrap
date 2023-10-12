package predelete

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"os"
	"tyk/tyk/bootstrap/constants"
	"tyk/tyk/bootstrap/data"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// maybe not needed?
func ExecutePreDeleteOperations() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	err = PreDeleteOperatorSecret(clientset)
	if err != nil {
		return err
	}

	err = PreDeleteEnterprisePortalSecret(clientset)
	if err != nil {
		return err
	}

	err = PreDeleteBootstrappingJobs(clientset)
	if err != nil {
		return err
	}

	return nil
}

func PreDeleteOperatorSecret(clientset *kubernetes.Clientset) error {
	fmt.Println("Running pre delete hook")
	secrets, err := clientset.CoreV1().Secrets(os.Getenv("TYK_POD_NAMESPACE")).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	found := false
	for _, value := range secrets.Items {
		if value.Name == os.Getenv("OPERATOR_SECRET_NAME") {
			err = clientset.CoreV1().Secrets(os.Getenv("TYK_POD_NAMESPACE")).Delete(context.TODO(), value.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
			found = true
			break
		}
	}

	if !found {
		fmt.Println("A previously created operator secret has not been identified")
	} else {
		fmt.Println("A previously created operator secret was identified and deleted")
	}
	return nil
}

func PreDeleteEnterprisePortalSecret(clientset *kubernetes.Clientset) error {
	fmt.Println("Running pre delete hook")
	ns := data.AppConfig.TykPodNamespace

	secrets, err := clientset.CoreV1().Secrets(ns).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	notFound := true
	for _, value := range secrets.Items {
		if data.AppConfig.EnterprisePortalSecretName == value.Name {
			err = clientset.CoreV1().Secrets(ns).
				Delete(context.TODO(), value.Name, metav1.DeleteOptions{})

			if err != nil {
				return err
			}
			fmt.Println("A previously created enterprise portal secret was identified and deleted")
			notFound = false
			break
		}
	}

	if notFound {
		fmt.Println("A previously created enterprise portal secret has not been identified")
	}

	return nil
}

func PreDeleteBootstrappingJobs(clientset *kubernetes.Clientset) error {
	ls := metav1.LabelSelector{
		MatchLabels: map[string]string{
			constants.TykBootstrapLabel: constants.TykBootstrapPostInstallLabel,
		},
	}

	jobs, err := clientset.
		BatchV1().
		Jobs(data.AppConfig.TykPodNamespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(ls.MatchLabels).String()})
	if err != nil {
		return err
	}

	found := false
	for _, value := range jobs.Items {
		deletePropagationType := metav1.DeletePropagationBackground

		err = clientset.
			BatchV1().
			Jobs(data.AppConfig.TykPodNamespace).
			Delete(context.TODO(), value.Name, metav1.DeleteOptions{PropagationPolicy: &deletePropagationType})
		if err != nil {
			return err
		}
		found = true
		break
	}

	if !found {
		fmt.Println("A previously created bootstrapping job has not been identified")
	} else {
		fmt.Println("A previously created bootstrapping job was identified and deleted")
	}
	return nil
}
