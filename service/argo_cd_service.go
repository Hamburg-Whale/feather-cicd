package service

import (
	"bytes"
	"feather/types"
	"fmt"
	"log"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func (s *Service) CreateProjectManifestRepo(projectId int64) error {
	const (
		repoName         = "feather-argocd"
		appSetFolderPath = "application-sets"
		appSetFileName   = "application-set.yaml"
		manifestFileName = "manifest.yaml"
	)

	applicationSetFilePath := fmt.Sprintf("%s/%s", appSetFolderPath, appSetFileName)

	res, err := s.repository.ProjectWithBaseCampInfo(projectId)
	if err != nil {
		return fmt.Errorf("Get BaseCamp failed: %w", err)
	}

	// projectManifestFilePath := fmt.Sprintf("%s/%s", res.ProjectName, manifestFileName)

	if err := s.ensureArgoCdRepo(res, repoName); err != nil {
		return err
	}

	if err := s.ensureApplicationSet(res, repoName, applicationSetFilePath); err != nil {
		return err
	}

	return nil
}

func (s *Service) ensureApplicationSet(res *types.ProjectWithBaseCampInfo, repoName string, filePath string) error {
	checkApplicationSetRepoReq := &types.CheckFileRequest{
		URL:      res.BaseCampURL,
		Token:    res.Token,
		Owner:    res.BaseCampOwner,
		Repo:     repoName,
		FilePath: filePath,
	}

	exists, err := s.fileExists(checkApplicationSetRepoReq)
	if err != nil {
		return fmt.Errorf("file check failed: %w", err)
	}

	if !exists {
		applicationSetName := fmt.Sprintf("%s-appset", res.BaseCampName)
		applicationSetURL := fmt.Sprintf("%s/%s.git", res.BaseCampURL, repoName)
		params := struct {
			ApplicationSetName string
			URL                string
		}{
			ApplicationSetName: applicationSetName,
			URL:                applicationSetURL,
		}

		tmpl, err := template.New("application-set.tmpl").Funcs(sprig.TxtFuncMap()).ParseFiles("assets/templates/argo/application-set.tmpl")
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, params); err != nil {
			return fmt.Errorf("failed to execute application-set template: %w", err)
		}

		generatedYAML := buf.String()

		author := &types.Author{
			Email: "feather@feather.com",
			Name:  "feather",
		}

		fileCommitDetails := &types.FileCommitDetails{
			Author:  *author,
			Content: generatedYAML,
			Message: "Create Application Set YAML",
		}

		createReq := &types.CreateFileRequest{
			URL:      res.BaseCampURL,
			Token:    res.Token,
			Owner:    res.BaseCampOwner,
			Repo:     repoName,
			FilePath: filePath,
			Details:  *fileCommitDetails,
		}

		if err := s.createFile(createReq); err != nil {
			return fmt.Errorf("failed to create application set file: %w", err)
		}
	}

	return nil
}

func (s *Service) ensureArgoCdRepo(res *types.ProjectWithBaseCampInfo, repoName string) error {

	checkArgoCdRepoReq := &types.CheckRepoRequest{
		URL:   res.BaseCampURL,
		Token: res.Token,
		Owner: res.BaseCampOwner,
		Name:  repoName,
	}

	exists, err := s.repoExists(checkArgoCdRepoReq)
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
	} else {
		log.Printf("Repository '%s' already exists.", repoName)
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
