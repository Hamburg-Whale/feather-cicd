package service

import (
	"encoding/json"
	"feather/types"
	"fmt"
	"log"
)

type GitService struct {
}

func (s *Service) CreateRepo(req *types.RepoFromTemplateRequest) (*types.Response, error) {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/generate", req.URL, req.Template.Owner, req.Template.Repo)

	payload := map[string]interface{}{
		"avatar":           req.Options.Avatar,
		"default_branch":   req.Options.DefaultBranch,
		"description":      req.Options.Description,
		"git_content":      req.Options.GitContent,
		"git_hooks":        req.Options.GitHooks,
		"labels":           req.Options.Labels,
		"name":             req.Name,
		"owner":            req.Owner,
		"private":          req.Private,
		"protected_branch": req.Options.ProtectedBranch,
		"topics":           req.Options.Topics,
		"webhooks":         req.Options.Webhooks,
	}

	res, err := s.DoJSONPost(repoURL, req.Token, payload)
	if err != nil {
		return nil, fmt.Errorf("repository creation failed: %w", err)
	}
	defer res.Body.Close()

	var result types.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if req.WebhookEnabled && req.Webhook != nil {
		if err := s.attachWebhook(req); err != nil {
			return nil, err
		}
	}

	log.Printf("Repository created: %s/%s", req.Owner, req.Name)
	return &result, nil
}

func (s *Service) attachWebhook(req *types.RepoFromTemplateRequest) error {
	hookURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", req.URL, req.Owner, req.Name)

	hookType := req.Webhook.Type
	if hookType == "" {
		hookType = "gitea"
	}

	payload := map[string]interface{}{
		"type": hookType,
		"config": map[string]string{
			"url":          req.Webhook.URL,
			"content_type": "json",
		},
		"events":        []string{"push"},
		"branch_filter": req.Webhook.BranchFilter,
		"active":        true,
	}

	if _, err := s.DoJSONPost(hookURL, req.Token, payload); err != nil {
		return fmt.Errorf("webhook creation failed: %w", err)
	}
	log.Printf("Webhook created for: %s/%s", req.Owner, req.Name)
	return nil
}
