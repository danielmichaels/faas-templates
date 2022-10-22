# OpenFaaS Templates

My collection of OpenFaaS templates.

> Functions as a Service

These templates should work on OpenFaaS whether running in a Kubernetes cluster, or on [faasd] which
is what these run on primarily.

You can test them out using docker, or by running [faasd] using [multipass](https://multipass.run).
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

I like to change the `<name>.yml` of each function to `stack.yml` so that I can call `faas-cli` 
commands without `-f <name>.yml` each time.

Setting `alias faas=faas-cli` is also handy which is in my `.zshrc` file, with completions. To 
setup completion, you can add the following (to zshrc):

`command -v faas-cli >/dev/null 2>&1 && source <(faas-cli completion --shell zsh)`

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

Note: `new-faas` is a [bash script] that can automate this for you (and) more.

[bash script]: https://raw.githubusercontent.com/danielmichaels/dot/master/scripts/new-faas

`export PROJECT=pollen-detector; mkdir -p $PROJECT; cd $_; faas-cli template store pull golang-middleware && faas-cli new --lang golang-middleware $PROJECT`

To see what other templates are possible run `faas-cli template store list` and replace 
`golang-middleware` with the template you want.

## Usage

To build, push to a registry, deploy to your OpenFaaS instance, you **must** first `cd` into the 
function you're explicitly wanting to use.

For instance, I want to build, push and deploy `banner-grab` to my instance, this is the process.

1. `cd banner-grab`
2. `faas-cli build`
3. `faas-cli push`
4. `faas-cli deploy`
5. **tip** do 2,3,4 in on power move; `faas-cli up`

## Hot reloading in development

Through the use of docker-compose, its now possible to do hot-reloading whilst in development. 
So far, this has proven effective for `golang-middleware` but `python` and `node` are also supported
via this [link](https://simonemms.com/blog/2020/08/12/live-reload-for-openfaas).

The only thing that does not work "*out of the box*" is `static` files. This is because the 
`docker-compose.yml` alters the build, making the paths different between production and dev.

To get around this, you can correct the pathing with a `if`, or helper function such as:

```go
// This requires DOCKER_COMPOSE=true to be set as an environment variable.
// Set this environment variable in the Makefile and docker-compose.environment.
func isDockerCompose() string {
	file := os.Getenv("DOCKER_COMPOSE")
	index := "./function/static/index.html"
	if file == "" {
		index = "./static/index.html"
	}
	return index
}

```

Otherwise, this method works quite well; `make dev/reload` to get started.

## Adding Secrets

Most templates have a Makefile which handles creating and using secrets in development. However,
in production, the secret must be added to kubernetes.

`faas-cli secret create <secret-name> --from-file=<secret.txt>`

Check the secret is in the correct namespace.

`kubectl get secrets -n openfaas-fn`

### Error, no golang-middleware?

The `golang-middleware` template does not live in the `openfaas` standard templates. This means 
you need to pull this each time.

`faas-cli template store pull golang-middleware`
## OpenFaaS?

These templates are for use in the [OpenFaaS] project, created by [Alex Ellis](https://alexellis.io)
. Learn more about the project [here][openfaas]

[openfaas]: https://openfaas.com
