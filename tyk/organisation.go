package tyk

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	constants2 "tyk/tyk/bootstrap/pkg/constants"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"

	"k8s.io/apimachinery/pkg/util/json"
)

var ErrOrgExists = errors.New("there shouldn't be any organisations, please " +
	"disable bootstrapping to avoid losing data or delete " +
	"already existing organisations")

func (s *Service) OrgExists() error {
	//s.l.Info(
	//	"looking if organisation exists",
	//	"Org Cname", s.appArgs.Tyk.Org.Cname,
	//	"Org Name", s.appArgs.Tyk.Org.Name,
	//)

	orgsApiEndpoint := s.appArgs.K8s.DashboardSvcUrl + constants.AdminOrganisationsEndpoint

	req, err := http.NewRequest(http.MethodGet, orgsApiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set(constants2.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(constants2.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	orgs := api.OrgAPIResp{}

	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return err
	}

	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == s.appArgs.Tyk.Org.Name ||
				organisation["cname"] == s.appArgs.Tyk.Org.Cname {
				//s.l.Info(
				//	"looking if organisation exists",
				//	"Org Cname", s.appArgs.Tyk.Org.Cname,
				//	"Org Name", s.appArgs.Tyk.Org.Name,
				//)

				return ErrOrgExists
			}
		}
	}

	return nil
}

func (s *Service) CreateOrganisation() error {
	//s.l.Info(
	//	"creating an organisation",
	//	"Name", s.appArgs.Tyk.Org.Name,
	//	"Cname", s.appArgs.Tyk.Org.Cname,
	//)

	createOrgData := api.CreateOrgReq{
		OwnerName:    s.appArgs.Tyk.Org.Name,
		CnameEnabled: true,
		Cname:        s.appArgs.Tyk.Org.Cname,
	}

	reqBodyBytes, err := json.Marshal(createOrgData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.appArgs.K8s.DashboardSvcUrl+constants.AdminOrganisationsEndpoint,
		bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set(constants2.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(constants2.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	createOrgResp := api.DashboardAPIResp{}

	err = json.Unmarshal(bodyBytes, &createOrgResp)
	if err != nil {
		return err
	}

	s.appArgs.Tyk.Org.ID = createOrgResp.Meta

	return nil
}
