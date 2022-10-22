# Timesheet Calculator

> Sends emails about my contract hours

This function calculates the remaining hours in my contract as a reminder for
what I need to be doing daily and monthly in order to maximise my income.

## Getting Started

This runs in OpenFaaS or faasd.

Get the complete list of possible options by calling `make`.

### Deploying

To deploy it, run either `faas-cli up -f pollen-detector.yml` or `make prod/up`.

### Developing

To develop locally, run `make dev/reload` which will create a container with all
the code from this folder and hot reload it on change.

## Secrets

**deployment**

Secrets must be manually created in the kubernetes cluster. To do that in this repo
run the following command.

```shell
cd secrets
for f in *; do                                                       
cat $f | faas-cli secret create $f 
done
```

**development**

The following files are needed in the `secrets` repo.

```shell
./secrets/b2appkey
./secrets/b2bucket
./secrets/b2keyid
./secrets/b2server
./secrets/mail-from
./secrets/mail-password
./secrets/mail-username
```
Docker compose will use this to load in the variables during development.

## Requirements

This function requires `cron-connector`. To install it into the cluster use 
`arkade install cron-connector`. 

Email SMTP access. I use fastmail. Note: fastmail port recommended 465 but my
VPS provider blocks 465, so we use 587 instead.
