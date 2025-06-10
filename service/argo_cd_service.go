package service

import (
	"feather/types"
	"fmt"
	"log"
)

func (s *Service) CreateProjectManifestRepo(projectId int64) error {
	const repoName = "feather-argocd"

	res, err := s.repository.ProjectWithBaseCampInfo(projectId)
	if err != nil {
		return fmt.Errorf("Get BaseCamp failed: %w", err)
	}

	checkReq := &types.CheckRepoRequest{
		URL:   res.BaseCampURL,
		Owner: res.BaseCampOwner,
		Name:  repoName,
		Token: res.Token,
	}

	exists, err := s.repoExists(checkReq)
	if err != nil {
		return fmt.Errorf("repository check failed: %w", err)
	}

	if !exists {
		createReq := &types.CreateRepoRequest{
			URL:         res.BaseCampURL,
			Description: "Repository for ArgoCD manifest management",
			Name:        repoName,
			Owner:       res.BaseCampOwner,
			Private:     false,
			Token:       res.Token,
		}

		if err := s.createRepo(createReq); err != nil {
			return fmt.Errorf("failed to create ArgoCD repository: %w", err)
		}
		log.Printf("Repository '%s' created successfully.", repoName)
	}

	return nil
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

func (s *Service) fileExists(req *types.CheckFileRequest) (bool, error) {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s", req.URL, req.Owner, req.Repo, req.FilePath)
	_, err := s.DoJSONGet(repoURL, req.Token)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) createFile(req *types.CreateFileRequest) error {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s", req.URL, req.Owner, req.Repo, req.FilePath)

	payload := map[string]interface{}{
		"Author":    req.Details.Author,
		"Branch":    req.Details.Branch,
		"NewBranch": req.Details.NewBranch,
		"Content":   req.Details.Content,
		"Message":   req.Details.Message,
	}

	_, err := s.DoJSONPost(repoURL, req.Token, payload)
	return err
}
