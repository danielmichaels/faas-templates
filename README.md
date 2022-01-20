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

## Layout

This repo has the following layout

```shell
.
├── banner-grab
│   └── handler.go
├── build
│   └── banner-grab
├── stack.yml
└── template
    ├── golang-http
    └── golang-middleware
```

Each new template should live inside its own directory. This keeps the root directory clear of 
several `function-name.yml` files. 

It does add a lot of `template` directories, as every function directory will need its own. 
These are ignored but will need to exist locally.

## How do I create a new function?

To create a new golang-middleware function using the cli:

`mkdir -p name && cd $_ && faas-cli template store pull golang-middleware && faas-cli new --lang golang-middleware name`

To see what other templates are possible run `faas-cli template store list` and replace 
`golang-middleware` with the template you want.


### Error, no golang-middleware?

The `golang-middleware` template does not live in the `openfaas` standard templates. This means 
you need to pull this each time.

`faas-cli template store pull golang-middleware`
## OpenFaaS?

These templates are for use in the [OpenFaaS] project, created by [Alex Ellis](https://alexellis.io)
. Learn more about the project [here][openfaas]

[openfaas]: https://openfaas.com
