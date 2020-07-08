## 0.4.0 (Unreleased)
## 0.3.6 (July 08, 2020)

IMPROVEMENTS:

* **Data Source: Service API:** New datasource, now you can manage service which created without using Terraform. (#118)


## 0.3.5 (July 05, 2020)

IMPROVEMENTS:

* **Service API:** New resource, now you can manage service resources using Terraform. (#115)

BUGFIX:

* **Api Integration:** Api integration update will no longer resets fields, which not managed via Terraform. (#119)

## 0.3.4 (May 20, 2020)

BUGFIX:

* **Api Integration:** Read function fixes the Allow_Write_Access field. Also default behaviour of `true` implemented for this field.

## 0.3.3 (May 18, 2020)
IMPROVEMENTS:

* If resource deleted manually, Provider find out while reading resources then approach to re-creates resources 

## 0.3.2 (May 14, 2020)

BUGFIX:

* **Team Routing Rule:** Some fields was optional but provider expects its to be mandatory. pr/96 fixes this.


## 0.3.1 (April 16, 2020)

BUGFIX:

* **Schedule Rotation:** Time restriction reading causes provider crash. This bug introduced in 0.3.0 version its fixed this release. 


## 0.3.0 (April 10, 2020)

BUGFIX:

* **Maintenance:** Fixed edit maintenance end date.
* **Integrations:** Added missing fields for states (owner_team_id, allow_configure_access).
* **Schedule Rotation:** Fix import https://github.com/terraform-providers/terraform-provider-opsgenie/pull/88 


## 0.2.9 (March 23, 2020)

IMPROVEMENTS:

* **Show warning if Opsgenie username (email addr) contains uppercase characters:** This lead to unexpected behaviour in the past.
* **Updated Resource opsgenie_team:** New optional argument *ignore_members* added to change team membership management behaviour (#65). The provider will add this argument to every new/existing opsgenie_team resource state with the default value (false).  

BUGFIX:

* **Updated documentations** 
* **Added missing resource fields in schedule** 


## 0.2.8 (February 07, 2020)

IMPROVEMENTS:

* **Updated opsgenie-go-sdk-v2** 

* **New Resource:** Notification Policy

BUGFIX:

* **Updated documentations** 


## 0.2.7 (January 21, 2020)

BUGFIX:

* Importing with parent resources fixed. It can be imported using parentID/resourceID syntax through cli.

## 0.2.6 (December 19, 2019)

IMPROVEMENTS:

* **Updated opsgenie-go-sdk-v2** 

* **New Resource:** Team Routing Rule

BUGFIXES:

* 'Random' NotifyType added.


## 0.2.5 (November 19, 2019)

IMPROVEMENTS:

* **Updated opsgenie-go-sdk-v2** 

BUGFIXES:

* OwnerTeam field added to EmailBasedIntegrations https://github.com/terraform-providers/terraform-provider-opsgenie/issues/40

## 0.2.4 (November 04, 2019)

IMPROVEMENTS:

* **Updated opsgenie-go-sdk-v2** 
* **Migrated to terraform-plugin-sdk** 

BUGFIXES:

* Global heartbeat creation fixed.

## 0.2.3 (September 25, 2019)

IMPROVEMENTS:

You can now refer existing resources on Opsgenie using datasources.

* **New Datasources:**  Schedule, Escalation, Heartbeat, Team, User
* **Edited Resources:** Schedule, Schedule Rotations, Team and User Contact

BUGFIXES:

* Some resources states are fixed.
* Edited date validations
* Test improvements

## 0.2.2 (September 18, 2019)

IMPROVEMENTS:

* **Edited Resource:** You can refer api based integrations api key to external resources by `.api_key` field. 

BUGFIXES:

* Documents updated

## 0.2.1 (September 13, 2019)

IMPROVEMENTS:

* **New Resource:** Heartbeat

BUGFIXES:

* Documents updated

## 0.2.0 (September 05, 2019)

IMPROVEMENTS:

* All resources updated using Opsgenie Go SDK v2
* User API
* Team API
* **New Resource:** User Contact API
* **New Resource:** Integration API (API and Email based)
* **New Resource:** Escalation API
* **New Resource:** Schedule API
* **New Resource:** Schedule Rotation API
* **New Resource:** Maintenance API

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
