# tfpdk
Terraform Provider Development Kit

NOTE: hard coded to [terraform-provider-azurerm](https://github.com/terraform-providers/terraform-provider-azurerm) currently... 

Run from the root of the provider project.

## TODO
- [x] untyped resource template
- [ ] `go fmt` outputs
- [ ] Add new resources / data sources to appropriate registration
- [x] typed and untyped Data Sources
- [ ] init a new provider - git clone [scaffold](https://github.com/hashicorp/terraform-provider-scaffolding)?
- [ ] optionally populate the Typed SDK into the init step? 
- [ ] Dynamically read provider name for item generation
- [ ] Populate `IDValidationFunc()` in template for IDs (Pandora)
- [ ] Clients?
- [ ] Autocomplete?


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

