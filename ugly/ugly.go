package ugly

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"golang.org/x/net/html"
)

var (
	reFlowToken = regexp.MustCompile("\"sFT\":\"([^,]+)\"")
	reSCtx      = regexp.MustCompile("\"sCtx\":\"([^,]+)\"")
	reCanary    = regexp.MustCompile("\"canary\":\"([^,]+)\"")
)

//MonitoringAzureAD client http login and password
type MonitoringAzureAD struct {
	client *http.Client
	login  string
	passwd string
}

//NewMonitoringAzureAD prepare the client http
func NewMonitoringAzureAD(login string, passwd string) *MonitoringAzureAD {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	return &MonitoringAzureAD{
		client: client,
		login:  login,
		passwd: passwd,
	}

}

//DirSyncManagement the data
type DirSyncManagement struct {
	IsPasswordSyncEnabled    bool
	IsPasswordSyncNormal     bool
	IsPasswordSyncRedWarning bool

	IsDirSyncEnabled      bool
	IsDirSyncObjectErrors bool
	IsDirSyncRedWarning   bool

	PasswordSyncLastSyncTime string
	DirSyncLastSyncTime      string
}

//MetaData read in html page
type MetaData struct {
	flowToken string
	sctx      string
	canary    string
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
func (mo *MonitoringAzureAD) getMetaData() (*MetaData, error) {

	resp, err := mo.client.Get("https://login.microsoftonline.com")
	defer closeBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Houston getMetaData is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	b := string(body)
	flowToken := reFlowToken.FindStringSubmatch(b)
	sctx := reSCtx.FindStringSubmatch(b)
	canary := reCanary.FindStringSubmatch(b)

	if len(flowToken) < 1 || len(sctx) < 1 || len(canary) < 1 {
		return nil, errors.New("Houston we have a problem with flowToken sctx canary")
	}

	return &MetaData{flowToken[1], sctx[1], canary[1]}, nil

}

func (mo *MonitoringAzureAD) GetDirSyncManagement() (*DirSyncManagement, error) {
	const (
		url             = "https://portal.office.com/admin/api/DirSyncManagement/manage"
		accept          = "Accept"
		applicationJSON = "application/json"
	)
	meta, err := mo.getMetaData()
	if err != nil {
		return nil, err
	}

	if errAuth := mo.authenticate(meta); errAuth != nil {
		return nil, errAuth
	}

	if errForceTheDoor := mo.openTheDoor(); errForceTheDoor != nil {
		return nil, errForceTheDoor
	}

	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header.Add(accept, applicationJSON)

	resp, err := mo.client.Do(r)
	defer closeBody(resp)
	if err != nil {
		return nil, err
	}

	dirSyncManagement := &DirSyncManagement{}
	if err := json.NewDecoder(resp.Body).Decode(dirSyncManagement); err != nil {
		return nil, err
	}
	return dirSyncManagement, nil
}

func (mo *MonitoringAzureAD) authenticate(meta *MetaData) error {
	const (
		urlValue = "https://login.microsoftonline.com/common/login"
	)
	params := url.Values{
		"login":     {mo.login},
		"passwd":    {mo.passwd},
		"canary":    {meta.canary},
		"ctx":       {meta.sctx},
		"flowToken": {meta.flowToken},
	}

	resp, err := mo.client.PostForm(urlValue, params)
	defer closeBody(resp)
	if resp.StatusCode != 200 {
		return errors.New("Houston authenticate is not 200")
	}
	return err
}

func (mo *MonitoringAzureAD) openTheDoor() error {
	const (
		urlValue = "https://portal.office.com/adminportal/home#/homepage"
	)
	resp, err := mo.client.Get(urlValue)
	defer closeBody(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Houston forceTheDoor is not 200")
	}
	findElement := func(n *html.Node, typeElement string) *html.Node {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == typeElement {
				return c
			}
		}
		return nil
	}

	nameAndValue := func(n *html.Node) (string, string) {
		var name = ""
		var value = ""
		for _, attInput := range n.Attr {
			if attInput.Key == "name" {
				name = attInput.Val
			}
			if attInput.Key == "value" {
				value = attInput.Val
			}
		}
		return name, value
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	htmlElement := findElement(doc, "html")
	if htmlElement == nil {
		return errors.New("Houston where is html?")
	}

	bodyElement := findElement(htmlElement, "body")
	if bodyElement == nil {
		return errors.New("Houston where is body?")
	}

	formElement := findElement(bodyElement, "form")
	if formElement == nil {
		return errors.New("Houston where is formElement?")
	}

	for _, att := range formElement.Attr {
		if att.Key == "action" {
			landingParams := url.Values{}
			for c := formElement.FirstChild; c != nil; c = c.NextSibling {
				if c.Data == "input" {
					name, value := nameAndValue(c)
					landingParams.Add(name, value)
				}
			}

			_, err = mo.client.PostForm(att.Val, landingParams)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("Houston nothing happened")
}
