# Environment
Programming Language : Go, JavaScript
Framework : Gin, React
IDE : Cursor, Goland
Dev Env : Kubernetes, Docker, Gitea, Argo Workflows, ArgoCD, Dive
Service Mesh : Istio
Observability : OpenTelemetry, Loki, Prometheus, Tempo, Kiali
AI : ChatGPT
# Summary
## Feather CI/CD Auto-Registration
### Seamlessly connect your GitHub repositories with Argo Workflows and ArgoCD for automated deployments
# Introduction

## Create Repository using a Template
## Gitea
Gitea API : Create a repository using a template
https://docs.gitea.com/en-us/api-reference/repository/generate-repo
### Gitea API : GenerateRepoOption

| avatar           | boolean<br><br>include avatar of the template repo                       |
| ---------------- | ------------------------------------------------------------------------ |
| default_branch   | string<br><br>Default branch of the new repository                       |
| description      | string<br><br>Description of the repository to create                    |
| git_content      | boolean<br><br>include git content of default branch in template repo    |
| git_hooks        | boolean<br><br>include git hooks in template repo                        |
| labels           | boolean<br><br>include labels in template repo                           |
| name*            | string  <br><br>Name of the repository to create                         |
| owner*           | string<br><br>The organization or person who will own the new repository |
| private          | boolean<br><br>Whether the repository is private                         |
| protected_branch | boolean<br><br>include protected branches in template repo               |
| topics           | boolean<br><br>include topics in template repo                           |
| webhooks         | boolean<br><br>include webhooks in template repo                         |
### Argo Workflows
