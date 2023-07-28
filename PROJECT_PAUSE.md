# Project Pause

**Pause Date:** 28/07/2023 

## What has been completed

The project currently creates a storage account in a specified resource group that is manually made. It then creates a Terraform storage container in the storage account that Terraform uses to store the state file. 

Terraform runs its `init` process a completes it successfully.

Terraform also runs its `apply` process but fails due to variables not being declared.

![terraform apply but failing](<Screenshot 2023-07-28 at 14.07.55.png>)

Template files have been moved from the `azure` package to the `run` package. This should allow the `.tmpl` files inside the `run/_templates/` package to create the Terraform files in the resulting project created by Vision. This hasn't been tested yet.

## What is next

1. Test the automated creation of the Terraform files using Vision. This will require placeholders. Check out the GCP project.

2. Try to automate the initialisation of the project. This includes:
    + setting all the environment variables 
        - Could this be done by a configuration file?
    + creating a "seed project"/"parent resource group"?
    + Decide whether all this is worth doing or it is just adding more steps to an already heavily manual process


