package predelete

import (
	"context"
	"fmt"
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

	err = PreDeletePortalSecret(clientset)
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
	secrets, err := clientset.
		CoreV1().
		Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	found := false
	for _, value := range secrets.Items {
		if value.Name == data.BootstrapConf.OperatorKubernetesSecretName {
			err = clientset.
				CoreV1().
				Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
				Delete(context.TODO(), value.Name, metav1.DeleteOptions{})

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

func PreDeletePortalSecret(clientset *kubernetes.Clientset) error {
	fmt.Println("Running pre delete hook")
	ns := data.BootstrapConf.K8s.ReleaseNamespace

	secrets, err := clientset.CoreV1().Secrets(ns).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	notFound := true
	for _, value := range secrets.Items {
		if data.BootstrapConf.DevPortalKubernetesSecretName == value.Name {
			err = clientset.CoreV1().Secrets(ns).
				Delete(context.TODO(), value.Name, metav1.DeleteOptions{})

			if err != nil {
				return err
			}
			fmt.Println("A previously created developer portal secret was identified and deleted")
			notFound = false
			break
		}
	}

	if notFound {
		fmt.Println("A previously created developer portal secret has not been identified")
	}

	return nil
}

// PreDeleteBootstrappingJobs deletes all jobs within the release namespace, that has specific label.
func PreDeleteBootstrappingJobs(clientset *kubernetes.Clientset) error {
	// Usually, the raw strings in label selectors are not recommended.
	jobs, err := clientset.
		BatchV1().
		Jobs(data.BootstrapConf.K8s.ReleaseNamespace).
		List(
			context.TODO(),
			metav1.ListOptions{
				LabelSelector: constants.TykBootstrapLabel,
			},
		)
	if err != nil {
		return err
	}

	var errCascading error
	for _, job := range jobs.Items {
		// Do not need to delete pre-delete job. It will be deleted by Helm.
		jobLabel := job.ObjectMeta.Labels[constants.TykBootstrapLabel]
		if jobLabel != constants.TykBootstrapPreDeleteLabel {
			deletePropagationType := metav1.DeletePropagationBackground

			err2 := clientset.
				BatchV1().
				Jobs(data.BootstrapConf.K8s.ReleaseNamespace).
				Delete(context.TODO(), job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePropagationType})
			if err2 != nil {
				errCascading = err2
			}
		}
	}

	return errCascading
}
