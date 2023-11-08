package tyk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	ic "tyk/tyk/bootstrap/pkg/constants"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"

	"k8s.io/apimachinery/pkg/util/json"
)

// CreateAdmin creates admin and sets its password based on the credentials populated in config.Config.
func (s *Service) CreateAdmin() error {
	adminData, err := s.createAdmin()
	if err != nil {
		return err
	}

	s.appArgs.Tyk.Admin.Auth = adminData.AuthCode

	err = s.setAdminPassword(adminData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) setAdminPassword(adminData NeededUserData) error {
	newPasswordData := api.ResetPasswordReq{
		NewPassword:     s.appArgs.Tyk.Admin.Password,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}

	reqBody, err := json.Marshal(newPasswordData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(constants.ApiUsersActionsResetEndpoint, s.appArgs.K8s.DashboardSvcUrl, adminData.UserId),
		bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set(ic.AuthorizationHeader, adminData.AuthCode)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("resetting password did not work")
	}

	return nil
}

type NeededUserData struct {
	AuthCode string
	UserId   string
}

func (s *Service) createAdmin() (NeededUserData, error) {
	s.l.Debug("Creating Admin User")

	reqBody := api.CreateUserReq{
		OrganisationId:  s.appArgs.Tyk.Org.ID,
		FirstName:       s.appArgs.Tyk.Admin.FirstName,
		LastName:        s.appArgs.Tyk.Admin.LastName,
		EmailAddress:    s.appArgs.Tyk.Admin.EmailAddress,
		Active:          true,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return NeededUserData{}, err
	}

	req, err := http.NewRequest(http.MethodPost, s.appArgs.K8s.DashboardSvcUrl+"/admin/users", bytes.NewReader(reqBytes))
	if err != nil {
		return NeededUserData{}, err
	}

	req.Header.Set(ic.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return NeededUserData{}, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return NeededUserData{}, err
	}

	if res.StatusCode == http.StatusForbidden {
		return NeededUserData{}, fmt.Errorf("user email already exists for this Org")
	}

	resp := api.CreateUserResp{}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return NeededUserData{}, err
	}

	s.l.Debug("Created Admin User successfully")

	return NeededUserData{UserId: resp.Meta.ID, AuthCode: resp.Message}, nil
}
