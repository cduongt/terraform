Terraform
=========

Fork for OCCI development.  
Current example infrastructure file is "metacloud-occi.tf" in main directory.  

Usage  
--------------------
1. Install Go (1.6+)
2. Set [`GOPATH`](https://golang.org/doc/code.html#GOPATH)
3. Clone this repository into `$GOPATH/src/github.com/hashicorp/terraform`
4. Install [`rocci-cli`](https://github.com/gwdg/rOCCI-cli)
5. Compile occi provider in repository

```sh
$ make plugin-dev PLUGIN=provider-occi
```
6. Edit example .tf file, especially path to voms proxy file and public key for contextualization
7. Run terraform

```sh
$ terraform apply
```
8. After you're done, destroy resources

```sh
$ terraform destroy
```
