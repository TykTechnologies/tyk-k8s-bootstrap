package k8s

import (
	"context"
	"fmt"
	"os"
	"tyk/tyk/bootstrap/data"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ExecutePreDeleteOperations() error {
	err := c.deleteOperatorSecret()
	if err != nil {
		return err
	}

	err = c.deleteEnterprisePortalSecret()
	if err != nil {
		return err
	}

	err = c.deleteBootstrappingJobs()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteOperatorSecret() error {
	fmt.Println("Running pre delete hook")
	secrets, err := c.clientSet.CoreV1().Secrets(os.Getenv("TYK_POD_NAMESPACE")).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	found := false
	for _, value := range secrets.Items {
		if value.Name == os.Getenv("OPERATOR_SECRET_NAME") {
			err = c.clientSet.CoreV1().Secrets(os.Getenv("TYK_POD_NAMESPACE")).Delete(context.TODO(), value.Name, metav1.DeleteOptions{})
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

func (c *Client) deleteEnterprisePortalSecret() error {
	fmt.Println("Running pre delete hook")
	ns := c.AppArgs.TykPodNamespace

	secrets, err := c.clientSet.CoreV1().Secrets(ns).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	notFound := true
	for _, value := range secrets.Items {
		if c.AppArgs.EnterprisePortalSecretName == value.Name {
			err = c.clientSet.CoreV1().Secrets(ns).
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

// deleteBootstrappingJobs deletes all jobs within the release namespace, that has specific label.
func (c *Client) deleteBootstrappingJobs() error {
	// Usually, the raw strings in label selectors are not recommended.
	jobs, err := c.clientSet.
		BatchV1().
		Jobs(c.AppArgs.TykPodNamespace).
		List(
			context.TODO(),
			metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s", data.TykBootstrapLabel),
			},
		)
	if err != nil {
		return err
	}

	var errCascading error
	for _, job := range jobs.Items {
		// Do not need to delete pre-delete job. It will be deleted by Helm.
		jobLabel := job.ObjectMeta.Labels[data.TykBootstrapLabel]
		if jobLabel != data.TykBootstrapPreDeleteLabel {
			deletePropagationType := metav1.DeletePropagationBackground

			err2 := c.clientSet.
				BatchV1().
				Jobs(c.AppArgs.TykPodNamespace).
				Delete(context.TODO(), job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePropagationType})
			if err2 != nil {
				errCascading = err2
			}
		}
	}

	return errCascading
}
