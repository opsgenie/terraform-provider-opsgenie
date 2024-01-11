
# Terraform Provider

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://www.datocms-assets.com/2885/1629941242-logo-terraform-main.svg" width="600px">


## Development

Everything needed to start local development and testing of the Provider plugin

### 1. Requirements

-	[Go](https://golang.org/doc/install) 1.18 (or higher, to build the provider plugin)
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x (To test the plugin)


### 2. Development Setup

#### Cloning the project

```bash
export GOPATH="${GOPATH:=$HOME/go}"

mkdir -p "$GOPATH/src/github.com/opsgenie"

cd "$GOPATH/src/github.com/opsgenie"

git clone git@github.com:opsgenie/terraform-provider-opsgenie
```


### 3. Compiling The Provider

```sh
export GOPATH="${GOPATH:=$HOME/go}"
cd "$GOPATH/src/github.com/opsgenie"

# Compile all versions of the provider and install it in GOPATH.
make build

# Only compile the local executable version (Faster)
make dev
```


#### Running tests on the Provider

Run all local unit tests.

```sh
make test
```


### 4. Using the Compiled Provider

#### Configure Terraform to use the compiled provider.

This configuration makes use of the [`dev_overrides`](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).

See the [Manual Setup][Manual Setup] for more details.

```bash
# Creates $HOME/.terraformrc
make setup
```

#### Manual setup

<details>
  <summary>Click to expand</summary>

1. Create the `.terraformrc` file on your in your home folder using `touch ~/.terraformrc`
2. Configure [`dev_overrides`](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) in your `~/.terraformrc` as show below:
    ```hcl
    provider_installation {
      dev_overrides {
        # Remember to replace <home dir> with your username
        "opsgenie/opsgenie" = "/home/<home dir>/go/bin"
      }
      direct {}
    }
    ```
3. Run `make build`

</details>


#### New OpsGenie Terraform project

1. Create a basic terraform project locally with a `main.tf` file:
    ```hcl
    terraform {
      required_providers {
        opsgenie = {
          source  = "opsgenie/opsgenie"
          version = ">=0.6.0" # version can be omitted
        }
      }
    }

    # Configure the Opsgenie Provider
    provider "opsgenie" {
      api_key = "<insert api_key>" # https://support.atlassian.com/opsgenie/docs/api-key-management/
      api_url = "api.opsgenie.com" # can be a stage instance url for devs
    }

    resource "opsgenie_team" {
      name        = "Dev-Provider test team"
      description = "New team made using in-development OpsGenie provider"
    }
    ```

2. And, Add respective terraform change files which you want to apply on your OG instance

3. Run respective terraform commands to test the provider as per your convenience
   Install the currently available provider with `tf init`
   `terraform plan` and `terraform init` will use providers from the configured paths in `$HOME/.terraformrc`
   `terraform` will output an error if no provider is found in the `dev_overrides` path. (`make build`)


#### Removing the 'dev_override' again

This allows you to use the normal release versions of the `opsgenie/opsgenie` provider.

*Note* Removes `$HOME/.terraformrc`

```bash
make clean
```
