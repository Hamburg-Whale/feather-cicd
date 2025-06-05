package service

import (
	"bytes"
	"encoding/json"
	"feather/types"
	"fmt"
	"io"
	"log"
	"net/http"

	"k8s.io/client-go/rest"
)

type Service struct {
	httpClient *http.Client
}

func NewService() *Service {
	return &Service{
		httpClient: &http.Client{},
	}
}

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
