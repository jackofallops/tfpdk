# tfpdk
Terraform Provider Development Kit
 
NOTE: Expects to be run from the root of a validly named Terraform provider e.g. `./terraform-provider-myprovider`

~~WARNING: currently needs the json output from `terraform providers schema -json` for the azurerm provider in `/tmp/azurerm-provider-out.json`~~ 
Relies on a locally installed terraform binary and appropriately configured [dev overrides file](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers). This 
allows the Terraform binary to skip the `init` step by informing it where your compiled provider binary can be found.

Run from the root of the provider project.

## TODO
- [x] untyped resource template
- [x] `go fmt` outputs
- [ ] Add new resources / data sources to appropriate registration
- [x] typed and untyped Data Sources
- [ ] init a new provider - git clone [scaffold](https://github.com/hashicorp/terraform-provider-scaffolding)?
- [ ] optionally populate the Typed SDK into the init step? 
- [ ] Dynamically read provider name for item generation
- [ ] Populate `IDValidationFunc()` in template for IDs (Pandora)
- [ ] Clients?
- [ ] Autocomplete?
- [x] Spike doc generation (just resources for now)
- [ ] update commands and templates to allow non-hashicorp providers and sources other than github (e.g. import paths)


## Commands

* `tfpdk resource` - Creates resource,

## Usage Examples

### Create a Typed Data Source for an existing Resource where the Resource's model is appropriate for use in the Data Source
```shell
tfpdk datasource -name=ShinyNewService -servicepackage=SomeCloudService -typed -useresourcemodel
```
Will create the following path `{providername}/internal/services/SomeClodService/shiny_new_service_data_source.go`

### Create an un-typed (traditional) resource 
```shell
tfpdk resource -name=ShinyNewService -servicepackage=SomeCloudService
```
Will create the following path `{providername}/internal/services/SomeClodService/shiny_new_service_resource.go`

### Create a typed updatable resource 
```shell
tfpdk resource -name=ShinyNewService -servicepackage=SomeCloudService -has-update -typed
```
Will create the following path `{providername}/internal/services/SomeClodService/shiny_new_service_resource.go`

