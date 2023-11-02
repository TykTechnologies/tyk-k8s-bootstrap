package helpers

import (
	"context"
	"fmt"
	"tyk/tyk/bootstrap/data"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func BootstrapTykOperatorSecret() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	secrets, err := clientset.
		CoreV1().
		Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for i := range secrets.Items {
		secret := secrets.Items[i]
		if secret.Name == data.BootstrapConf.OperatorKubernetesSecretName {
			err = clientset.
				CoreV1().
				Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
				Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})

			if err != nil {
				return err
			}

			fmt.Println("A previously created operator secret was identified and deleted")

			break
		}
	}

	err = CreateTykOperatorSecret(clientset)
	if err != nil {
		return err
	}

	return nil
}

func CreateTykOperatorSecret(clientset *kubernetes.Clientset) error {
	secretData := map[string][]byte{
		TykAuth: []byte(data.BootstrapConf.Tyk.UserAuth),
		TykOrg:  []byte(data.BootstrapConf.Tyk.OrgId),
		TykMode: []byte(TykModePro),
		TykUrl:  []byte(data.BootstrapConf.K8s.DashboardSvcUrl),
	}

	objectMeta := metav1.ObjectMeta{Name: data.BootstrapConf.OperatorKubernetesSecretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err := clientset.
		CoreV1().
		Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
		Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func BootstrapTykPortalSecret() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	secrets, err := clientset.CoreV1().Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for i := range secrets.Items {
		secret := secrets.Items[i]

		if data.BootstrapConf.DevPortalKubernetesSecretName == secret.Name {
			err = clientset.CoreV1().Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
				Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println("A previously created portal secret was identified and deleted")

			break
		}
	}

	if data.BootstrapConf.DevPortalKubernetesSecretName != "" {
		err = CreateTykPortalSecret(clientset, data.BootstrapConf.DevPortalKubernetesSecretName)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateTykPortalSecret(clientset *kubernetes.Clientset, secretName string) error {
	secretData := map[string][]byte{
		TykAuth: []byte(data.BootstrapConf.Tyk.UserAuth),
		TykOrg:  []byte(data.BootstrapConf.Tyk.OrgId),
	}

	objectMeta := metav1.ObjectMeta{Name: secretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err := clientset.
		CoreV1().
		Secrets(data.BootstrapConf.K8s.ReleaseNamespace).
		Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
