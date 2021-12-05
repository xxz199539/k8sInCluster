package src

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *K8sClient) HandleDelDeployment(config *appsv1.Deployment, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	client := s.AppsV1().Deployments(namespace)
	if err := client.Delete(context.Background(), config.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *K8sClient) HandleCreateDeployment(config *appsv1.Deployment, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	client := s.AppsV1().Deployments(namespace)
	if err := client.Delete(context.Background(), config.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}