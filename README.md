# 🚀 Feather CI/CD Auto-Registration

> **Seamlessly connect your GitHub repositories with Argo Workflows and ArgoCD for automated deployments**

![build](https://github.com/user-attachments/assets/ce2e1094-c89e-4c49-8a1c-1c16741dedee)

---

## 🌐 Environment

| Category        | Tech Stack                                                                 |
|----------------|------------------------------------------------------------------------------|
| **Language**    | Go, JavaScript                                                              |
| **Framework**   | Gin (Go), React (JS)                                                        |
| **IDE**         | Cursor, GoLand                                                              |
| **Dev Env**     | Kubernetes, Docker, Gitea, Argo Workflows, ArgoCD, Dive                    |
| **Service Mesh**| Istio                                                                       |
| **Observability**| OpenTelemetry, Loki, Prometheus, Tempo, Kiali                             |
| **AI**          | ChatGPT                                                                     |

---

## 🧩 Summary

Feather streamlines the CI/CD pipeline setup by **auto-registering repositories** created via Gitea templates into **Argo Workflows** and **ArgoCD**, enabling **hands-free, automated deployments** from the moment a project is scaffolded.

---

## 🛠️ Features

- 🔨 **Create repositories from templates** via the Gitea API
- ⚙️ **Auto-configure Argo Workflows** for CI pipelines
- 🚀 **Auto-sync with ArgoCD** for GitOps-based CD
- 🔍 **Fully observable** with Prometheus, Tempo, Loki, and Kiali

---

## 📦 Create Repository Using a Template

Feather leverages the [Gitea GenerateRepo API](https://docs.gitea.com/en-us/api-reference/repository/generate-repo) to scaffold new repositories from templates with customizable options.

### 📑 Gitea API: `GenerateRepoOption`

| Field            | Type      | Description                                                                  |
|------------------|-----------|------------------------------------------------------------------------------|
| `name`*          | `string`  | Name of the repository to create                                             |
| `owner`*         | `string`  | The organization or person who will own the new repository                  |
| `description`    | `string`  | Description of the new repository                                            |
| `default_branch` | `string`  | Default branch name                                                          |
| `private`        | `boolean` | Whether the repository is private                                            |
| `avatar`         | `boolean` | Include avatar of the template repo                                          |
| `git_content`    | `boolean` | Include Git content of default branch in template repo                       |
| `git_hooks`      | `boolean` | Include Git hooks from the template repo                                     |
| `labels`         | `boolean` | Include issue labels from the template repo                                  |
| `protected_branch`| `boolean`| Include protected branches from the template repo                            |
| `topics`         | `boolean` | Include topic tags                                                           |
| `webhooks`       | `boolean` | Include webhooks from the template repo                                      |

---

## ⚙️ Argo Workflows

Feather registers newly created repositories to **Argo Workflows**, enabling CI pipelines to be executed automatically. This includes:

- 🧪 Test automation
- 🛠️ Build steps
- ✅ Status reporting

---

## 🔁 ArgoCD

Repositories are also linked with **ArgoCD**, providing:

- 📦 Continuous Deployment through GitOps
- 🔄 Automatic sync of application manifests
- 🔒 Declarative security and rollback support

---

## 📸 Preview

Coming soon...

---

## 🤝 Contributing

Contributions are welcome! Please open issues or submit PRs to improve Feather or its integrations.

---

## 📄 License

[MIT License](./LICENSE)

---

*Made with ❤️ by developers for developers*
