package tyk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"tyk/tyk/bootstrap/tyk/api"
)

func (s *Service) createUser(dashboardUrl, orgId string) (string, error) {
	userData, err := s.getUserData(dashboardUrl, orgId)

	if err != nil {
		return "", err
	}

	fmt.Println(userData)

	err = s.setUserPassword(userData.UserId, userData.AuthCode, dashboardUrl)
	if err != nil {
		return "", err
	}

	return userData.AuthCode, nil
}

func (s *Service) setUserPassword(userId, authCode, dashboardUrl string) error {
	newPasswordData := api.ResetPasswordRequest{
		NewPassword:     s.appArgs.TykAdminPassword,
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

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		fmt.Println(string(b))
		return errors.New("resetting password did not work")
	}

	return nil
}

func (s *Service) GenerateDashboardCredentials() error {
	orgId, err := s.createOrganisation(s.appArgs.DashboardUrl)
	if err != nil {
		return err
	}

	s.appArgs.OrgId = orgId

	userAuth, err := s.createUser(s.appArgs.DashboardUrl, orgId)
	if err != nil {
		return err
	}

	s.appArgs.UserAuth = userAuth

	return nil
}

type NeededUserData struct {
	AuthCode string
	UserId   string
}

func (s *Service) getUserData(dashboardUrl, orgId string) (NeededUserData, error) {
	reqBody := api.CreateUserRequest{
		OrganisationId:  orgId,
		FirstName:       s.appArgs.TykAdminFirstName,
		LastName:        s.appArgs.TykAdminLastName,
		EmailAddress:    s.appArgs.TykAdminEmailAddress,
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

	req.Header.Set("admin-auth", s.appArgs.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return NeededUserData{}, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	getUserResponse := api.CreateUserResponse{}
	err = json.Unmarshal(bodyBytes, &getUserResponse)
	if err != nil {
		return NeededUserData{}, err
	}

	fmt.Println("getuserresponse", getUserResponse)

	return NeededUserData{UserId: getUserResponse.Meta.ID, AuthCode: getUserResponse.Message}, nil
}
