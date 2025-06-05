package service

import (
	"bytes"
	"encoding/json"
	"feather/repository"
	"feather/types"
	"fmt"
	"io"
	"log"
	"net/http"

	"k8s.io/client-go/rest"
)

type Service struct {
	httpClient *http.Client
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		httpClient: &http.Client{},
		repository: repository,
	}
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

func GetKubeConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func (s *Service) DoJSONGet(url, token string) (*http.Response, error) {
	return s.doJSONRequest("GET", url, token, nil)
}

func (s *Service) DoJSONPost(url, token string, payload interface{}) (*http.Response, error) {
	return s.doJSONRequest("POST", url, token, payload)
}

func (s *Service) doJSONRequest(method, url, token string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if res.StatusCode >= 300 {
		defer res.Body.Close()
		respBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, string(respBody))
	}

	return res, nil
}

func (s *Service) CreateUser(email string, password string, nickname string) error {
	err := s.repository.CreateUser(email, password, nickname)
	if err != nil {
		log.Println("회원 생성에 실패했습니다. : ", "err", err.Error())
		return fmt.Errorf("회원 생성 실패: %w", err)
	}
	return nil
}

func (s *Service) User(userId int64) (*types.User, error) {
	res, err := s.repository.User(userId)
	if err != nil {
		log.Println("회원 조회에 실패했습니다. : ", "err", err.Error())
		return nil, fmt.Errorf("회원 조회 실패: %w", err)
	}
	return res, nil
}

func (s *Service) CreateBaseCamp(name string, url string, token string, userId int64) error {
	err := s.repository.CreateBaseCamp(name, url, token, userId)
	if err != nil {
		log.Println("베이스캠프 생성에 실패했습니다. : ", "err", err.Error())
		return fmt.Errorf("베이스캠프 생성 실패: %w", err)
	}
	return nil
}

func (s *Service) BaseCampsByUserId(userId int64) ([]*types.BaseCamp, error) {
	res, err := s.repository.BaseCampsByUserId(userId)
	if err != nil {
		log.Println("베이스캠프 조회에 실패했습니다. : ", "err", err.Error())
		return nil, fmt.Errorf("베이스캠프 조회 실패: %w", err)
	}
	return res, nil
}

func (s *Service) BaseCamp(baseCampId int64) (*types.BaseCamp, error) {
	res, err := s.repository.BaseCamp(baseCampId)
	if err != nil {
		log.Println("베이스캠프 조회에 실패했습니다. : ", "err", err.Error())
		return nil, fmt.Errorf("베이스캠프 조회 실패: %w", err)
	}
	return res, nil
}

func (s *Service) CreateProject(name string, url string, owner string, private bool, baseCampId int64) error {
	err := s.repository.CreateProject(name, url, owner, private, baseCampId)
	if err != nil {
		log.Println("프로젝트 생성에 실패했습니다. : ", "err", err.Error())
		return fmt.Errorf("프로젝트 생성 실패: %w", err)
	}
	return nil
}

func (s *Service) ProjectsByBaseCampId(baseCampId int64) ([]*types.Project, error) {
	res, err := s.repository.ProjectsByBaseCampId(baseCampId)
	if err != nil {
		log.Println("프로젝트 조회에 실패했습니다. : ", "err", err.Error())
		return nil, fmt.Errorf("프로젝트 조회 실패: %w", err)
	}
	return res, nil
}

func (s *Service) Project(projectId int64) (*types.Project, error) {
	res, err := s.repository.Project(projectId)
	if err != nil {
		log.Println("프로젝트 조회에 실패했습니다. : ", "err", err.Error())
		return nil, fmt.Errorf("프로젝트 조회 실패: %w", err)
	}
	return res, nil
}
