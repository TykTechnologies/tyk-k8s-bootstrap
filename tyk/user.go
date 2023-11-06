package tyk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	constants2 "tyk/tyk/bootstrap/pkg/constants"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"
)

func (s *Service) CreateAdmin() error {
	adminData, err := s.createAdmin()
	if err != nil {
		return err
	}

	err = s.setAdminPassword(adminData)
	if err != nil {
		return err
	}

	s.appArgs.Tyk.Admin.Auth = adminData.AuthCode

	return nil
}

//func (s *Service) UserExists(userEmail string) (*NeededUserData, error) {
//	s.l.Info("checking if users exists")
//
//	req, err := http.NewRequest(http.MethodGet, s.appArgs.K8s.DashboardSvcUrl+constants.ApiUsers, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header.Set(data.AuthorizationHeader, s.appArgs.Tyk.Admin.Secret)
//	req.Header.Set(data.ContentTypeHeader, "application/json")
//
//	res, err := s.httpClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	bodyBytes, err := io.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	users := api.ListUsersResp{}
//
//	fmt.Printf("%#v\n", string(bodyBytes))
//
//	err = json.Unmarshal(bodyBytes, &users)
//	if err != nil {
//		return nil, err
//	}
//
//	for i := range users {
//		user := users[i]
//		if user.EmailAddress == userEmail {
//			s.appArgs.Tyk.Admin.Auth = user.AccessKey
//			s.appArgs.Tyk.Org.ID = user.OrgID
//			return &NeededUserData{UserId: user.ID, AuthCode: user.AccessKey}, nil
//		}
//	}
//
//	return nil, nil
//}

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

	req.Header.Set(constants2.AuthorizationHeader, adminData.AuthCode)
	req.Header.Set(constants2.ContentTypeHeader, "application/json")

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

var ErrUserExists = fmt.Errorf("User email already exists for this Org")

func (s *Service) createAdmin() (NeededUserData, error) {
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

	req.Header.Set(constants2.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(constants2.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return NeededUserData{}, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return NeededUserData{}, err
	}

	if res.StatusCode == http.StatusForbidden {
		return NeededUserData{}, ErrUserExists
	}

	resp := api.CreateUserResp{}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return NeededUserData{}, err
	}

	return NeededUserData{UserId: resp.Meta.ID, AuthCode: resp.Message}, nil
}
