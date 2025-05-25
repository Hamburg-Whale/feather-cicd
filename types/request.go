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

type CreateRepoReq struct {
	WebhookFlag      bool             `json:"webhook_flag"`
	Url              string           `json:"url" binding:"required"`
	Token            string           `json:"token" binding:"required"`
	TemplateOwner    string           `json:"template_owner" binding:"required"`
	TemplateRepo     string           `json:"template_repo" binding:"required"`
	Avatar           bool             `json:"avatar"`
	DefaultBranch    string           `json:"default_branch"`
	Description      string           `json:"description"`
	GitContent       bool             `json:"git_content"`
	GitHooks         bool             `json:"git_hooks"`
	Labels           bool             `json:"labels"`
	Name             string           `json:"name"`
	Owner            string           `json:"owner"`
	Private          bool             `json:"private"`
	ProtectedBranch  bool             `json:"protected_branch"`
	Topics           bool             `json:"topics"`
	Webhooks         bool             `json:"webhooks"`
	CreateWebhookReq CreateWebhookReq `json:"create_webhook_req"`
}

type CreateWebhookReq struct {
	Type         string `json:"type"`
	BranchFilter string `json:"branch_filter" binding:"required"`
	Url          string `json:"url" binding:"required"`
}

type CreateJobBasedJavaReq struct {
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Jdk           string `json:"jdk"`
	BuildTool     string `json:"build_tool"`
	Url           string `json:"url"`
	ImageRegistry string `json:"image_registry"`
	ImageName     string `json:"image_name"`
	ImageTag      string `json:"image_tag"`
}
