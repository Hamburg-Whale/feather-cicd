package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildCommand(jdk, builder, url, dockerRepo, imageName, tag string) (string, error) {
	data := struct {
		JDK, Builder, URL, DockerRepo, ImageName, Tag string
	}{
		JDK:        jdk,
		Builder:    builder,
		URL:        url,
		DockerRepo: dockerRepo,
		ImageName:  imageName,
		Tag:        tag,
	}

	tmpl, err := template.ParseFiles("build.tmpl")
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
func CreateArgoWorkflowsJobBasedSpringBoot(name, namespace, jdk, builder, url, dockerRepo, imageName, tag string) error {
	config, err := GetKubeConfig()
	if err != nil {
		log.Fatalf("failed to load in-cluster config: %v", err)
	}

	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(namespace)

	command, err := buildCommand(jdk, builder, url, dockerRepo, imageName, tag)

	workflow := &wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: name + "-",
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
