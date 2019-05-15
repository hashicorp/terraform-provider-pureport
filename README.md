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

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Usage
---------------------

```
# For example, restrict pureport version in 0.1.x
provider "pureport" {
  version = "~> 0.1"
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

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
