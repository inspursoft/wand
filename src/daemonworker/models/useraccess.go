package models

type UserAccess struct {
	ID          int64
	Username    string
	AccessToken string
	IsOrg       int
}
