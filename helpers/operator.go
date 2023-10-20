package helpers

import (
	"context"
	"fmt"
	"github.com/TykTechnologies/tyk-k8s-bootstrap/data"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
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

	secrets, err := clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
		List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, value := range secrets.Items {
		if value.Name == data.AppConfig.OperatorSecretName {
			err = clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
				Delete(context.TODO(), value.Name, v1.DeleteOptions{})
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
		TykAuth: []byte(data.AppConfig.UserAuth),
		TykOrg:  []byte(data.AppConfig.OrgId),
		TykMode: []byte(TykModePro),
		TykUrl:  []byte(data.AppConfig.DashboardUrl),
	}

	objectMeta := v1.ObjectMeta{Name: data.AppConfig.OperatorSecretName}

	secret := v12.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}
	_, err := clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
		Create(context.TODO(), &secret, v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func BootstrapTykEnterprisePortalSecret() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	secrets, err := clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
		List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, value := range secrets.Items {
		if data.AppConfig.EnterprisePortalSecretName == value.Name {
			err = clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
				Delete(context.TODO(), value.Name, v1.DeleteOptions{})
			if err != nil {
				return err
			}
			fmt.Println("A previously created enterprise portal secret was identified and deleted")
			break
		}
	}

	err = CreateTykEnterprisePortalSecret(clientset)
	if err != nil {
		return err
	}

	return nil
}

func CreateTykEnterprisePortalSecret(clientset *kubernetes.Clientset) error {
	secretData := map[string][]byte{
		TykAuth: []byte(data.AppConfig.UserAuth),
		TykOrg:  []byte(data.AppConfig.OrgId),
	}

	objectMeta := v1.ObjectMeta{Name: data.AppConfig.EnterprisePortalSecretName}

	secret := v12.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}
	_, err := clientset.CoreV1().Secrets(data.AppConfig.TykPodNamespace).
		Create(context.TODO(), &secret, v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
