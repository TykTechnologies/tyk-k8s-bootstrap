package k8s

import (
	"context"
	"tyk/tyk/bootstrap/pkg/constants"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExecutePreDeleteOperations executes operations needed in pre-delete chart hook one by one.
func (c *Client) ExecutePreDeleteOperations() error {
	c.l.Info("Running pre delete hook")

	if err := c.deleteSecret(c.appArgs.OperatorKubernetesSecretName, true); err != nil {
		return err
	}

	if err := c.deleteSecret(c.appArgs.DevPortalKubernetesSecretName, true); err != nil {
		return err
	}

	if err := c.deleteBootstrappingJobs(); err != nil {
		return err
	}

	c.l.Info("Run pre delete hook successfully")

	return nil
}

// deleteSecret deletes the Kubernetes secret with given name.
func (c *Client) deleteSecret(secretName string, ignoreNotFound bool) error {
	err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) && ignoreNotFound {
			return nil
		}

		return err
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
