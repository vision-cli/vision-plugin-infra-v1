TF_CONFIG_DIR=config
TF_ARGS = -var-file=$(TF_ENVIRONMENTS_DIR)/$(ENVIRONMENT).tfvars \
          -var="db_user_password=$(DB_PASSWORD)" \
          -var="oauth2_client_secret=$(OAUTH2_CLIENT_SECRET)"

init:
	@terraform init -backend-config="$(TF_ENVIRONMENTS_DIR)/$(ENVIRONMENT)-backend.conf"

apply:
	@terraform apply $(TF_ARGS) $(TF_AUTO_APPROVE)

plan:
	@terraform plan $(TF_ARGS)


# Use to avoid pipeline hanging when vars not defined
plan_ci:
	@terraform plan --input=false $(TF_ARGS)

destroy:
	@TF_WARN_OUTPUT_ERRORS=1 terraform destroy $(TF_ARGS)
