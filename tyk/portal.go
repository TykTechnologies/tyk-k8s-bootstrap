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

// BootstrapClassicPortal bootstraps Tyk Classic Portal.
func (s *Service) BootstrapClassicPortal() error {
	if s.appArgs.Tyk.Admin.Auth == "" {
		s.l.Warn("Missing Admin Auth configuration may cause Authorization failures on Tyk",
			"If Tyk Dashboard bootstrapping is disabled, please provide Admin Auth through environment variables")
	}

	err := s.createPortalDefaultSettings()
	if err != nil {
		return err
	}

	err = s.initialiseCatalogue()
	if err != nil {
		return err
	}

	err = s.createPortalHomePage()
	if err != nil {
		return err
	}

	err = s.setPortalCname()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) setPortalCname() error {
	s.l.Debug("Setting portal cname")

	if s.appArgs.Tyk.Org.Cname == "" {
		s.l.Warn("Missing Organisation Cname configuration may cause failures on Tyk requests",
			"Please provide Org Cname")
	}

	cnameReq := api.CnameReq{Cname: s.appArgs.Tyk.Org.Cname}

	reqBody, err := json.Marshal(cnameReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		s.appArgs.K8s.DashboardSvcUrl+constants.ApiPortalCnameEndpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set(constants2.AuthorizationHeader, s.appArgs.Tyk.Admin.Auth)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to set portal cname")
	}

	s.l.Debug("Set portal cname successfully")

	return nil
}

func (s *Service) initialiseCatalogue() error {
	s.l.Debug("Initialising Catalogue")

	if s.appArgs.Tyk.Org.ID == "" {
		s.l.Warn("Missing Organisation ID configuration may cause failures on Tyk requests",
			"Please provide Org ID")
	}

	initCatalog := api.InitCatalogReq{OrgId: s.appArgs.Tyk.Org.ID}

	reqBody, err := json.Marshal(initCatalog)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.appArgs.K8s.DashboardSvcUrl+constants.ApiPortalCatalogueEndpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set(constants2.AuthorizationHeader, s.appArgs.Tyk.Admin.Auth)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	resp := api.DashboardAPIResp{}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	s.l.Debug("Initialized Catalogue successfully")

	return nil
}

func (s *Service) createPortalHomePage() error {
	s.l.Debug("Creating portal homepage")

	homepageContents := portalHomepageReq()

	reqBody, err := json.Marshal(homepageContents)
	if err != nil {
		return err
	}

	reqData := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, s.appArgs.K8s.DashboardSvcUrl+constants.ApiPortalPagesEndpoint, reqData)
	if err != nil {
		return err
	}

	req.Header.Set(constants2.AuthorizationHeader, s.appArgs.Tyk.Admin.Auth)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	resp := api.DashboardAPIResp{}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	s.l.Debug("Created portal homepage successfully")

	return nil
}

func portalHomepageReq() api.PortalHomepageReq {
	return api.PortalHomepageReq{
		IsHomepage:   true,
		TemplateName: "",
		Title:        "Developer portal name",
		Slug:         "/",
		Fields: api.PortalFields{
			JumboCTATitle:       "Tyk Developer Portal",
			SubHeading:          "Sub Header",
			JumboCTALink:        "#cta",
			JumboCTALinkTitle:   "Your awesome APIs, hosted with Tyk!",
			PanelOneContent:     "Panel 1 content.",
			PanelOneLink:        "#panel1",
			PanelOneLinkTitle:   "Panel 1 Button",
			PanelOneTitle:       "Panel 1 Title",
			PanelThereeContent:  "",
			PanelThreeContent:   "Panel 3 content.",
			PanelThreeLink:      "#panel3",
			PanelThreeLinkTitle: "Panel 3 Button",
			PanelThreeTitle:     "Panel 3 Title",
			PanelTwoContent:     "Panel 2 content.",
			PanelTwoLink:        "#panel2",
			PanelTwoLinkTitle:   "Panel 2 Button",
			PanelTwoTitle:       "Panel 2 Title",
		},
	}
}

func (s *Service) createPortalDefaultSettings() error {
	s.l.Debug("Creating bootstrap default settings")

	// TODO(buraksekili): DashboardSvcUrl can be populated via environment variables. So, the URL
	// might have trailing slashes. Constructing the URL with raw string concatenating is not a good
	// approach here. Needs refactoring.
	req, err := http.NewRequest(
		http.MethodPut,
		s.appArgs.K8s.DashboardSvcUrl+constants.ApiPortalConfigurationEndpoint,
		nil,
	)
	req.Header.Set(constants2.AuthorizationHeader, s.appArgs.Tyk.Admin.Auth)

	if err != nil {
		return err
	}

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	s.l.Debug("Created bootstrap default settings successfully")

	return nil
}
