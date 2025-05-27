package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"

	"feather/types"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func (service *Service) CreateArgoSensor(req *types.CreateJobBasedJavaReq) error {
	config, err := GetKubeConfig()
	if err != nil {
		log.Fatalf("failed to load in-cluster config: %v", err)
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("dynamic client 생성 실패: %v", err)
	}

	workflowScript, err := CreateArgoWorkflowScript(req)
	if err != nil {
		log.Fatalf("Workflow Script 생성 실패 : %v", err)
	}

	sensorParams := struct {
		Namespace          string
		SensorName         string
		ServiceAccountName string
		EventSourceName    string
		EventName          string
		TriggerName        string
		WorkflowNamePrefix string
		WorkflowName       string
		WorkflowScript     string
		DockerUser         string
	}{
		Namespace:          req.Namespace,
		SensorName:         req.Name + "-sensor",
		ServiceAccountName: "gitea-sensor-sa",
		EventSourceName:    "gitea-webhook",
		EventName:          "gitea",
		TriggerName:        "argo-workflow-trigger",
		WorkflowNamePrefix: req.Name + "-",
		WorkflowName:       "build-and-push",
		WorkflowScript:     workflowScript,
		DockerUser:         "DockerUser",
	}

	tmpl, err := template.New("sensor.tmpl").Funcs(sprig.TxtFuncMap()).ParseFiles("./sensor.tmpl")
	if err != nil {
		return fmt.Errorf("sensor.tmpl 파일 읽기 실패: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, sensorParams)
	if err != nil {
		return fmt.Errorf("템플릿 실행 실패: %w", err)
	}

	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(buf.Bytes(), &obj); err != nil {
		return fmt.Errorf("Sensor YAML 디코딩 실패: %w", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "sensors",
	}

	res, err := dyn.Resource(gvr).Namespace(sensorParams.Namespace).Create(context.Background(), &obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Sensor 생성 실패: %w", err)
	}

	log.Printf("Sensor 생성 완료: %s\n", res.GetName())
	return nil
}
