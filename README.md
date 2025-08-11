# CDM-CLI
The cdm-cli aims to create a smooth and intuitive development experience for engineers writing new checks. 
It is currently in early development and inclides one command `addCheck` which enables developers to quickly generate the template for a new check inside the [cdm-checks](https://github.com/IamSamD/cdm-checks) shared library repository. 

## Creating a new check
I will use an example of creating a new check that will check the upgrade status of an Azure AKS cluster, failing and notifying the SRE team when the cluster is due an upgrade. 

1. Clone the [cdm-checks](https://github.com/IamSamD/cdm-checks) repo if you do not have it already.

2. In the root of the repo:
```bash
./cdm-cli addCheck --provider azure --resource aks --check-name upgrade
```

This will scaffold a new check template in the repo under `azure/aks/upgrade`

```bash
.
└── azure
    └── aks
        └── upgrade
            ├── check
            │   └── check.go
            ├── go.mod
            ├── go.sum
            └── main.go
```

Now refer to the [cdm-checks documentation](https://github.com/IamSamD/cdm-checks/blob/main/README.md) for guidance on how to write the check

## Future Updates
In future updates we plan to add a command for scaffolding the initial repository and pipeline for implementing cdm-checks for a client.

We are also considering the design choice of the cdm-cli actiing as an orchestrator to run checks in the cdm pipeline.
This could abrstact away further configuration ingestoion and 'boostrapping' of a check before execution, keeping the checks themselves simplistic and focussed on a single purpose. 