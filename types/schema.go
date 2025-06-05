package types

type User struct {
	ID int64 `json:"id" binding:"required"`

	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type BaseCamp struct {
	ID int64 `json:"id" binding:"required"`

	Name  string `json:"name" binding:"required"`
	URL   string `json:"url" binding:"required"`
	Token string `json:"token" binding:"required"`

	User_ID int64 `json:"user_id" binding:"required"`
}

type Project struct {
	ID int64 `json:"id" binding:"required"`

	Name    string `json:"name" binding:"required"`
	URL     string `json:"url" binding:"required"`
	Owner   string `json:"owner" binding:"required"`
	Private bool   `json:"private,omitempty"`

	BaseCamp_ID int64 `json:"basecamp_id" binding:"required"`
}
