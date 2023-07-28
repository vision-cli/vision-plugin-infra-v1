TF_ENVIRONMENTS_DIR = config
TF_ARGS = -var-file=$(TF_ENVIRONMENTS_DIR)/$(ENVIRONMENT).tfvars -var="db_pass=$(POSTGRES_PASSWORD)" -var="auth_pass=$(AUTH_PASSWORD)"

init:
	@terraform init -backend-config="$(TF_ENVIRONMENTS_DIR)/$(ENVIRONMENT)-backend.conf.tmpl"

apply:
	@terraform apply $(TF_ARGS) $(TF_AUTO_APPROVE)

plan:
	@terraform plan $(TF_ARGS)


# Use to avoid pipeline hanging when vars not defined
plan_ci:
	@terraform plan --input=false $(TF_ARGS)

destroy:
	@TF_WARN_OUTPUT_ERRORS=1 terraform destroy $(TF_ARGS)
