package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tykModePro = "pro"
	tykModeKey = "TYK_MODE"
	tykAuthKey = "TYK_AUTH"
	tykOrgKey  = "TYK_ORG"
	tykURLKey  = "TYK_URL"
)

// BootstrapTykOperatorSecret bootstrap a Kubernetes Secret utilized by Tyk Operator.
// If the system has the secret created already, it deletes the existing one and recreates
// a secret for Tyk Operator.
func (c *Client) BootstrapTykOperatorSecret() error {
	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for i := range secrets.Items {
		secret := secrets.Items[i]
		if secret.Name == c.appArgs.OperatorKubernetesSecretName {
			err = c.clientSet.
				CoreV1().
				Secrets(c.appArgs.K8s.ReleaseNamespace).
				Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})

			if err != nil {
				return err
			}

			fmt.Println("A previously created operator secret was identified and deleted")

			break
		}
	}

	secretData := map[string][]byte{
		tykAuthKey: []byte(c.appArgs.Tyk.Admin.Auth),
		tykOrgKey:  []byte(c.appArgs.Tyk.Org.ID),
		tykModeKey: []byte(tykModePro),
		tykURLKey:  []byte(c.appArgs.K8s.DashboardSvcUrl),
	}

	objectMeta := metav1.ObjectMeta{Name: c.appArgs.OperatorKubernetesSecretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err = c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// BootstrapTykPortalSecret creates a secret required by Tyk Developer Portal pod which
// is not going to be ready until this secret is created. If there is a secret created already,
// it deletes the existing one and recreates the secret.
func (c *Client) BootstrapTykPortalSecret() error {
	// TODO(buraksekili): do we need to list the secrets? Can we just try getting the secret?
	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for i := range secrets.Items {
		secret := secrets.Items[i]

		if c.appArgs.DevPortalKubernetesSecretName == secret.Name {
			err = c.clientSet.
				CoreV1().
				Secrets(c.appArgs.K8s.ReleaseNamespace).
				Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println("A previously created portal secret was identified and deleted")

			break
		}
	}

	if c.appArgs.DevPortalKubernetesSecretName != "" {
		secretData := map[string][]byte{
			tykAuthKey: []byte(c.appArgs.Tyk.Admin.Auth),
			tykOrgKey:  []byte(c.appArgs.Tyk.Org.ID),
		}

		objectMeta := metav1.ObjectMeta{Name: c.appArgs.DevPortalKubernetesSecretName}

		secret := v1.Secret{
			ObjectMeta: objectMeta,
			Data:       secretData,
		}

		_, err = c.clientSet.
			CoreV1().
			Secrets(c.appArgs.K8s.ReleaseNamespace).
			Create(context.TODO(), &secret, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
