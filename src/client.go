package src

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type YamlItem struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	MetaData   interface{} `json:"metadata"`
}

type K8sClient struct {
	*kubernetes.Clientset
	Config *restclient.Config
}

func NewK8sClient(ip string, configPath ...string) (*K8sClient, error) {
	// if no args means config path is /root/.kube/config
	var path = os.Getenv("Home") + "/.kube/config"
	if len(configPath) == 0 {
		path = configPath[0]
	}
	config, err := clientcmd.BuildConfigFromFlags(ip, path)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{
		Clientset: client,
		Config:    config,
	}, nil
}

func (s *K8sClient) HandleAppConfig(path string) ([]interface{}, error) {
	if s.RESTClient() == nil {
		return nil, errors.New("RestClient is nil")
	}
	configs := make([]interface{}, 0)
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	items := strings.Split(string(fileBytes), "---")
	for _, item := range items {
		if item == "" {
			continue
		}
		var y YamlItem
		doc := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(item)), 100000)
		if err = doc.Decode(&y); err != nil {
			return nil, err
		}
		switch y.Kind {
		case StatefulSet:
			statefulSetSchema := &appsv1.StatefulSet{}
			if err := doc.Decode(&statefulSetSchema); err != nil {
				return nil, err
			}
			configs = append(configs, statefulSetSchema)
		case Service:
			serviceSchema := &apiv1.Service{}
			if err := doc.Decode(&serviceSchema); err != nil {
				return nil, err
			}
			configs = append(configs, serviceSchema)
		case ConfigMap:
			configMapSchema := &apiv1.ConfigMap{}
			if err := doc.Decode(&configMapSchema); err != nil {
				return nil, err
			}
			configs = append(configs, configMapSchema)
		case Deployment:
			DeploymentSchema := &appsv1.Deployment{}
			if err := doc.Decode(&DeploymentSchema); err != nil {
				return nil, err
			}
			configs = append(configs, DeploymentSchema)
		}
	}
	return configs, nil
}

// HandleListPod kubectl get wide -o
func (s *K8sClient)HandleListPod() ([][]interface{}, error) {
	var err error
	podList := make([][]interface{}, 0)
	podList = append(podList, []interface{}{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE", "NOMINATED NODE", "READINESS GATES"})
	if s.RESTClient() == nil {
		return nil, errors.New("RestClient is nil")
	}
	pods, err := s.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		var readinessGates  interface{}
		nominatedNodeName := ""
		restartCount := 0
		readyCount := 0
		for _, status := range pod.Status.ContainerStatuses {
			restartCount += int(status.RestartCount)
			if status.Ready {
				//fmt.Println(count.State.Running.StartedAt.UnixNano())
				readyCount ++
			}
		}
		if pod.Spec.ReadinessGates == nil {
			readinessGates = "<none>"
		}else {
			readinessGates = pod.Spec.ReadinessGates
		}
		if pod.Status.NominatedNodeName == "" {
			pod.Status.NominatedNodeName = "<none>"
		}
		if pod.Status.PodIP == "" {
			pod.Status.PodIP = "<node>"
		}
		if pod.Spec.NodeName == "" {
			pod.Spec.NodeName = "<node>"
		}
		var startTime int64
		if pod.Status.StartTime == nil {
			startTime = pod.CreationTimestamp.Unix()
		}else {
			startTime = pod.Status.StartTime.Unix()
		}
		podList = append(podList, []interface{}{
			pod.Name, fmt.Sprintf("%d/%d", readyCount, len(pod.Status.ContainerStatuses)),
			string(pod.Status.Phase), strconv.Itoa(restartCount), handleTimeTrans(time.Now().Unix()-
				startTime), pod.Status.PodIP, pod.Spec.NodeName, readinessGates,
			nominatedNodeName,
		})
	}
	return podList, nil
}


func handleTimeTrans(timeDiff int64) string {
	if timeDiff < 60 {
		return fmt.Sprintf("%ds", timeDiff)
	}
	if timeDiff < 3600 {
		return fmt.Sprintf("%dm%ds", timeDiff / 60, 320 % 60)
	}
	if timeDiff < 3600 * 24 {
		return fmt.Sprintf("%dh%dm", timeDiff / 3600, (timeDiff % 3600) / 60)
	}
	return fmt.Sprintf("%dd%dh", timeDiff / (3600 * 24), (timeDiff % (3600 * 24)) / 3600)
}

