package soap

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	Graph                  = "https://graph.windows.net"
	provisioningWebService = "https://provisioningapi.microsoftonline.com/provisioningwebservice.svc"
	applicationSoap        = "application/soap+xml"
)

type MsolConnect struct {
	ClientId       string
	BearerToken    string
	TrackingHeader string
	MessageID      string
}
type MsolConnectResponse struct {
	Header *HeaderResponse `xml:"Header"`
}

type HeaderResponse struct {
	BecContext *BecContextResponse `xml:"BecContext"`
}

type BecContextResponse struct {
	DataBlob string `xml:"DataBlob"`
}

type GetCompanyInformationResponse struct {
}
type GetCompanyInformation struct {
	DataBlob       string
	BearerToken    string
	TrackingHeader string
	MessageID      string
}

//MonitoringAzureAD monitore Azure AD
type MonitoringAzureAD struct {
	verbose      bool
	client       *http.Client
	login        string
	passwd       string
	clientID     string
	tenantDomain string
	tenantID     string
}

//Token to talk with API
type Token struct {
	//Type Bearer
	AccessToken string `json:"access_token"`
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

//NewMonitoringAzureAD create monitore
func NewMonitoringAzureAD(login string, passwd string, clientID string, tenantDomain, tenantID string, verbose bool) *MonitoringAzureAD {
	client := &http.Client{}

	return &MonitoringAzureAD{
		verbose:      verbose,
		client:       client,
		login:        login,
		passwd:       passwd,
		clientID:     clientID,
		tenantDomain: tenantDomain,
		tenantID:     tenantID,
	}

}

func (m *MonitoringAzureAD) AccessToken(resource string) (string, error) {
	const (
		//TODO try with common
		loginURLPrefix = "https://login.microsoftonline.com/"
		loginURLApi    = "/oauth2/token"
	)
	resp, err := m.client.PostForm(fmt.Sprint(loginURLPrefix, m.tenantID, loginURLApi), url.Values{
		"username":   {m.login},
		"password":   {m.passwd},
		"resource":   {resource},
		"scope":      {"openid"},
		"grant_type": {"password"},
		"client_id":  {m.clientID},
	})
	defer closeBody(resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed %d", resp.StatusCode)
	}

	token := &Token{}
	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return "", err
	}

	return token.AccessToken, nil

}

func payload(name string, templateContent string, datum interface{}) (io.Reader, error) {
	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return nil, err
	}

	read, write := io.Pipe()
	go func(write *io.PipeWriter, tmpl *template.Template, datum interface{}) {
		if err := tmpl.Execute(write, datum); err != nil {
			//Scheiße we have a error
			write.CloseWithError(err)
		} else {
			write.Close()
		}
	}(write, tmpl, datum)
	return read, nil

}

func (m *MonitoringAzureAD) postSoap(reader io.Reader, datum interface{}) error {
	resp, err := m.client.Post(provisioningWebService, applicationSoap, reader)
	defer closeBody(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		payload, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("postSoap failed %d %s", resp.StatusCode, payload)
	}
	//TODO aahhhhhhhh
	if m.verbose {
		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)
		err = xml.NewDecoder(tee).Decode(datum)
		log.Println(string(buf.Bytes()))
		return err
	}
	return xml.NewDecoder(resp.Body).Decode(datum)
}

func (m *MonitoringAzureAD) MsolConnect(token string) (string, error) {
	const soapMsolConnect = `<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing"><s:Header><a:Action s:mustUnderstand="1">http://provisioning.microsoftonline.com/IProvisioningWebService/MsolConnect</a:Action><a:MessageID>urn:uuid:{{.MessageID}}</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><UserIdentityHeader xmlns="http://provisioning.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BearerToken xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">Bearer {{.BearerToken}}</BearerToken><LiveToken i:nil="true" xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService"/></UserIdentityHeader><ClientVersionHeader xmlns="http://provisioning.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><ClientId xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">{{.ClientId}}</ClientId><Version xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">1.2.166.0</Version></ClientVersionHeader><ContractVersionHeader xmlns="http://becwebservice.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BecVersion xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">Version38</BecVersion></ContractVersionHeader><TrackingHeader xmlns="http://becwebservice.microsoftonline.com/">{{.TrackingHeader}}</TrackingHeader><a:To s:mustUnderstand="1">https://provisioningapi.microsoftonline.com/provisioningwebservice.svc</a:To></s:Header><s:Body><MsolConnect xmlns="http://provisioning.microsoftonline.com/"><request xmlns:b="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><b:BecVersion>Version4</b:BecVersion><b:TenantId i:nil="true"/><b:VerifiedDomain i:nil="true"/></request></MsolConnect></s:Body></s:Envelope>`

	msol := &MsolConnect{
		ClientId:       m.tenantID,
		BearerToken:    token,
		MessageID:      "",
		TrackingHeader: "",
	}

	reader, err := payload("soapMsolConnect", soapMsolConnect, msol)
	if err != nil {
		return "", err
	}

	connectResponse := &MsolConnectResponse{}
	if err = m.postSoap(reader, connectResponse); err != nil {
		return "", err
	}

	if connectResponse.Header != nil && connectResponse.Header.BecContext != nil && connectResponse.Header.BecContext.DataBlob != "" {
		return connectResponse.Header.BecContext.DataBlob, nil
	}

	return "", fmt.Errorf("GetCompanyInformation failed")
}

func (m *MonitoringAzureAD) GetCompanyInformation(token, dataBlob string) (interface{}, error) {
	const soapGetCompanyInformation = `<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing"><s:Header><a:Action s:mustUnderstand="1">http://provisioning.microsoftonline.com/IProvisioningWebService/GetCompanyInformation</a:Action><a:MessageID>urn:uuid:{{.MessageID}}</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><UserIdentityHeader xmlns="http://provisioning.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BearerToken xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">Bearer {{.BearerToken}}</BearerToken><LiveToken i:nil="true" xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService"/></UserIdentityHeader><BecContext xmlns="http://becwebservice.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><DataBlob xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">{{.DataBlob}}</DataBlob></BecContext><ClientVersionHeader xmlns="http://provisioning.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><ClientId xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">50afce61-c917-435b-8c6d-60aa5a8b8aa7</ClientId><Version xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">1.2.166.0</Version></ClientVersionHeader><ContractVersionHeader xmlns="http://becwebservice.microsoftonline.com/" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BecVersion xmlns="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService">Version38</BecVersion></ContractVersionHeader><TrackingHeader xmlns="http://becwebservice.microsoftonline.com/">{{.TrackingHeader}}</TrackingHeader><a:To s:mustUnderstand="1">https://provisioningapi.microsoftonline.com/provisioningwebservice.svc</a:To></s:Header><s:Body><GetCompanyInformation xmlns="http://provisioning.microsoftonline.com/"><request xmlns:b="http://schemas.datacontract.org/2004/07/Microsoft.Online.Administration.WebService" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><b:BecVersion>Version16</b:BecVersion><b:TenantId i:nil="true"/><b:VerifiedDomain i:nil="true"/></request></GetCompanyInformation></s:Body></s:Envelope>`

	getCompanyInformation := &GetCompanyInformation{
		DataBlob:       dataBlob,
		BearerToken:    token,
		MessageID:      "",
		TrackingHeader: "",
	}

	reader, err := payload("soapGetCompanyInformation", soapGetCompanyInformation, getCompanyInformation)
	if err != nil {
		return "", err
	}

	response := &GetCompanyInformationResponse{}
	if err = m.postSoap(reader, response); err != nil {
		return nil, err
	}

	return response, err
}
