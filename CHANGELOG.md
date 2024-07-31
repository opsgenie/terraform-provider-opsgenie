## 0.6.37 (July 30, 2024)
* BUGFIX: [#446](https://github.com/opsgenie/terraform-provider-opsgenie/pull/446)
  * **Integration Policy**
    * Fixed perpetual time drift in the conditions field in Integration policy
## 0.6.36 (July 7, 2024)
* BUGFIX:  [#411](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/411), [#440](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/440)
  * **API backoff mechanism**
    * Add options to customize the API backoff mechanism
  * **Integration Action Policy**
    * Fixed integration_action filter conditions field by giving support for extra-properties
* IMPROVEMENTS   [#296](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/296)
  * **Integration Policy**
    * Updated Integration request to allow configuration access flag
## 0.6.35 (December 18, 2023)
* BUGFIX: [#413](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/413), [#416](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/416) 
  * **time_restriction:**
    * Fixed drift for empty time_restriction
    * Fixed time_restriction argument in team_routing_rule resource documentation
* IMPROVEMENTS: [#405](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/405)
  * **README**
    * Update the README setup instructions and Make commands
## 0.6.34 (November 06, 2023)
* BUGFIX: [#404](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/404)
  * **Notification Policy:**
    * Fixed perpetual drift for policies when filters contain more than one condition.
  * **time_restriction:**
    * Fixed perpetual drift for `notification/alert policies`, `notification/team_routing rules` and `schedule_rotation` containing `time_restriction` blocks.
  * **Notification Rule:**
    * Fixed perpetual drift for rules when they contain more than one step.

* IMPROVEMENTS: [#404](https://github.com/opsgenie/terraform-provider-opsgenie/pulls/404)
  * **time_restriction:**
    * Added further schema validation to make it easier to type valid `time_restriction` blocks when using the `terraform-ls` language server.
  * **Notification Policy:**
    * Added further schema validation to make it easier to add multiple `action` blocks when using the `terraform-ls` language server.

## 0.6.28 (July 13, 2023)
* BUGFIX:
  * **API Integration:**
    * Fixes an issue where owner team could not be updated when the API integration is linked with an Integration action.

## 0.6.27 (July 11, 2023)
* IMPROVEMENTS:
  * **Alert Policy:**
    * Added support for `escalation` and `schedule` type in `responders` field.
  * **Dev Loop**
    * Added a hook in goreleaser to generate the local terraform binary/exe to ease debugging.

## 0.6.26 (June 29, 2023)
* BUGFIX:
  * **Notification Rule:**
    * Fixed an issue where the users could not set the send_after value to 0 while creating a notification rule.

## 0.6.25 (June 23, 2023)
* IMPROVEMENTS:
  * **Integration Action:** Updated documentation to show how to create alert integration action with multiline description using chomp function.
  * **Schedule** Updated documentation by removing unused `rules` attribute.
  * Updated `README` to correctly render Terraform logo.
* BUGFIX:
  * **Alert Policy:** 
    * Fixed ignore_original_actions and ignore_original_details being switched.
    * Added `{{description}}` as default value for `alert_description` field to solve [#290](https://github.com/opsgenie/terraform-provider-opsgenie/issues/290).
  * **Team:** Fixed error message to better explain the restrictions around team names.

## 0.6.24 (May 26, 2023)
IMPROVEMENTS:
* **Notification Policy:** Allow zeroes in until_hour/until_minute under notification_policy.delay_action

## 0.6.23 (May 26, 2023)
BUGFIX:
* Bump up opsgenie-go-sdk-v2 to v1.2.1.
* **Integration Action:** Added recipients to update in integration_actions responders.

## 0.6.22 (May 23, 2022)
IMPROVEMENTS:
* Bump up opsgenie-go-sdk-v2 to v1.2.18.
* **API Integration:** Added owner team update.

## 0.6.21 (May 23, 2023)
IMPROVEMENTS:
* Bump up opsgenie-go-sdk-v2 to v1.2.16.

## 0.6.20 (January 30, 2023)
BUGFIX:
* **Service Incident Rule:** Conditions is a set not an ordered list

## 0.6.19 (January 25, 2023)
BUGFIX:
* **Team Routing Rule:** Stop sending routing rule order update if it's a default rule
* **Escalation:** Fix updating escalation's owner team issue when using team based integration api key

## 0.6.18 (November 16, 2022)
BUGFIX:
* **Escalation:** Import repeat field bug fixed.

## 0.6.11 (February 18, 2022)
* fix team set problem


## 0.6.10 (February 18, 2022)
* Update Routing rule will update the order,too.
  **Note:** There are still bug in creation of Routing Rules from scratch due to concurrency problem. Following command will solve the problem. Apologize for the confusion
  
            terraform apply -refresh-only
  

## 0.6.8 (January 5, 2022)
* Alert Policy alert_description fix implemented

## 0.6.7 (December 13, 2021)
* GoLang version increased

## 0.6.6 (December 13, 2021)
BUGFIX:
* **Routing Rule:** Add is_default field
* **Heartbeat:** Add "." string support

## 0.6.5 (June 15, 2021)
BUGFIX:
* **Schedule:** Timezone diff problem fixed.

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
