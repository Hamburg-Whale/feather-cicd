package service

import (
	"context"
	"feather/types"
	"fmt"
	"log"
	"strings"
	"text/template"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateArgoWorkflowScript(req *types.CreateJobBasedJavaReq) (string, error) {
	workflowScript, err := buildCommand(req)
	if err != nil {
		return "", fmt.Errorf("failed to reate workflow script: %w", err)
	}
	return workflowScript, nil
}

func buildCommand(req *types.CreateJobBasedJavaReq) (string, error) {
	data := struct {
		JDK, BuildTool, URL, ImageRegistry, ImageName, ImageTag string
	}{
		JDK:           req.Jdk,
		BuildTool:     req.BuildTool,
		URL:           req.Url,
		ImageRegistry: req.ImageRegistry,
		ImageName:     req.ImageName,
		ImageTag:      req.ImageTag,
	}

	tmpl, err := template.ParseFiles("assets/templates/argo/ci.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template file: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

/*
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:

	generateName: $name
	namespace: $namespace

spec:

	entrypoint: build-and-push
	templates:
	- name: build-and-push
		container:
			image: docker:24.0.5-dind
			command: ["sh", "-c"]
			args: ~template~
			securityContext:
			  privileged: true
*/
func (service *Service) CreateArgoWorkflowsJobBasedSpringBoot(req *types.CreateJobBasedJavaReq) error {
	config, err := GetKubeConfig()
	if err != nil {
		log.Fatalf("failed to load in-cluster config: %v", err)
	}

	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(req.Namespace)

	command, err := buildCommand(req)

	workflow := &wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: req.Name + "-",
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: "build-and-push",
			Templates: []wfv1.Template{
				{
					Name: "build-and-push",
					Container: &corev1.Container{
						Image:   "docker:24.0.5-dind",
						Command: []string{"sh", "-c"},
						Args:    []string{command},
						SecurityContext: &corev1.SecurityContext{
							Privileged: boolPtr(true),
						},
					},
				},
			},
		},
	}

	result, err := wfClient.Create(context.TODO(), workflow, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	fmt.Println("Workflow submitted:", result.Name)
	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
