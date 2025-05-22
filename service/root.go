package service

import (
	"bytes"
	"encoding/json"
	"feather/types"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Service struct {
}

// NewService Repository를 주입받아 새로운 Service를 생성합니다.
func NewService() *Service {
	return &Service{}
}

// CreateUser 서비스에 대한 신규 사용자를 생성합니다.
type CreateUserRes struct {
	Email string `json:"email"`
}

func (service *Service) CreateUser(email string, password string) (*types.Response, error) {
	_, err := fmt.Println("email: ", email, "password: ", password)
	if err != nil {
		log.Println("회원 생성에 실패했습니다. : ", err.Error())
		return nil, err
	}

	return types.NewRes(200, CreateUserRes{
		Email: email,
	}, "User created successfully"), nil
}

// func (service *Service) AuthUser(req *types.AuthUserReq) (*types.Response, error) {
// 	url := req.Url
// 	prefix := "/api/v1/user"

// 	authReq, err := http.NewRequest("GET", url + prefix, nil)
// 	if err != nil {
// 		log.Println("요청 생성 실패 : ", err.Error())
// 		return types.NewRes(401, nil, "요청 생성 실패")
// 	}

// 	authReq.Header.Set("Authorization", "token " + req.Token)

// 	client := &http.Client{}
// 	res, err := client.Do(authReq)
// 	if err != nil {
// 		log.Println("요청 실패 : ", err.Error())
// 		return types.NewRes(401, nil, "요청 실패")
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		fmt.Printf("인증 실패: HTTP %d\n", res.StatusCode)
// 		return types.NewRes(401, nil, "인증 실패")
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println("응답 읽기 실패:", err)
// 		os.Exit(1)
// 	}

// 	var expectedUser types.GiteaUser
// 	if err := json.Unmarshal(body, &expectedUser); err != nil {
// 		fmt.Println("JSON 파싱 실패:", err)
// 		os.Exit(1)
// 	}

// 	if expectedUser.Login == req.Username {
// 		fmt.Println("토큰과 사용자 이름이 일치")
// 	} else {
// 		fmt.Printf("사용자 이름 불일치")
// 	}

// 	return types.NewRes(200, nil, "사용자 인증 성공")
// }

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

func (service *Service) CreateRepo(req *types.CreateRepoReq) (*types.Response, error) {
	repoURL := fmt.Sprintf("%s/repos/%s/%s/generate", req.Url, req.TemplateOwner, req.TemplateRepo)

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
	hookURL := fmt.Sprintf("%s/repos/%s/%s/hooks", baseURL, owner, repo)

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
