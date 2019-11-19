## 0.2.5 (Unreleased)

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
