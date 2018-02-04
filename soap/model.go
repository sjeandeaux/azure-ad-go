package soap

import "time"

type MsolConnect struct {
	ClientId       string
	BearerToken    string
	TrackingHeader string
	MessageID      string
}
type MsolConnectEnvelope struct {
	Header *MsolConnectHeaderResponse `xml:"Header"`
}

type MsolConnectHeaderResponse struct {
	BecContext *MsolConnectBecContextResponse `xml:"BecContext"`
}

type MsolConnectBecContextResponse struct {
	DataBlob string `xml:"DataBlob"`
}

type GetCompanyInformation struct {
	DataBlob       string
	BearerToken    string
	TrackingHeader string
	MessageID      string
}

type GetCompanyInformationEnvelope struct {
	Body GetCompanyInformationBody `xml:"Body"`
}

type GetCompanyInformationBody struct {
	Response GetCompanyInformationResponse `xml:"GetCompanyInformationResponse"`
}

type GetCompanyInformationResponse struct {
	Result GetCompanyInformationResult `xml:"GetCompanyInformationResult"`
}

type GetCompanyInformationResult struct {
	ReturnValue GetCompanyInformationReturnValue `xml:"ReturnValue"`
}

type GetCompanyInformationReturnValue struct {
	LastDirSyncTime time.Time `xml:"LastDirSyncTime"`
}

type HasObjectsWithDirSyncProvisioningErrors struct {
	DataBlob       string
	BearerToken    string
	TrackingHeader string
	MessageID      string
}

type HasObjectsWithDirSyncProvisioningErrorsEnvelope struct {
	Body HasObjectsWithDirSyncProvisioningErrorsBody `xml:"Body"`
}

type HasObjectsWithDirSyncProvisioningErrorsBody struct {
	Response HasObjectsWithDirSyncProvisioningErrorsResponse `xml:"HasObjectsWithDirSyncProvisioningErrorsResponse"`
}

type HasObjectsWithDirSyncProvisioningErrorsResponse struct {
	Result HasObjectsWithDirSyncProvisioningErrorsResult `xml:"HasObjectsWithDirSyncProvisioningErrorsResult"`
}

type HasObjectsWithDirSyncProvisioningErrorsResult struct {
	ReturnValue bool `xml:"ReturnValue"`
}
