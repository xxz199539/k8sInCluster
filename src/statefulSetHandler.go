package src

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

func (s *K8sClient) HandleCreateStatefulSet(config *appsv1.StatefulSet, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	fulSetClient := s.AppsV1().StatefulSets(namespace)
	if _, err := fulSetClient.Create(context.Background(), config, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}


func (s *K8sClient) HandleDelStatefulSet(config *appsv1.StatefulSet, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}
	fulSetClient := s.AppsV1().StatefulSets(namespace)
	if err := fulSetClient.Delete(context.Background(), config.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
