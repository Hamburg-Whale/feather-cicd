package service

import (
	"feather/types"
	"fmt"
)

func (service *Service) CreateProjectManifestRepo(req *types.CreateRepoReq) error {
	checkRepoReq := &types.CheckRepoReq{
		Url:   req.Url,
		Owner: req.Owner,
		Name:  "feather-argocd",
		Token: req.Token,
	}

	if !repoIsExist(checkRepoReq) {
		argoRepoReq := &types.CreateRepoReq{
			Url:         req.Url,
			Description: "ArgoCD 매니페스트 관리용 리포지토리",
			Name:        "feather-argocd",
			Owner:       req.Owner,
			Private:     false,
		}
		createRepo(argoRepoReq)
	}

	return nil
}

func repoIsExist(req *types.CheckRepoReq) bool {
	repoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s", req.Url, req.Owner, req.Name)

	_, err := DoJSONGet(repoURL, req.Token)
	if err != nil {
		return true
	}
	return false
}

func createRepo(req *types.CreateRepoReq) error {
	repoURL := fmt.Sprintf("%s/api/v1/user/repos", req.Url)

	payload := map[string]interface{}{
		"Description": req.Description,
		"Name":        req.Name,
		"Private":     req.Private,
	}

	_, err := DoJSONPost(repoURL, req.Token, payload)
	if err != nil {
		return fmt.Errorf("ArgoCD 리포지토리 생성 실패: %w", err)
	}
	return nil
}
