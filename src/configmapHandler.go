package src

import (
	"context"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *K8sClient) HandleDelConfigMap(config *apiv1.ConfigMap, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	client := s.CoreV1().ConfigMaps(namespace)
	if err := client.Delete(context.Background(), config.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}


func (s *K8sClient) HandleCreateConfigMap(config *apiv1.ConfigMap, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	client := s.CoreV1().ConfigMaps(namespace)
	if _, err := client.Create(context.Background(), config, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}