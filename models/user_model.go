package models

type RegisterUserModel struct{
	ID int `json:"id"`;
	FullName string `json:"full_name"`;
	Email string `json:"email"`;
	ProfilePic string `json:"profile_pic"`;
	Password string `json:"password"`;
}

type GetUserProfileModel struct{
	FullName string `json:"full_name"`;
	Email string `json:"email"`;
	ProfilePic string `json:"profile_pic"`;
}