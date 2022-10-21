# Timesheet Calculator

> what do?

## Getting Started

### Deploying

### Developing

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
