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

func (s *Service) BoostrapPortal() error {
	err := s.createPortalDefaultSettings()
	if err != nil {
		return err
	}

	err = s.initialiseCatalogue()
	if err != nil {
		return err
	}

	err = s.createPortalHomepage()
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

	cnameReq := api.CnameRequest{Cname: s.appArgs.TykPortalCname}
	reqBody, err := json.Marshal(cnameReq)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)

	req, err := http.NewRequest("PUT", s.appArgs.DashboardSvcAddr+ApiPortalCnameEndpoint, reqData)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.appArgs.TykUserAuth)

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

	initCatalog := api.InitCatalogReq{OrgId: s.appArgs.TykOrgId}
	reqBody, err := json.Marshal(initCatalog)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", s.appArgs.DashboardSvcAddr+ApiPortalCatalogueEndpoint, reqData)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.appArgs.TykUserAuth)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}
	resp := api.DashboardGeneralResponse{}
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

func (s *Service) createPortalHomepage() error {
	fmt.Println("Creating portal homepage")

	homepageContents := portalHomePageRequest()
	reqBody, err := json.Marshal(homepageContents)
	if err != nil {
		return err
	}

	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", s.appArgs.DashboardSvcAddr+ApiPortalPagesEndpoint, reqData)
	req.Header.Set("Authorization", s.appArgs.TykUserAuth)

	if err != nil {
		return err
	}

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func portalHomePageRequest() api.PortalHomepageRequest {
	return api.PortalHomepageRequest{
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

	req, err := http.NewRequest("POST", s.appArgs.DashboardSvcAddr+ApiPortalConfigurationEndpoint, nil)
	req.Header.Set("Authorization", s.appArgs.TykUserAuth)

	if err != nil {
		return err
	}

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
