package src

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *K8sClient) HandleCreateService(config *apiv1.Service, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	serviceClient := s.CoreV1().Services(namespace)
	if _, err := serviceClient.Create(context.Background(), config, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}


func (s *K8sClient) HandleDelService(config *appsv1.StatefulSet, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	serviceClient := s.CoreV1().Services(namespace)
	if err := serviceClient.Delete(context.Background(), config.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
