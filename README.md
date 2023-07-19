# ![logo](./images/vision-logo.svg "Vision") &nbsp; Vision Plugin - Infra

This plugin creates a standard infra template

Vision plugins require golang (https://go.dev) to be installed

Install the plugin with

```
go install github.com/vision-cli/vision-plugin-infra-v1
```

You will now see the infra plugin commands on the vision cli

```
vision --help
```

Before running the plugin, you must set the following environment variables in your active terminal:
```
  --- Azure ---
  export AZURE_SUBSCRIPTION_ID=...

  --- Terraform ---
  export ENVIRONMENT=...
  export POSTGRES_PASSWORD=...
  export AUTH_PASSWORD=...
  export TF_AUTO_APPROVE=...

```

You are now ready to run the plugin. You can run the plugin using:

```
cat message.json | vision-infra-plugin-v1
```