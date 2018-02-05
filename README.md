# Monitor Azure AD

## Soap

* GetCompanyInformation

## ugly

The steps:
* GET https://login.microsoftonline.com to get canary, sctx and flowToken
* POST FORM https://login.microsoftonline.com/common/login with canary, sctx and flowToken, login and passwd
* GET https://portal.office.com/adminportal/home#/homepage
  * POST the homepage form
* GET https://portal.office.com/admin/api/DirSyncManagement/manage

References:
* https://gist.github.com/dlundgren/f4778c235eabd6467d6c5a9f727f9a7c
* https://github.com/outsideopen/nagios-plugins
