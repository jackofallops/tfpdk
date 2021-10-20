# Configuration File

Create using:
```shell
tfpdk config
```

Which will generate the `.tfpdk.hcl` file in the root of the project, containing:

```hcl
service_packages_path           = "internal/services"
provider_github_org             = "hashicorp"
docs_path                       = "docs"
resource_docs_directory_name    = "r"
data_source_docs_directory_name = "d"
use_typed_sdk                   = false
```

These default values can be updated as required for the provider under development and will not be overwritten by re-running the command again.  To return to default value, simply delete the file and re-run.
