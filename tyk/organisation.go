package tyk

import (
	"bytes"
	"io"
	"net/http"
	"tyk/tyk/bootstrap/pkg/config"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/json"

	ic "tyk/tyk/bootstrap/pkg/constants"
)

// OrgExists checks if the given Tyk Organisation is created or not by
// checking owner_name and cname of the Organisation.
func (s *Service) OrgExists() (bool, error) {
	orgsApiEndpoint := s.appArgs.K8s.DashboardSvcUrl + constants.AdminOrganisationsEndpoint

	req, err := http.NewRequest(http.MethodGet, orgsApiEndpoint, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set(ic.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	orgs := api.OrgAPIResp{}

	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return false, err
	}

	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == s.appArgs.Tyk.Org.Name ||
				organisation["cname"] == s.appArgs.Tyk.Org.Cname {
				s.l.WithFields(logrus.Fields{
					"organisationName":  s.appArgs.Tyk.Org.Name,
					"organisationCname": s.appArgs.Tyk.Org.Cname,
				}).Debug("Organisation exists")

				return true, nil
			}
		}
	}

	s.l.WithFields(logrus.Fields{
		"organisationName":  s.appArgs.Tyk.Org.Name,
		"organisationCname": s.appArgs.Tyk.Org.Cname,
	}).Debug("Organisation does not exist")

	return false, nil
}

// CreateOrganisation creates organisation based on the information populated in the config.Config.
func (s *Service) CreateOrganisation() error {
	s.l.WithFields(logrus.Fields{
		"organisationName":  s.appArgs.Tyk.Org.Name,
		"organisationCname": s.appArgs.Tyk.Org.Cname,
	}).Debug("Creating Organisation")

	createOrgData := api.CreateOrgReq{
		OwnerName:    s.appArgs.Tyk.Org.Name,
		CnameEnabled: true,
		Cname:        s.appArgs.Tyk.Org.Cname,
	}

	// Enable hybrid in the organisation while setting up the MDCB Control Plane.
	// For reference: https://tyk.io/docs/tyk-multi-data-centre/setup-controller-data-centre/
	if s.appArgs.Tyk.Hybrid.Enabled {
		createOrgData.HybridEnabled = true
		createOrgData.EventOptions = eventOptions(s.appArgs.Tyk.Hybrid, s.appArgs.Tyk.Admin.EmailAddress)
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

	req.Header.Set(ic.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

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

	s.l.WithFields(logrus.Fields{
		"organisationName":  s.appArgs.Tyk.Org.Name,
		"organisationCname": s.appArgs.Tyk.Org.Cname,
	}).Debug("Created Organisation successfully")

	return nil
}

func eventOptions(hconf config.HybridConf, defaultEmail string) map[string]api.EventConfig {
	m := make(map[string]api.EventConfig)

	const (
		keyEventKey       = "key_event"
		hashedKeyEventKey = "hashed_key_event"
	)

	if hconf.HashedKeyEvent != nil {
		m[hashedKeyEventKey] = *hconf.HashedKeyEvent
	} else {
		m[hashedKeyEventKey] = api.EventConfig{Email: defaultEmail}
	}

	if hconf.KeyEvent != nil {
		m[keyEventKey] = *hconf.KeyEvent
	} else {
		m[keyEventKey] = api.EventConfig{Email: defaultEmail}
	}

	return m
}
