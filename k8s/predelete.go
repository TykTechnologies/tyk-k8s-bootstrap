package k8s

import (
	"context"
	"fmt"
	"tyk/tyk/bootstrap/pkg/constants"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExecutePreDeleteOperations executes operations needed in pre-delete chart hook one by one.
func (c *Client) ExecutePreDeleteOperations() error {
	if err := c.deleteOperatorSecret(); err != nil {
		return err
	}

	if err := c.deletePortalSecret(); err != nil {
		return err
	}

	if err := c.deleteBootstrappingJobs(); err != nil {
		return err
	}

	return nil
}

// deleteOperatorSecret deletes the Kubernetes secret created specifically for Tyk Operator.
func (c *Client) deleteOperatorSecret() error {
	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	found := false

	for i := range secrets.Items {
		value := secrets.Items[i]
		if value.Name == c.appArgs.OperatorKubernetesSecretName {
			err = c.clientSet.
				CoreV1().
				Secrets(c.appArgs.K8s.ReleaseNamespace).
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

// deletePortalSecret deletes the Kubernetes secret created specifically for Tyk Developer Portal.
func (c *Client) deletePortalSecret() error {
	fmt.Println("Running pre delete hook")

	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	notFound := true

	for i := range secrets.Items {
		value := secrets.Items[i]
		if c.appArgs.DevPortalKubernetesSecretName == value.Name {
			err = c.clientSet.CoreV1().Secrets(c.appArgs.K8s.ReleaseNamespace).
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

// deleteBootstrappingJobs deletes all jobs within the release namespace, that has specific label.
func (c *Client) deleteBootstrappingJobs() error {
	// Usually, the raw strings in label selectors are not recommended.
	jobs, err := c.clientSet.
		BatchV1().
		Jobs(c.appArgs.K8s.ReleaseNamespace).
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

	for i := range jobs.Items {
		job := jobs.Items[i]

		// Do not need to delete pre-delete job. It will be deleted by Helm.
		jobLabel := job.ObjectMeta.Labels[constants.TykBootstrapLabel]
		if jobLabel != constants.TykBootstrapPreDeleteLabel {
			deletePropagationType := metav1.DeletePropagationBackground

			err2 := c.clientSet.
				BatchV1().
				Jobs(c.appArgs.K8s.ReleaseNamespace).
				Delete(context.TODO(), job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePropagationType})
			if err2 != nil {
				errCascading = err2
			}
		}
	}

	return errCascading
}
