package api

import "time"

type DashboardAPIResp struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    string `json:"Meta"`
}

type OrgAPIResp struct {
	Organisations []map[string]interface{} `json:"organisations"`
	Pages         int                      `json:"pages"`
}

type GetUserResp struct {
	APIModel        struct{}          `json:"api_model"`
	FirstName       string            `json:"first_name"`
	LastName        string            `json:"last_name"`
	EmailAddress    string            `json:"email_address"`
	Password        string            `json:"password"`
	OrgID           string            `json:"org_id"`
	Active          bool              `json:"active"`
	ID              string            `json:"id"`
	AccessKey       string            `json:"access_key"`
	UserPermissions map[string]string `json:"user_permissions"`
}

type ListUsersResp []GetUserResp

type CreateUserResp struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    struct {
		APIModel        struct{} `json:"api_model"`
		FirstName       string   `json:"first_name"`
		LastName        string   `json:"last_name"`
		EmailAddress    string   `json:"email_address"`
		OrgID           string   `json:"org_id"`
		Active          bool     `json:"active"`
		ID              string   `json:"id"`
		AccessKey       string   `json:"access_key"`
		UserPermissions struct {
			IsAdmin string `json:"IsAdmin"`
		} `json:"user_permissions"`
		GroupID         string        `json:"group_id"`
		PasswordMaxDays int           `json:"password_max_days"`
		PasswordUpdated time.Time     `json:"password_updated"`
		PWHistory       []interface{} `json:"PWHistory"`
		CreatedAt       time.Time     `json:"created_at"`
	} `json:"Meta"`
}
