package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/controllers"
)

func functionalTestScalingInterface() {
	kubeClientLocal := controllers.NewKubeClientLocal()

	// deployment name to container name
	fmt.Printf("## Testing deployment to container conversion latency ##\n")
	timeStart := time.Now()
	response, err := kubeClientLocal.GetContainerNameFromDeployment(context.TODO(), "test-deployment", "test-namespace")
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("response: %v\ntime taken for response: %vms\n\n", response, time.Now().Sub(timeStart).Milliseconds())

	// get request value
	fmt.Printf("## Testing get request for cpu and mem latency ##\n")
	timeStart = time.Now()
	cpu, mem, err := kubeClientLocal.GetRequestValue(context.TODO(), "test-deployment", "test-namespace")
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("response: cpu->%v, mem->%v \ntime taken for response: %vms\n\n", cpu, mem, time.Now().Sub(timeStart).Milliseconds())

	// set request value
	fmt.Printf("## Testing set request for cpu and mem latency ##\n")
	timeStart = time.Now()
	cpu_float, mem_float := float32(300.00), float32(300.00)
	err = kubeClientLocal.SetRequestValue(context.TODO(), "test-deployment", "test-namespace", cpu_float, mem_float)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("time taken for response: %vms\n\n", time.Now().Sub(timeStart).Milliseconds())

	// set limit value
	fmt.Printf("## Testing set limit for cpu and mem latency ##\n")
	timeStart = time.Now()
	cpu_float, mem_float = float32(500.00), float32(500.00)
	err = kubeClientLocal.SetLimitValue(context.TODO(), "test-deployment", "test-namespace", cpu_float, mem_float)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("time taken for response: %vms\n\n", time.Now().Sub(timeStart).Milliseconds())

	// set replica value
	fmt.Printf("## Testing set replica latency ##\n")
	timeStart = time.Now()
	replica_value := int32(5)
	err = kubeClientLocal.SetReplicaValue(context.TODO(), "test-deployment", "test-namespace", replica_value)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("time taken for response: %vms\n\n", time.Now().Sub(timeStart).Milliseconds())
}

func main() {
	functionalTestScalingInterface()
}