package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/sjeandeaux/azure-ad-go/information"
	"github.com/sjeandeaux/azure-ad-go/log"
	"github.com/sjeandeaux/azure-ad-go/soap"
	"github.com/sjeandeaux/azure-ad-go/ugly"
)

//commandLineArgs all parameters in command line
type commandLineArgs struct {
	user         string
	password     string
	clientID     string
	tenantDomain string
	tenantID     string
	verbose      bool
}

var commandLine = &commandLineArgs{}

//init configuration
func init() {
	log.Logger.SetOutput(os.Stdout)
	flag.StringVar(&commandLine.user, "user", os.Getenv("AZURE_AD_USER"), "user for Azure AD")
	flag.StringVar(&commandLine.password, "password", os.Getenv("AZURE_AD_PASSWORD"), "user for Azure AD")
	flag.StringVar(&commandLine.clientID, "clientID", os.Getenv("AZURE_AD_CLIENT_ID"), "The client ID")
	flag.StringVar(&commandLine.tenantDomain, "tenantDomain", os.Getenv("AZURE_AD_TENANT_DOMAIN"), "Th tenant domain")
	flag.StringVar(&commandLine.tenantID, "tenantID", os.Getenv("AZURE_AD_TENANT_ID"), "The tenant ID")
	verbose, _ := strconv.ParseBool(os.Getenv("AZURE_AD_VERBOSE"))
	flag.BoolVar(&commandLine.verbose, "verbose", verbose, "verbose")
	flag.Parse()
}

func main() {
	log.Logger.Println(information.Print())
	uglyGet := ugly.NewMonitoringAzureAD(commandLine.user, commandLine.password)
	if dir, err := uglyGet.GetDirSyncManagement(); err != nil {
		log.Logger.Panicln(err)
	} else {
		log.Logger.Println(fmt.Sprintf("%+v", dir))
	}

	//soap
	soapGet := soap.NewMonitoringAzureAD(commandLine.user, commandLine.password, commandLine.clientID, commandLine.tenantDomain, commandLine.tenantID, commandLine.verbose)
	//aaaaaaaaaaaaahhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh!!!!!!!!!!!!!!!!!!!!!
	if token, err := soapGet.AccessToken(soap.Graph); err != nil {
		log.Logger.Panicln(err)
	} else {
		if dataBlob, err := soapGet.MsolConnect(token); err != nil {
			log.Logger.Panicln(err)
		} else {
			if datum, err := soapGet.GetCompanyInformation(token, dataBlob); err != nil {
				log.Logger.Panicln(err)
			} else {
				log.Logger.Printf("%+v", datum)
			}

			if datum, err := soapGet.HasObjectsWithDirSyncProvisioningErrors(token, dataBlob); err != nil {
				log.Logger.Panicln(err)
			} else {
				log.Logger.Printf("%+v", datum)
			}
		}
	}
}
