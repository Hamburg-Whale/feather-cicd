package types

type CreateUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Gitea Auth with Token
// https://demo.gitea.com/api/swagger#/user/userGetCurrent
type AuthUserReq struct {
	Url      string `json:"url" binding:"required"`
	Username string `json:"username" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// Gitea Create a repository using a template
// https://docs.gitea.com/en-us/api-reference/repository/generate-repo
/*
GenerateRepoOption{
	description:
	GenerateRepoOption options when creating repository using a template

	avatar	[...]
	default_branch	[...]
	description	[...]
	git_content	[...]
	git_hooks	[...]
	labels	[...]
	name*	[...]
	owner*	[...]
	private	[...]
	protected_branch	[...]
	topics	[...]
	webhooks	[...]
}
*/

type RepoFromTemplateRequest struct {
	// General Info
	URL     string `json:"url" binding:"required"`
	Token   string `json:"token" binding:"required"`
	Owner   string `json:"owner" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Private bool   `json:"private,omitempty"`

	// Template Info
	Template TemplateInfo `json:"template" binding:"required"`

	// Optional Features
	Options TemplateRepoOptions `json:"options,omitempty"`

	// Webhook
	WebhookEnabled bool           `json:"webhook_enabled"`
	Webhook        *WebhookConfig `json:"webhook,omitempty"`
}

type TemplateInfo struct {
	Owner string `json:"owner" binding:"required"`
	Repo  string `json:"repo" binding:"required"`
}

type TemplateRepoOptions struct {
	Avatar          bool   `json:"avatar,omitempty"`
	DefaultBranch   string `json:"default_branch,omitempty"`
	Description     string `json:"description,omitempty"`
	GitContent      bool   `json:"git_content,omitempty"`
	GitHooks        bool   `json:"git_hooks,omitempty"`
	Labels          bool   `json:"labels,omitempty"`
	ProtectedBranch bool   `json:"protected_branch,omitempty"`
	Topics          bool   `json:"topics,omitempty"`
	Webhooks        bool   `json:"webhooks,omitempty"`
}

type WebhookConfig struct {
	Type         string `json:"type,omitempty"` // default is "gitea"
	BranchFilter string `json:"branch_filter" binding:"required"`
	URL          string `json:"url" binding:"required"`
}

type JobBasedJavaRequest struct {
	Name          string `json:"name" binding:"required"`
	Namespace     string `json:"namespace" binding:"required"`
	JDK           string `json:"jdk" binding:"required"`
	BuildTool     string `json:"build_tool" binding:"required"`
	URL           string `json:"url" binding:"required"`
	ImageRegistry string `json:"image_registry" binding:"required"`
	ImageName     string `json:"image_name" binding:"required"`
	ImageTag      string `json:"image_tag" binding:"required"`
}

type CreateRepoRequest struct {
	URL         string `json:"url" binding:"required"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name" binding:"required"`
	Owner       string `json:"owner" binding:"required"`
	Private     bool   `json:"private,omitempty"`
	Token       string `json:"token" binding:"required"`
}

type CheckRepoRequest struct {
	URL   string `json:"url" binding:"required"`
	Token string `json:"token" binding:"required"`
	Owner string `json:"owner" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

type CheckFileRequest struct {
	URL      string `json:"url" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
	Repo     string `json:"repo" binding:"required"`
	FilePath string `json:"file_path" binding:"required"`
}
