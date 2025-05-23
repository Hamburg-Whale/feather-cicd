package service

import (
	"bytes"
	"encoding/json"
	"feather/types"
	"fmt"
	"io"
	"net/http"
)

func (service *Service) CreateRepo(req *types.CreateRepoReq) (*types.Response, error) {
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

	res, err := doJSONPost(repoURL, req.Token, payload)
	if err != nil {
		return nil, fmt.Errorf("리포지토리 생성 실패: %w", err)
	}
	defer res.Body.Close()

	var result types.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if req.WebhookFlag {
		err := service.createWebhook(req.Url, req.Token, req.Owner, req.Name, &req.CreateWebhookReq)
		if err != nil {
			return nil, fmt.Errorf("웹훅 생성 실패: %w", err)
		}
	}

	return &result, nil
}

func (service *Service) createWebhook(baseURL, token, owner, repo string, req *types.CreateWebhookReq) error {
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

	_, err := doJSONPost(hookURL, token, payload)
	if err != nil {
		return fmt.Errorf("웹훅 요청 실패: %w", err)
	}
	return nil
}

func doJSONPost(url, token string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("JSON 직렬화 실패: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("요청 실패: %w", err)
	}

	if res.StatusCode >= 300 {
		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("요청 실패: %s", string(bodyBytes))
	}

	return res, nil
}
