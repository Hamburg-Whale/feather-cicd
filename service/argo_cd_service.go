package service

import (
	"bytes"
	"encoding/base64"
	"feather/repository"
	"feather/types"
	"fmt"
	"log"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type ArgoCdService interface {
	CreateProjectManifestRepo(projectId int64) error
	ensureApplicationSet(res *types.ProjectWithBaseCampInfo, repoName string, filePath string) error
	ensureArgoCdRepo(res *types.ProjectWithBaseCampInfo, repoName string) error
}

type argoCdServiceImpl struct {
	repository *repository.Repository
	gitService GitService
}

func NewArgoCdService(repository *repository.Repository, gitService GitService) ArgoCdService {
	return &argoCdServiceImpl{
		repository: repository,
		gitService: gitService,
	}
}

func (s *argoCdServiceImpl) CreateProjectManifestRepo(projectId int64) error {
	const (
		repoName         = "feather-argocd"
		appSetFolderPath = "application-sets"
		appSetFileName   = "application-set.yaml"
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

func (s *argoCdServiceImpl) ensureProjectManifest(res *types.ProjectWithBaseCampInfo, repoName string) error {
	const manifestFileName = "manifest.yaml"
	filePath := fmt.Sprintf("%s/%s/%s", repoName, res.ProjectName, manifestFileName)

	checkProjectManifestReq := &types.CheckFileRequest{
		URL:      res.BaseCampURL,
		Token:    res.Token,
		Owner:    res.BaseCampOwner,
		Repo:     repoName,
		FilePath: filePath,
	}

	exists, err := s.gitService.FileExists(checkProjectManifestReq)
	if err != nil {
		return fmt.Errorf("file check failed: %w", err)
	}
	log.Print("File Check Complete \n")

	if exists {
		log.Printf("Project Manifest already exists at %s", filePath)
		return nil
	}
}

func (s *argoCdServiceImpl) ensureApplicationSet(res *types.ProjectWithBaseCampInfo, repoName string, filePath string) error {
	checkApplicationSetDirReq := &types.CheckFileRequest{
		URL:      res.BaseCampURL,
		Token:    res.Token,
		Owner:    res.BaseCampOwner,
		Repo:     repoName,
		FilePath: filePath,
	}

	exists, err := s.gitService.FileExists(checkApplicationSetDirReq)
	if err != nil {
		return fmt.Errorf("file check failed: %w", err)
	}
	log.Print("File Check Complete \n")

	if exists {
		log.Printf("ApplicationSet file already exists at %s", filePath)
		return nil
	}

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
	encodedYaml := base64.StdEncoding.EncodeToString([]byte(generatedYAML))

	author := &types.Author{
		Email: "feather@feather.com",
		Name:  "feather",
	}

	fileCommitDetails := &types.FileCommitDetails{
		Author:  *author,
		Content: encodedYaml,
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

	if err := s.gitService.CreateFile(createReq); err != nil {
		return fmt.Errorf("failed to create application set file: %w", err)
	}

	return nil
}

func (s *argoCdServiceImpl) ensureArgoCdRepo(res *types.ProjectWithBaseCampInfo, repoName string) error {

	checkArgoCdRepoReq := &types.CheckRepoRequest{
		URL:   res.BaseCampURL,
		Token: res.Token,
		Owner: res.BaseCampOwner,
		Name:  repoName,
	}

	exists, err := s.gitService.RepoExists(checkArgoCdRepoReq)
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

		if err := s.gitService.CreateRepo(createReq); err != nil {
			return fmt.Errorf("failed to create ArgoCD repository: %w", err)
		}
		log.Printf("Repository '%s' created successfully.", repoName)
	} else {
		log.Printf("Repository '%s' already exists.", repoName)
	}

	return nil
}
