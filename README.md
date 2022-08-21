# MSDS432-Phase3-Chicago Business Intelligence
This repository contains source code to complete Phase-3 requirements for the Chicago Business Intelligence for Strategic Planning project.
This Readme provides the step by step guide to install and deploy source code to Google Cloud. 

## Files Overview
1. Folder 'src': It has following 3 folders under pkg Folder
    - store: This folder has all the Go code for pulling data from Chicago portal and push to postgres data lake 
    - types: This folder stores the metadata and types of data tables 
    - api: http handler 
2. Cloudbuild.yaml: This file has the steps for google cloud build process. In essence, it deploys 2 microservices:
    - go-microservice: Builds a docker image using Dockerfile for go-microservice, push it to google container registry and deploy it to cloud run
    - pdadmin: Pulls a docker image for pgadmin, push it to google container registry and deploy it to cloud run
3. Dockerfile: This is the docker file required for go-microservice
4. go.mod/sum: To fulfill Go packages dependencies
5. main.go: Go main file firing source code under 'src' Folder

## Installation & Deployment Guide
1. Postgres database setup
    - Create database instance of postgres using the following command.
      > gcloud sql instances create cbipostgres --database-version=POSTGRES_14 --cpu=2 --memory=7680MB --region=us-central
    - Create sql users on the database instance using the following command.
      > gcloud sql users set-password postgres --instance=cbipostgres --password=root
    - Create a database for our microservice using the following command.
      > gcloud sql databases create chicago_business_intelligence --instance=mypostgres
2. Setting up continuous deployment using Cloud Build
    - Under 'Cloud Build' under cloud console, create deployment trigger with source as this repository.
    - put trigger event as "push to main branch"
    - Ensure all the relevant IAM permissions/roles are correct. This is very important!!
3. Get the Postgres DB instance connection name
    - Go to the "cbipostgres" instance created in the previous steps and copy the instance connection name. Ex: cbi-yunus:us-central1:cbipostgres
    - Go to line 24 in the main.go source code file and update the connection string with copied instance connection name as shown below.
      >user=postgres dbname=chicago_business_intelligence password=root host=/cloudsql/cbi-yunus:us-central1:cbipostgres sslmode=disable port = 5432
4. Update cloudbuild.yaml file with correct Project-ID. Ex: cbi-yunus
5. Push these changes to github main branch to trigger cloud build immediately and wait for the build to complete.
6. Verify "go-microservice" and "pgaadmin" services are up and running
7. Verify data being updated in postgres data lake
    - From Cloud Run, click on pgadmin, copy the highlighted URL
    - Open the URL in a Browser and Login to pgadmin console
    - Add server with appropriate server name, hostname and login/password as declared before for postgres.
    - After login, click on Chicago_business intelligence
    - Click on schemas/tables and verify you see at least one of the CBI tables with data

## Congratulations! You finished the tutorial!!
