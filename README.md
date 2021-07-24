# tfpdk
Terraform Provider Development Kit

NOTE: hard coded to [terraform-provider-azurerm](https://github.com/terraform-providers/terraform-provider-azurerm) currently... 

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