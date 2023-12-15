package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/constants"
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

func functionalTestingPrometheusInterface() {
	// Creating the connection to the prometheus running in the localhost
	newProm := controllers.NewPrometheusInstance("http://localhost:9090")

	timeStart := time.Now()
	// create the promquery for memory
	promQueryMemory := controllers.BuildPromQueryForMemory("default", "1m", "prometheus")

	// get the data
	promResponseMemory, err := newProm.GetPrometheusData(context.TODO(), promQueryMemory, "memory")
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("The memory values are: %v", *promResponseMemory)
	timeEnd := time.Now()

	timeTaken := timeEnd.Sub(timeStart).Milliseconds()
	fmt.Printf("Time taken to get memory data from prometheus: %vms", timeTaken)

	fmt.Printf("\n\n")

	timeStart = time.Now()
	// create the promquery for cpu
	promQueryCPU := controllers.BuildPromQueryForCPU("default", "1m", "prometheus")

	// get the data
	promResponseCPU, err := newProm.GetPrometheusData(context.TODO(), promQueryCPU, "cpu")
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("The cpu values are: %v", *promResponseCPU)
	timeEnd = time.Now()

	timeTaken = timeEnd.Sub(timeStart).Milliseconds()
	fmt.Printf("Time taken to get cpu data from prometheus: %vms", timeTaken)
}

func testError() {
	listener := controllers.NewPrometheusInstance("http://localhost:9090")
	ns := "default"
	container := "prometheus"
	currTime := time.Now()
	startTime := currTime.AddDate(0, 0, constants.TrainingDataDuration)
	promQueryCPU := controllers.BuildPromQueryForCPU(ns, "2m", container)
	promQueryMem := controllers.BuildPromQueryForMemory(ns, "2m", container)
	// fmt.Printf("%s\n%s\n\n\n", promQueryCPU, promQueryMem)
	//fmt.Printf("%s\n\n\n", promQueryMem)
	_, _ = listener.GetPrometheusDataWithinRange(context.TODO(), promQueryMem, startTime, currTime, constants.StepsMinutesInterval*time.Minute, "memory")
	fmt.Printf("\n\n")
	_, _ = listener.GetPrometheusDataWithinRange(context.TODO(), promQueryCPU, startTime, currTime, constants.StepsMinutesInterval*time.Minute, "cpu")

	// cpuList := promCPU.PromItemList
	// for i, v := range promMem.PromItemList {
	// 	fmt.Printf("%v, %v\n", v, cpuList[i])
	// }
	//fmt.Print(promMem.PromDataType)
	//for _, v := range promMem.PromItemList {
	//	fmt.Print(v, "\n")
	//}
}

func testLocalLambdaExecution() {
	cmd := exec.Command("python", "./lambda/create_lambda/lambda_function.py", "prometheus-default")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", output)
}

func main() {
	// functionalTestScalingInterface()
	// functionalTestingPrometheusInterface()
	// testError()
	testLocalLambdaExecution()
}
