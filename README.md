Terraform
=========

Fork for OCCI development.
Documentation can be found [here](https://github.com/cduongt/terraform/blob/occi/website/source/docs/providers/occi/r/virtual_machine.html.markdown).

Prerequisities
--------------------
Working VOMS proxy file(/tmp/x509up_u1000)
[Contextualisation file](https://wiki.egi.eu/wiki/FAQ10_EGI_Federated_Cloud_User#Contextualisation)

Installation
--------------------
1. Install Go (1.7+)
2. Set [`GOPATH`](https://golang.org/doc/code.html#GOPATH)
3. Clone this repository into `$GOPATH/src/github.com/hashicorp/terraform`
4. Install [`rocci-cli`](https://github.com/gwdg/rOCCI-cli)
5. Change directory to Terraform repo
6. Change branch to current development branch (`git checkout occi`)
7. Compile Terraform core
```sh
$ make core-dev
```
8. Compile occi provider
```sh
$ make plugin-dev PLUGIN=provider-occi
```

Usage
--------------------

1. Create .tf infrastructure file
2. Run terraform
```sh
$ terraform apply
```
3. You can check created VMs
```sh
$ terraform show
```
4. After you're done, destroy resources
```sh
$ terraform destroy
```
