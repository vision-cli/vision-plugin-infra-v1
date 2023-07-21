# Cloud Infrastructure

## Deployed Infrastructure

The platform consists of:

- Terraform code to deploy the following infra components

```
   |-----------------------------------------------------|
   |                  Https Load Balancer                |     // Https ingress to the platform
   |-----------------------------------------------------|
         |                              |
         |                 |------------------------|
         |                 |  Identity Aware Proxy  |          // Authentication proxy
         |                 |------------------------|
         |                              |
         |           |----------| |----------| |----------|
         |           |   NEG    | |    NEG   | |    NEG   |     // Network endpoint groups for cloud run containers
         |           |----------| |----------| |----------|
         |                 |            |            |
   |----------|      |----------| |----------| |----------|
   | cloud    |      | cloud    | | cloud    | | cloud    |
   | run      |      | run      | | run      | | run      |    // Cloud functions for platform services
   | gcs-proxy|      | gcs-proxy| | graphql  | | (various)|    // Only egress to VPN
   |----------|      |----------| |----------| |----------|
         |                     |                 |
  |-------------|       |-------------|   |-------------|    |---------|
  |   Storage   |       |   Storage   |   |   private   |----| Jumpbox |  // VPN jumpbox
  |   (public)  |       |  (private)  |   |   network   |    |---------|
  |-------------|       |-------------|   |-------------|----VPC peering  // VPC peering IP
                                                 |
                                          |-------------|
                                          |   Database  |
                                          |  (private)  |
                                          |-------------|

```

## Pre-build

Before runing terraform, you will need a terraforn seed project. The purpose of this project is to:

- store the terraform state in a storage bucket so it can be shared
- hold the service account that is used to create the dev and prod projects
- hold the container registry
- hold the oauth2 client, because Google cant automate this, it must be done manually before an IAP can be created

1.  Create a gcp project called 'seed'
    Typically create this in a folder for your project. The folder will evenutally contain your seed project as well as dev and prod
2.  Update config/org.tfvars with relevant details such as folder id, billing acct, etc

3.  Use the helper script from https://github.com/terraform-google-modules/terraform-google-project-factory.git to create a service account using the command below

    ```
    ./helpers/setup-sa.sh -o <organization id> -p <project id> [-b <billing account id>] [-f <folder id>] [-n <service account name>]
    ```

    Copy the resulting credentials.json to the terraform config folder
    Make a note of the service account email, you will need to provide it later to the github actions runner

You need the container registry (artifact registory in google speak). Use this to create the project

4.  Make sure you are logged in and authe'ed against the seed project

    ```
    gcloud auth application-default login --project seed
    ```

5.  Create a storage bucket called 'tfstate-&lt;unique id>'.
    Typically use the unique id from your project name.
    Set the region as required.
    Warning! It is highly recommended that you enable Object Versioning on the GCS bucket to allow for state recovery in the case of accidental deletions and human error.

6.  Update the backend.conf in the config folder to use the bucket just created

7.  In the Google console go to the oauth consent screen and set it to "external". You will have to create a brand as well

8.  Update config/org.tfvars with the oauth client id and secret

9.  Create an artifact registry for docker images. You may need to enable the Artifact Registry API. Consider enabling vulnerability scanning

## Build

1.  Initialise terraform with

    ```
    make init
    ```

2.  Apply the terraform with

    ```
    make apply
    ```

3.  Configure DNS with the load balancer IPs returned by terraform, viz
    atos-digital.net -> prod load balancer IP
    dev.atos-digital.net -> dev load balancer IP

4.  Make a note of the project service account email. Add this user to the seed project wih rights to read and write containers in seed project.

## Post-build

Create the following
secrets.DEV_DB_PASSWORD
secrets.DEV_GCP_WORKLOAD_IDP
secrets.DEV_GCP_SERV_ACCOUNT

secrets.PROD_DB_PASSWORD
secrets.PROD_GCP_WORKLOAD_IDP
secrets.PROD_GCP_SERV_ACCOUNT

secrets.OAUTH2_CLIENT_SECRET

## Remote access

Use the jumpbox to access the private network and componets attached to the network e.g. database

The best way to connect is using the gcloud cli. Find the jumpbox in the console and cut-and-paste the connection command. It will be something like

```
gcloud compute ssh --zone "europe-west2-c" "jumpbox" --project "<your project name>"
```

You should install the postgres tools with

```
sudo apt-get update && sudo apt-get install postgresql
```

Connect to the database with

```
psql -h INSTANCE_IP -U USERNAME DATABASE
```

for example

```
psql -h 10.x.x.x -U dbuser defaultdb
```
