package tyk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	constants2 "tyk/tyk/bootstrap/pkg/constants"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"

	"k8s.io/apimachinery/pkg/util/json"
)

// BootstrapClassicPortal bootstraps Tyk Classic Portal.
func (s *Service) BootstrapClassicPortal() error {
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
	fmt.Println("Setting portal cname")

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

	return nil
}

func (s *Service) initialiseCatalogue() error {
	fmt.Println("Initialising Catalogue")

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
		fmt.Println(err)
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createPortalHomePage() error {
	fmt.Println("Creating portal homepage")

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
		fmt.Println(err)
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

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
	fmt.Println("Creating bootstrap default settings")

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

	return nil
}
