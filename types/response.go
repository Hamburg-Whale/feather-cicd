package types

import (
	"strings"
	"time"
)

type header struct {
	Result int    `json:"result"`
	Data   string `json:"data"`
}

func newHeader(result int, data ...string) *header {
	return &header{
		Result: result,
		Data:   strings.Join(data, ", "),
	}
}

type Response struct {
	*header
	Result interface{} `json:"result"`
}

func NewRes(result int, res interface{}, data ...string) *Response {
	return &Response{
		header: newHeader(result, data...),
		Result: res,
	}
}

type GiteaUser struct {
	Active            bool      `json:"active"`
	AvatarURL         string    `json:"avatar_url"`
	Created           time.Time `json:"created"`
	Description       string    `json:"description"`
	Email             string    `json:"email"`
	FollowersCount    int64     `json:"followers_count"`
	FollowingCount    int64     `json:"following_count"`
	FullName          string    `json:"full_name"`
	HTMLURL           string    `json:"html_url"`
	ID                int64     `json:"id"`
	IsAdmin           bool      `json:"is_admin"`
	Language          string    `json:"language"`
	LastLogin         time.Time `json:"last_login"`
	Location          string    `json:"location"`
	Login             string    `json:"login"`
	LoginName         string    `json:"login_name"`
	ProhibitLogin     bool      `json:"prohibit_login"`
	Restricted        bool      `json:"restricted"`
	SourceID          int64     `json:"source_id"`
	StarredReposCount int64     `json:"starred_repos_count"`
	Visibility        string    `json:"visibility"`
	Website           string    `json:"website"`
}
