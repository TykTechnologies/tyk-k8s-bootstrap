package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"tyk/tyk/bootstrap/data"

	"k8s.io/apimachinery/pkg/util/json"
)

func CreateUser(client http.Client, dashboardUrl, orgId string) (string, error) {
	userData, err := GetUserData(client, dashboardUrl, orgId)
	if err != nil {
		return "", err
	}

	err = SetUserPassword(client, userData.UserId, userData.AuthCode, dashboardUrl)
	if err != nil {
		return "", err
	}

	return userData.AuthCode, nil
}

type ResetPasswordStruct struct {
	NewPassword     string            `json:"new_password"`
	UserPermissions map[string]string `json:"user_permissions"`
}

func SetUserPassword(client http.Client, userId, authCode, dashboardUrl string) error {
	newPasswordData := ResetPasswordStruct{
		NewPassword:     data.BootstrapConf.Tyk.Admin.Password,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}

	reqBody, err := json.Marshal(newPasswordData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(ApiUsersActionsResetEndpoint, dashboardUrl, userId),
		bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set(data.AuthorizationHeader, authCode)
	req.Header.Set(data.ContentTypeHeader, "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("resetting password did not work")
	}

	return nil
}

func GenerateDashboardCredentials(client http.Client) error {
	orgId, err := CreateOrganisation(client, data.BootstrapConf.K8s.DashboardSvcUrl)
	if err != nil {
		return err
	}

	data.BootstrapConf.Tyk.Org.ID = orgId

	userAuth, err := CreateUser(client, data.BootstrapConf.K8s.DashboardSvcUrl, orgId)
	if err != nil {
		return err
	}

	data.BootstrapConf.Tyk.Admin.Auth = userAuth

	return nil
}

type CreateUserRequest struct {
	OrganisationId  string            `json:"org_id"`
	FirstName       string            `json:"first_name"`
	LastName        string            `json:"last_name"`
	EmailAddress    string            `json:"email_address"`
	Active          bool              `json:"active"`
	UserPermissions map[string]string `json:"user_permissions"`
}

type CreateUserResponse struct {
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

type NeededUserData struct {
	AuthCode string
	UserId   string
}

func GetUserData(client http.Client, dashboardUrl, orgId string) (NeededUserData, error) {
	reqBody := CreateUserRequest{
		OrganisationId:  orgId,
		FirstName:       data.BootstrapConf.Tyk.Admin.FirstName,
		LastName:        data.BootstrapConf.Tyk.Admin.LastName,
		EmailAddress:    data.BootstrapConf.Tyk.Admin.EmailAddress,
		Active:          true,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return NeededUserData{}, err
	}

	req, err := http.NewRequest(http.MethodPost, dashboardUrl+"/admin/users", bytes.NewReader(reqBytes))
	if err != nil {
		return NeededUserData{}, err
	}

	req.Header.Set(data.AdminAuthHeader, data.BootstrapConf.Tyk.Admin.Secret)
	req.Header.Set(data.ContentTypeHeader, "application/json")

	res, err := client.Do(req)
	if err != nil {
		return NeededUserData{}, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	getUserResponse := CreateUserResponse{}

	err = json.Unmarshal(bodyBytes, &getUserResponse)
	if err != nil {
		return NeededUserData{}, err
	}

	return NeededUserData{UserId: getUserResponse.Meta.ID, AuthCode: getUserResponse.Message}, nil
}
