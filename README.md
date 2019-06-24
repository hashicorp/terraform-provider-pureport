Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Pureport Inc, team at [Pureport](https://www.pureport.com).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Usage
---------------------

```
# For example, restrict pureport version in 0.4.x
provider "pureport" {
  version = "~> 0.4"
}
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-pureport`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-pureport
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-pureport
$ make build
```

Using the provider
----------------------

## 3rd Party plugin installation

Copy the terraform-provider-pureport plugin in to the terraform third-party plugins directory.

| OS                | Location                        |
|-------------------|:--------------------------------|
| Windows           | %APPDATA%\terraform.d\plugins   |
| All other systems | ~/.terraform.d/plugins          |

More information about this can be found [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*).
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-pureport
...
```

This provider uses `golangci-lint` for checking static analysis of the source code. This needs to be
installed separate from the other golang modules required to build the provider.

```sh
$ make tools
$ make lint
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

You can also install the plugin which will build and copy the plugin to your terraform third party
plugin directory. You'll need to re-initialize terraform in module directory after installing the
new plugin.

```sh
$ make install
$ cd <some_module>/
$ terraform init
```

Acceptance Test Setup
---------------------------

When preparing to run the acceptance tests, some initial manual setup will be required for each
cloud provider to ensure we are able to deploy the cloud infrastructure for testing.

An example environment setup script is available in this repository in
`examples/envsetup.sh.examples`. You can modify this file with your cloud provider information
and then source it in to your shell prior to deploying and running the acceptance tests.

After the credentials have been setup, you can run Terraform Configuration in `test-infra` to deploy
the required cloud provider resources.

## Azure

For Azure, you will need to create a Resource Group with the name "terraform-acceptance-tests" and
also a Service Principle as instructed by the `azurerm` provider. Instructions can be found [here](https://www.terraform.io/docs/providers/azurerm/auth/service_principal_client_secret.html).

After running the test infrastructure, please copy the service key from the output to your
environment setup script. This will be needed for the acceptance tests.

## Google Cloud

For Google Cloud, you will need to have a valid account and a project created that you can deploy
resource in to.

## AWS

For Amazon Web Services, you will need to have a valid IAM identity with permission to create
resources.

