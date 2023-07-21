package plugin

import (
	"errors"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/common/execute"
	"github.com/vision-cli/common/marshal"
	"github.com/vision-cli/common/tmpl"
	"github.com/vision-cli/vision-plugin-infra-v1/placeholders"
	"github.com/vision-cli/vision-plugin-infra-v1/run"
)

var Usage = api_v1.PluginUsageResponse{
	Version: "0.1.0",
	Use:     "gcp",
	Short:   "manage gcp infratrusture",
	Long:    "manage gcp infrastructure using a standard template",
	Example: `Before runing vision gcp create, you will need a gcp seed project. The purpose of this project is to:

- store the terraform state in a storage bucket so it can be shared
- hold the service account that is used to create the dev and prod projects
- hold the container registry
- hold the oauth2 client, because Google cant automate this, it must be done manually before an IAP can be created

You should already have this seed project from project create as you need the container registry parameter to create the project. 

1. In the seed project:

   1.1. Create a storage bucket called 'tfstate-<unique id>', 
        where <unique id> is a unique-str for your project - you can see this in vision.json.
        Set the region as required.
        Warning! It is highly recommended that you enable Object Versioning on the GCS bucket 
		         to allow for state recovery in the case of accidental deletions and human error.

   1.2. In the Google console go to the oauth consent screen and set it to "external". 
        You will have to create a brand as well. Make a note of the client id and secret (you will need this for step 4).

2. Now you can run:
  
      vision gcp create

3. Update the dev.tfvars and prod.tfvars files in infra/gcp/platform/tf/config. 
   You need to update details such as billing account, folder, etc.

4. The first infrastructure creation must be done manually in order to create the workload identity used by the github workflow.
 
   4.1 Make sure you are logged in and authe'ed against the seed project

          gcloud auth application-default login --project seed

   4.2 From the /infra/gcp/platform/tf folder initialise terraform with

          ENVIRONMENT=dev make init
  
   4.3. Apply the terraform with

          ENVIRONMENT=dev DB_PASSWORD=<dev db password> OAUTH2_CLIENT_SECRET=<oauth client secret from 1.2> make apply
  
   4.4. Configure DNS with the load balancer IPs returned by terraform, viz

3. Create a PR to merge the code and github workflow into your repo. Once merged, 
   the workflow will run and create the dev and prod projects.
   Make a note of the terraform output variables which you need for step 4 below.

4. Create the following secrets in the project's github repo:
   - DEV_DB_PASSWORD - dev db password from 4.3 above
   - DEV_GCP_WORKLOAD_IDP - this is the idp for the dev workload (from step 4.3)
   - DEV_GCP_SERV_ACCOUNT - this is the service account for the dev workload (from step 4.3)
   - PROD_DB_PASSWORD - prod db password from 4.3 above
   - PROD_GCP_WORKLOAD_IDP - this is the idp for the prod workload (from step 4.3)
   - PROD_GCP_SERV_ACCOUNT - this is the service account for the prod workload (from step 4.3)
   - OAUTH2_CLIENT_SECRET - this is the oauth2 client secret from the seed project (step 1.2 above)
`,
	Subcommands:    []string{"create"},
	Flags:          []api_v1.PluginFlag{},
	RequiresConfig: true,
}

var DefaultConfig = api_v1.PluginConfigResponse{
	Defaults: []api_v1.PluginConfigItem{},
}

func Handle(input string, e execute.Executor, t tmpl.TmplWriter) string {
	req, err := marshal.Unmarshal[api_v1.PluginRequest](input)
	if err != nil {
		return errorResponse(err)
	}
	result := ""
	switch req.Command {
	case api_v1.CommandUsage:
		result, err = marshal.Marshal[api_v1.PluginUsageResponse](Usage)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandConfig:
		result, err = marshal.Marshal[api_v1.PluginConfigResponse](DefaultConfig)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandRun:
		if len(req.Args) == 0 || req.Args[placeholders.ArgsCommandIndex] == "" {
			return errorResponse(errors.New("missing cli command"))
		}
		switch req.Args[placeholders.ArgsCommandIndex] {
		case "create":
			p, err := placeholders.SetupPlaceholders(req)
			if err != nil {
				return errorResponse(err)
			}
			err = run.Create(p, e, t)
			if err != nil {
				return errorResponse(err)
			}
		default:
			return errorResponse(errors.New("unknown cli command"))
		}
		resp := api_v1.PluginResponse{
			Result: "SUCCESS!",
			Error:  "",
		}
		result, err = marshal.Marshal[api_v1.PluginResponse](resp)
		if err != nil {
			return errorResponse(err)
		}
	default:
		return errorResponse(errors.New("unknown command"))
	}
	return result
}

func errorResponse(err error) string {
	res, err := marshal.Marshal[api_v1.PluginResponse](api_v1.PluginResponse{
		Result: "",
		Error:  err.Error(),
	})
	if err != nil {
		panic(err.Error())
	}
	return res
}
