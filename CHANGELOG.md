## 0.6.4 (April 14, 2021)
BUGFIX:
* **Notification Policy:** De-duplication Action problem fixed.

## 0.6.3 (February 11, 2021)
BUGFIX:
* **Alert Policy:** Global alert policy import fixed

## 0.6.2 (January 29, 2021)
BUGFIX:
* **User:** Timezone diff problem fixed.

## 0.6.1 (January 21, 2021)
**BREAKING CHANGE**

**Terraform Plugin SDK upgraded to v2.**

BUGFIX:
* **Notification Rule :** Reading the time_restrictions corrected.
* **Team Routing Rule :** Reading the time_restrictions corrected.
* **Schedule Rotations :** Reading the time_restrictions corrected.

## 0.6.0 (January 19, 2021)
**BREAKING CHANGE**

Terraform Plugin SDK upgraded to v2. 
Acceptance tests are need Terraform 12.26+ versions.

BUGFIX:
* **Notification Rule :** Reading the time_restrictions corrected.
* **Team Routing Rule :** Reading the time_restrictions corrected.
* **Incident Service Rule :** Tags and details are implemented. Before it won't work due to schema <-> List conversion.
* **Incident Template:** Stakeholder properties and impacted services read functions fixed to comply with their schema types.
* **Integrations:** Responders only available when owner team isn't set. Therefore provider now only read and add responders to requests if owner team is not available.
* **User:** UserAddress.city field has fixed. In future user resource will change to adopt Atlassian Opsgenie Platform changes.

## 0.5.7 (December 24, 2020)
IMPROVEMENTS:
* **Integration Actions :** add support for ignore action


## 0.5.6 (December 24, 2020)
IMPROVEMENTS:
* **Integration Actions :** add custom priority for create action (#177)


## 0.5.5 (December 16, 2020)
Improvement:
* Add support for webhook integration (#197)


## 0.5.4 (December 10, 2020)
BUGFIX:
* Able to set key in extra-properties field of condition for service incident rule (#204)


## 0.5.3 (December 4, 2020)
BUGFIX:
* GO SDK v2 version synced to support none & escalation in schedule rotations 
* Docs update


## 0.5.2 (November 19, 2020)
BUGFIX:
* Fix opsgenie_notification_rule: Fix for issue (#188) 
* Update integration name restrictions (#187)
* Docs update (#194)


## 0.5.1 (October 16, 2020)
IMPROVEMENTS:
* Added missing options to user resource (#179-#180)
* Added actions to opsgenie_alert_policy (#186)


## 0.5.0 (September 18, 2020)
NEW RESOURCE:
* Custom user role implemented (#161)

## 0.4.9 (September 17, 2020)
NEW RESOURCE:
* Incident Template implemented (#178)


## 0.4.8 (September 4, 2020)
NEW RESOURCE:
* Notification rule implemented (#121)


## 0.4.7 (August 26, 2020)

IMPROVEMENTS:
* **Team :** allow users to delete default resources while creating team.


## 0.4.6 (August 21, 2020)

BUGFIX:
* **Integration Actions :** allow integration action import(#151)

## 0.4.5 (August 20, 2020)

IMPROVEMENTS:
* **Integration Actions :** Go-Sdk-v2 updated to support all custom field names.

## 0.4.4 (August 18, 2020)

IMPROVEMENTS:
* **Integration Actions :** filter conditions set eventType as field (#148)

## 0.4.3 (August 14, 2020)

IMPROVEMENTS:
* **Integration Actions :** add priority for create action (#157)

## 0.4.2 (August 13, 2020)

BUGFIX:
* **Integration Actions :** allow extra_properties (#152)

## 0.4.1 (July 29, 2020)

Opsgenie Provider repository changes

## 0.4.0 (July 28, 2020)

IMPROVEMENTS:

* **Integration Actions Api :** New resource, now you can manage integration actions via Terraform. (#139)

BUGFIX:

* **Team API :** Allow dot character (#137)

## 0.3.9 (July 20, 2020)

BUGFIX:

* **Alert Policy:** Fixed TF crash because of field type (#132)

## 0.3.8 (July 16, 2020)

IMPROVEMENTS:

* **Service Incident Rule Api :** New resource, now you can manage incident rules via Terraform. (#130)

BUGFIX:

* **Team Routing Rule:** Fixed not field and edited condition validation (#125)

## 0.3.7 (July 13, 2020)

IMPROVEMENTS:

* **Alert Policy API:** New resource, now you can manage alert policies.

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
