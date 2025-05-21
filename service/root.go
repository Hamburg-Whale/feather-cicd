package service

import (
	"feather/types"
	"fmt"
	"log"
)

type Service struct {
}

// NewService Repository를 주입받아 새로운 Service를 생성합니다.
func NewService() *Service {
	return &Service{}
}

// CreateUser 서비스에 대한 신규 사용자를 생성합니다.
func (service *Service) CreateUser(email string, password string) (*types.CreateUserReq, error) {
	_, err := fmt.Println("email: ", email, "password: ", password)
	if err != nil {
		log.Println("회원 생성에 실패했습니다. : ", "err", err.Error())
		return nil, err
	}
	return &types.CreateUserReq{
		Email:    email,
		Password: password,
	}, nil
}
