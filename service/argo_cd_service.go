package service

import (
	"feather/types"
	"fmt"
	"log"
)

func (s *Service) CreateProjectManifestRepo(req *types.CreateRepoRequest) error {
	const repoName = "feather-argocd"

	checkReq := &types.CheckRepoRequest{
		URL:   req.URL,
		Owner: req.Owner,
		Name:  repoName,
		Token: req.Token,
	}

	exists, err := s.repoExists(checkReq)
	if err != nil {
		return fmt.Errorf("repository check failed: %w", err)
	}

	if !exists {
		createReq := &types.CreateRepoRequest{
			URL:         req.URL,
			Description: "ArgoCD 매니페스트 관리용 리포지토리",
			Name:        repoName,
			Owner:       req.Owner,
			Private:     false,
			Token:       req.Token,
		}

		if err := s.createRepo(createReq); err != nil {
			return fmt.Errorf("failed to create ArgoCD repository: %w", err)
		}
		log.Printf("Repository '%s' created successfully.", repoName)
	}

	return nil
}

func (s *Service) fileExists(req *types.CheckFileRequest) (bool, error) {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s", req.URL, req.Owner, req.Repo, req.FilePath)
	_, err := s.DoJSONGet(repoURL, req.Token)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) repoExists(req *types.CheckRepoRequest) (bool, error) {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s", req.URL, req.Owner, req.Name)
	_, err := s.DoJSONGet(repoURL, req.Token)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) createRepo(req *types.CreateRepoRequest) error {
	repoURL := fmt.Sprintf("%s/api/v1/user/repos", req.URL)
	payload := map[string]interface{}{
		"Description": req.Description,
		"Name":        req.Name,
		"Private":     req.Private,
	}
	_, err := s.DoJSONPost(repoURL, req.Token, payload)
	return err
}
