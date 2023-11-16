package controllers

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ScalingFunctions interface {
	GetRequestValue(ctx context.Context, deployment_name string, namespace_name string) (string, string, error)
	SetReplicaValue(ctx context.Context, deployment_name string, namespace_name string, replica_value int32) error
	SetRequestValue(ctx context.Context, deployment_name string, namespace_name string, cpu_request_value float32, memory_request_value float32) error
	SetLimitValue(ctx context.Context, deployment_name string, namespace_name string, cpu_limit_value float32, memory_limit_value float32) error
	GetContainerNameFromDeployment(ctx context.Context, deployment_name string, namespace_name string) (string, error)
}

type KubeClient struct {
	kubeClientset *kubernetes.Clientset
}

func NewKubeClientLocal() *KubeClient {
	// Use the current context in kubeconfig
	clientConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil
	}

	kubeclient := KubeClient{
		kubeClientset: clientset,
	}
	return &kubeclient
}

func NewKubeClient() *KubeClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error loading in-cluster config: %v\n", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error building kubernetes client: %v\n", err)
	}
	kubeclient := KubeClient{
		kubeClientset: clientset,
	}
	return &kubeclient
}

func (kc *KubeClient) GetContainerNameFromDeployment(ctx context.Context, deployment_name string, namespace_name string) (string, error) {
	deploymentsClient := kc.kubeClientset.AppsV1().Deployments(namespace_name)

	deployment, err := deploymentsClient.Get(ctx, deployment_name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("Error getting deployment %s: %v", deployment_name, err)
	}

	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return "", fmt.Errorf("No containers found in the deployment %s", deployment_name)
	}

	return deployment.Spec.Template.Spec.Containers[0].Name, nil
}

func (kc *KubeClient) SetReplicaValue(ctx context.Context, deployment_name string, namespace_name string, replica_value int32) error {
	deploymentsClient := kc.kubeClientset.AppsV1().Deployments(namespace_name)

	deployment, err := deploymentsClient.Get(ctx, deployment_name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting deployment %s: %v", deployment_name, err)
	}

	deployment.Spec.Replicas = &replica_value
	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Error updating deployment %s: %v", deployment_name, err)
	}

	return nil
}

func (kc *KubeClient) SetRequestValue(ctx context.Context, deployment_name string, namespace_name string, cpu_request_value float32, memory_request_value float32) error {
	deploymentsClient := kc.kubeClientset.AppsV1().Deployments(namespace_name)

	deployment, err := deploymentsClient.Get(ctx, deployment_name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting deployment %s: %v", deployment_name, err)
	}

	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("No containers found in the deployment %s", deployment_name)
	}

	// Update CPU and Memory resource requests for the first container in the pod template
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = v1.ResourceList{
		v1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu_request_value*1000), resource.DecimalSI),
		v1.ResourceMemory: *resource.NewQuantity(int64(memory_request_value), resource.DecimalSI),
	}

	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Error updating deployment %s: %v", deployment_name, err)
	}

	return nil
}

func (kc *KubeClient) SetLimitValue(ctx context.Context, deployment_name string, namespace_name string, cpu_limit_value float32, memory_limit_value float32) error {
	deploymentsClient := kc.kubeClientset.AppsV1().Deployments(namespace_name)

	deployment, err := deploymentsClient.Get(ctx, deployment_name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting deployment %s: %v", deployment_name, err)
	}

	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("No containers found in the deployment %s", deployment_name)
	}

	// Update CPU and Memory resource limits for the first container in the pod template
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = v1.ResourceList{
		v1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu_limit_value*1000), resource.DecimalSI),
		v1.ResourceMemory: *resource.NewQuantity(int64(memory_limit_value), resource.DecimalSI),
	}

	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Error updating deployment %s: %v", deployment_name, err)
	}

	return nil
}

func (kc *KubeClient) GetRequestValue(ctx context.Context, deployment_name string, namespace_name string) (string, string, error) {
	deploymentsClient := kc.kubeClientset.AppsV1().Deployments(namespace_name)

	deployment, err := deploymentsClient.Get(ctx, deployment_name, metav1.GetOptions{})
	if err != nil {
		return "", "", fmt.Errorf("Error getting deployment %s: %v", deployment_name, err)
	}

	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return "", "", fmt.Errorf("No containers found in the deployment %s", deployment_name)
	}

	cpuRequest := deployment.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU]
	memoryRequest := deployment.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory]

	return cpuRequest.String(), memoryRequest.String(), nil
}
