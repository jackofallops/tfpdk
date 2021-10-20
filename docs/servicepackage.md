# ServicePackage

The `tfpdk servicepackage` command is used to create a boilerplate structure for a new Service in the provider.

Example:
```shell
tfpdk servicepackage -servicepackage=MyNewServicePackage
```
will create the following structure:
```shell
.
└── internal
    └── services
        └── MyNewServicePackage
            ├── client
            │   └── client.go
            └── registration.go
```

* `client.go` is used for instantiating Service Clients to use for operations within the resource CRUD functions. The file created contains a basic structure to guide this.

* `registration.go` is used to declare the resources and data sources made available by this service package. The code within is a basic, compliant and compilable framework into which resources and Data Sources can be defined.
~> **Note:** A future version of this tool will be able to manage this file on the user's behalf, but for now it must be manually updated.