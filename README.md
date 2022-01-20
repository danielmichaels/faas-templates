# OpenFaaS Templates

My collection of OpenFaaS templates.

> Functions as a Service

These templates should work on OpenFaaS whether running in a Kubernetes cluster, or on [faasd] which
is what these run on primarily.

You can test them out using docker, or by running [faasd] using [mutlipass](https://multipass.run).
For more instructions on how to get started refer to
the [faasd multipass instructions](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md)

[faasd]: https://github.com/openfaas/faasd

## Expectations

1. `faas-cli` is installed on your host
2. A `faasd` or OpenFaaS instance available for deployment

### Helpful tips

I find it easiest to set the following environment variables so that I do not have to explicitly 
declare them on every invocation of `faas-cli`. 

- `OPENFAAS_PREFIX=docker.io/danielmichaels` or use another registry. You'll need to `docker login` as well
- `OPENFAAS_URL=http://$FAAS_DOMAIN`; the domain of your instance. Defaults to localhost:8080

## OpenFaaS?

These templates are for use in the [OpenFaaS] project, created by [Alex Ellis](https://alexellis.io)
. Learn more about the project [here][openfaas]

[openfaas]: https://openfaas.com
