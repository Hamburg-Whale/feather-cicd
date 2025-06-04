package service

import (
	"encoding/json"
	"feather/types"
	"fmt"
)

func (service *Service) CreateRepo(req *types.CreateRepoBasedTemplateReq) (*types.Response, error) {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/generate", req.Url, req.TemplateOwner, req.TemplateRepo)

	payload := map[string]interface{}{
		"avatar":           req.Avatar,
		"default_branch":   req.DefaultBranch,
		"description":      req.Description,
		"git_content":      req.GitContent,
		"git_hooks":        req.GitHooks,
		"labels":           req.Labels,
		"name":             req.Name,
		"owner":            req.Owner,
		"private":          req.Private,
		"protected_branch": req.ProtectedBranch,
		"topics":           req.Topics,
		"webhooks":         req.Webhooks,
	}

	res, err := DoJSONPost(repoURL, req.Token, payload)
	if err != nil {
		return nil, fmt.Errorf("리포지토리 생성 실패: %w", err)
	}
	defer res.Body.Close()

	var result types.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if req.WebhookFlag {
		err := createWebhook(req.Url, req.Token, req.Owner, req.Name, &req.CreateWebhookReq)
		if err != nil {
			return nil, fmt.Errorf("웹훅 생성 실패: %w", err)
		}
	}

	return &result, nil
}

func createWebhook(baseURL, token, owner, repo string, req *types.CreateWebhookReq) error {
	hookURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", baseURL, owner, repo)

	hookType := req.Type
	if hookType == "" {
		hookType = "gitea"
	}

	payload := map[string]interface{}{
		"type": hookType,
		"config": map[string]string{
			"url":          req.Url,
			"content_type": "json",
		},
		"events":        []string{"push"},
		"branch_filter": req.BranchFilter,
		"active":        true,
	}

	_, err := DoJSONPost(hookURL, token, payload)
	if err != nil {
		return fmt.Errorf("웹훅 요청 실패: %w", err)
	}
	return nil
}
