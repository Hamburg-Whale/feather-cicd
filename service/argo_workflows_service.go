package service

import (
	"context"
	"fmt"
	"log"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateArgoWorkflowsJob(dockerRepo string, imageName string, tag string, jdk string, builder string, name string, url string, namespace string) error {
	config, err := GetKubeConfig()
	if err != nil {
		log.Fatalf("failed to load in-cluster config: %v", err)
	}

	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(namespace)

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
						Args: []string{fmt.Sprintf(`
							apk add --no-cache git %s %s docker-cli &&
							git clone %s.git app &&
							cd app &&
							./gradlew build &&
							echo "$DOCKER_TOKEN" | docker login -u $DOCKER_USER --password-stdin &&
							docker build -t %s/%s:%s . &&
							docker push %s/%s:%s
						`, jdk, builder, url, dockerRepo, imageName, tag, dockerRepo, imageName, tag)},
						SecurityContext: &corev1.SecurityContext{
							Privileged: boolPtr(true),
						},
					},
				},

				/**
				apiVersion: argoproj.io/v1alpha1
				kind: Workflow
				metadata:
				generateName: sdk-demo-
				spec:
				entrypoint: build-and-push
				templates:
					- name: build-and-push
					  container:
						image: docker:24.0.5-dind
						command: ["sh", "-c"]
						args:
						- |
							apk add --no-cache git openjdk17 gradle docker-cli &&
							git clone https://github.com/example/springboot-app.git app &&
							cd app &&
							./gradlew build &&
							echo "$DOCKER_TOKEN" | docker login -u $DOCKER_USER --password-stdin &&
							docker build -t mydockeruser/springboot-app:v1 . &&
							docker push mydockeruser/springboot-app:v1
						securityContext:
						privileged: true
				**/
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
