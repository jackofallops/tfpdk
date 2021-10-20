# tfpdk
Terraform Provider Development Kit

DISCLAIMER: The author primarily works on the [AzureRM Provider](https://github.com/hashicorp/terraform-provider-azurerm) so this tool is *heavily* biased towards the code styles and requirements of that provider currently.  Future development will look to genericise the output (eventually 😜) 
 
NOTE: Expects to be run from the root of a validly named Terraform provider e.g. `./terraform-provider-myprovider`

NOTE: Relies on a locally installed terraform binary and appropriately configured [dev overrides file](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers). This 
allows the Terraform binary to skip the `init` step by informing it where your compiled provider binary can be found. (you still currently need a local .tf file with your [provider configured](https://www.terraform.io/docs/language/providers/configuration.html) though, sorry!)

Run from the root of the provider project.

## TODO
- [x] untyped resource template
- [x] `go fmt` outputs
- [ ] Add new resources / data sources to appropriate registration
- [x] typed and untyped Data Sources
- [x] init a new provider - git clone [scaffold](https://github.com/hashicorp/terraform-provider-scaffolding)?
- [ ] optionally populate the Typed SDK into the init step? 
- [x] Dynamically read provider name for item generation
- [ ] Populate `IDValidationFunc()` in template for IDs (Pandora)
- [ ] Clients?
- [ ] Autocomplete?
- [x] Spike doc generation (just resources for now)
- [ ] migrate to plugin sdk to read schema directly for doc gen as Terraform's JSON output lacks necessary detail. (e.g. timeouts, ForceNew flags etc)
- [x] update commands and templates to allow non-hashicorp providers ~~and sources other than github~~ (e.g. import paths)


## Commands

* `tfpdk resource` - Creates resource,

## Usage Examples

### Create a Typed Data Source for an existing Resource where the Resource's model is appropriate for use in the Data Source
```shell
tfpdk datasource -name ShinyNewService -servicepackage SomeCloudService -typed -useresourcemodel
```
Will create the following path `{providername}/internal/services/SomeCloudService/shiny_new_service_data_source.go`

### Create an un-typed (traditional) resource 
```shell
tfpdk resource -name ShinyNewService -servicepackage SomeCloudService
```
Will create the following path `{providername}/internal/services/SomeCloudService/shiny_new_service_resource.go`

### Create a typed updatable resource 
```shell
tfpdk resource -name ShinyNewService -servicepackage SomeCloudService -has-update -typed
```
Will create the following path `{providername}/internal/services/SomeCloudService/shiny_new_service_resource.go`

### Document your new resource (relies on the `Description` fields in your schema, so populate those for best effect)
```shell
tfpdk document -type resource -name ShinyNewService -id "00000000-0000-0000-0000-000000000000"
```

## Command Documentation
`tfpdk [command]`
* [init](docs/init.md)
* [config](docs/config.md)
* [servicepackage](docs/servicepackage.md)
* [resource](docs/resource.md)
* [datasource](docs/datasource.md)
* [document](docs/document.md)