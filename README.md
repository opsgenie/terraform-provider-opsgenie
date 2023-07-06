
Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://www.datocms-assets.com/2885/1629941242-logo-terraform-main.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.14.2 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/opsgenie/terraform-provider-opsgenie`

```sh
$ mkdir -p $GOPATH/src/github.com/opsgenie; cd $GOPATH/src/github.com/opsgenie
$ git clone git@github.com:opsgenie/terraform-provider-opsgenie
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/opsgenie/terraform-provider-opsgenie
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-opsgenie
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

Testing the Provider from Local Registry Version
------------------------------------------------
* Create a `.terraformrc` file on your in your home folder using `vi ~/.terraformrc`
* Add the local provider registry conf in `.terraformrc`
```
provider_installation {
  filesystem_mirror {
    path    = "~/terraform/providers"
    include = ["test.local/*/*"]
  }
  direct {
    exclude = ["test.local/*/*"]
  }
}
```
* Run `make build` on local (it will internally trigger a hook to write to `test.local` registry located in your `~/terraform/providers` folder)
* You can create a terraform basic project of your own locally with `main.tf` file
```
terraform {
  required_providers {
    opsgenie = {
      source  = "test.local/opsgenie/opsgenie"
      version = "<local_version>"
    }
  }
}

# Configure the Opsgenie Provider
provider "opsgenie" {
  api_key = <api_key>
  api_url = "api.opsgenie.com" # can be a stage instance url for devs
}
```
* And, Add respective terraform change files which you want to apply on your OG instance
* Run respective terraform commands to test the provider as per your convenience