package models

type LoginUserModel struct{
	Email string `json:"email"`;
	Password string `json:"password"`;
}

type LoginResponse struct {
    ID         int    `json:"id"`
    FullName   string `json:"full_name"`
    Email      string `json:"email"`
    ProfilePic string `json:"profile_pic"`
}