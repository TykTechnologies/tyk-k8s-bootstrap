package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BootstrapTykPortalSecret creates a secret required by Tyk Developer Portal pod which
// is not going to be ready until this secret is created. If there is a secret created already,
// it deletes the existing one and recreates the secret.
func (c *Client) BootstrapTykPortalSecret() error {
	err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		Delete(context.TODO(), c.appArgs.DevPortalKubernetesSecretName, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

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

	return nil
}
