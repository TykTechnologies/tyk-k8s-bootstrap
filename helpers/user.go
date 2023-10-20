package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TykTechnologies/tyk-k8s-bootstrap/data"
	"io"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

func CreateUser(client http.Client, dashboardUrl string, orgId string) (string, error) {
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

func SetUserPassword(client http.Client, userId string, authCode string, dashboardUrl string) error {
	newPasswordData := ResetPasswordStruct{
		NewPassword:     data.AppConfig.TykAdminPassword,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}
	reqBody, err := json.Marshal(newPasswordData)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(ApiUsersActionsResetEndpoint, dashboardUrl, userId),
		bytes.NewReader(reqBody))
	req.Header.Set("authorization", authCode)
	req.Header.Set("Content-Type", "application/json")
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
	orgId, err := CreateOrganisation(client, data.AppConfig.DashboardUrl)
	if err != nil {
		return err
	}

	data.AppConfig.OrgId = orgId

	userAuth, err := CreateUser(client, data.AppConfig.DashboardUrl, orgId)
	if err != nil {
		return err
	}

	data.AppConfig.UserAuth = userAuth

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
		APIModel struct {
		} `json:"api_model"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		EmailAddress    string `json:"email_address"`
		OrgID           string `json:"org_id"`
		Active          bool   `json:"active"`
		ID              string `json:"id"`
		AccessKey       string `json:"access_key"`
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

func GetUserData(client http.Client, dashboardUrl string, orgId string) (NeededUserData, error) {
	reqBody := CreateUserRequest{
		OrganisationId:  orgId,
		FirstName:       data.AppConfig.TykAdminFirstName,
		LastName:        data.AppConfig.TykAdminLastName,
		EmailAddress:    data.AppConfig.TykAdminEmailAddress,
		Active:          true,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return NeededUserData{}, err
	}
	reqData := bytes.NewReader(reqBytes)
	req, err := http.NewRequest("POST", dashboardUrl+"/admin/users", reqData)
	if err != nil {
		return NeededUserData{}, err
	}

	req.Header.Set("admin-auth", data.AppConfig.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")
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
