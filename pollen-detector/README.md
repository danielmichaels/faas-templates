# Pollen Detector

> Get daily updates for Canberra's pollen count

Canberra is known for its terrible pollen counts during spring and into early
summer. As a hay fever sufferer, I sometimes forget to take my morning anti-histamine.
This function sends a daily reminder if the pollen count is above `Low`.

## Getting Started

This runs in OpenFaaS or faasd. 

To develop locally, run `make dev/reload` which will create a container with all
the code from this folder and hot reload it on change.

To deploy it, run either `faas-cli up -f pollen-detector.yml` or `make prod/up`.

Get the complete list of possible options by calling `make`.

## Secrets

In development, secrets are stored locally on disk. This functions requires; 
`chat-id` and `tg-token` for sending Telegram messages - the function will error
without these. Running `make dev/reload` will make the files accessible for local 
development, but they must exist and be populated within this folder.

In production, the secrets must be added to the `openfaas-fn` namespace. This can be
done with the `faas-cli`.

```shell
# create a secret
faas-cli secret create <name>

# update a secret
faas-cli secret update <name>
```

Secrets must be referenced within the code using a helper function like this:

```go
// getSecret retrieves the secret from openfaas and makes it available for use.
func getSecret(secretName string) ([]byte, error) {
	secret, err := ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", secretName))
	if err != nil {
		return nil, err
	}
	return secret, nil
}
```

## Requirements

This function requires `cron-connector`. To install it into the cluster use 
`arkade install cron-connector`. 
