package k8s

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tyk/tyk/bootstrap/tyk"
)

func (c *Client) BootstrapTykOperatorSecret() error {
	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.AppArgs.ReleaseNamespace).
		List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	for _, value := range secrets.Items {
		if value.Name == c.AppArgs.OperatorSecretName {
			err = c.clientSet.
				CoreV1().
				Secrets(c.AppArgs.ReleaseNamespace).
				Delete(context.TODO(), value.Name, metaV1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println("A previously created operator secret was identified and deleted")
			break
		}
	}

	err = c.createTykOperatorSecret()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createTykOperatorSecret() error {
	secretData := map[string][]byte{
		tyk.TykAuth: []byte(c.AppArgs.TykUserAuth),
		tyk.TykOrg:  []byte(c.AppArgs.TykOrgId),
		tyk.TykMode: []byte(tyk.TykModePro),
		tyk.TykUrl:  []byte(c.AppArgs.DashboardSvcAddr),
	}

	objectMeta := metaV1.ObjectMeta{Name: c.AppArgs.OperatorSecretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err := c.clientSet.
		CoreV1().
		Secrets(c.AppArgs.ReleaseNamespace).
		Create(context.TODO(), &secret, metaV1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) BootstrapTykEnterprisePortalSecret() error {
	secrets, err := c.clientSet.
		CoreV1().
		Secrets(c.AppArgs.ReleaseNamespace).
		List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	for _, value := range secrets.Items {
		if c.AppArgs.EnterprisePortalSecretName == value.Name {
			err = c.clientSet.
				CoreV1().
				Secrets(c.AppArgs.ReleaseNamespace).
				Delete(context.TODO(), value.Name, metaV1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println("A previously created enterprise portal secret was identified and deleted")
			break
		}
	}

	err = c.createTykEnterprisePortalSecret()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createTykEnterprisePortalSecret() error {
	secretData := map[string][]byte{
		tyk.TykAuth: []byte(c.AppArgs.TykUserAuth),
		tyk.TykOrg:  []byte(c.AppArgs.TykOrgId),
	}

	objectMeta := metaV1.ObjectMeta{Name: c.AppArgs.EnterprisePortalSecretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err := c.clientSet.
		CoreV1().
		Secrets(c.AppArgs.ReleaseNamespace).
		Create(context.TODO(), &secret, metaV1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
