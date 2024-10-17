package models

type AddStoryRequest struct {
    Story map[string]interface{} `json:"story"`
}

type UserStoryAddSuccessModel struct{
    Status  bool        `json:"status"`
	Message string      `json:"message"`
}