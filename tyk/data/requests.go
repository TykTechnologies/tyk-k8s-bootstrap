package data

type CreateOrgRequest struct {
	OwnerName    string `json:"owner_name"`
	CnameEnabled bool   `json:"cname_enabled"`
	Cname        string `json:"cname"`
}

type CreateUserRequest struct {
	OrganisationId  string            `json:"org_id"`
	FirstName       string            `json:"first_name"`
	LastName        string            `json:"last_name"`
	EmailAddress    string            `json:"email_address"`
	Active          bool              `json:"active"`
	UserPermissions map[string]string `json:"user_permissions"`
}

type ResetPasswordRequest struct {
	NewPassword     string            `json:"new_password"`
	UserPermissions map[string]string `json:"user_permissions"`
}

type InitCatalogReq struct {
	OrgId string `json:"org_id"`
}

type CnameRequest struct {
	Cname string `json:"cname"`
}

type PortalHomepageRequest struct {
	IsHomepage   bool         `json:"is_homepage"`
	TemplateName string       `json:"template_name"`
	Title        string       `json:"title"`
	Slug         string       `json:"slug"`
	Fields       PortalFields `json:"fields"`
}
