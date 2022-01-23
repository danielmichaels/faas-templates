# Blog Example

This function will look up a user supplied IP address with [ipinfo.io]. From the IP address it 
will then retrieve the weather for that IP addresses city using [openweathermap].

## Usage

### Secrets

1. `weather-key` is the [openweathermap] API key. This must be set or the function will fail to build.
2. `ipinfo-key` is the [ipinfo.io] API key. It must be set.

**Setting Secrets**

Secrets must be set in the `faasd` instance. To create or update existing secrets you must 
either pass in a value from a file, or through stdin.

Using a file, assuming the API key is in a file at the root of this directory.

`faas-cli secret create ipinfo-key --from-file=ipinfo-key.txt`

See the [OpenFaaS docs] for more details.

### Environment Variables

`ORIGINS` can be set with any valid URL. It does have defaults set in the `stack.yml` file. This 
is needed for CORS when called from a domain other than the `faasd` instance.

A note about environment variables. When referencing an envar inside a `faasd` function, refer 
to the `environment` yaml file value, not the environment variable. For example,

```yaml
# Define the environment variable in your shell as usual.
ORIGINS=http://localhost:8000 && faas-cli up`
# Use the environment key of `origins` from stack.yml in your code not ORIGINS.
environment:
      origins: "${ORIGINS}"
# in the code, how I've referenced it
origins := os.Getenv("origins")
```
### Commands

To build:

`faas-cli build` or `make prod/build`

Push to docker

`faas-cli push` or `make prod/push`

Deploy to `faasd`

`faas-cli deploy` or `make prod/deploy`

Do all these steps

`faas up` or `make prod/up`

What `make` commands are available?

Run `make` or `make help` to see the options.

[ipinfo.io]: https://ipinfo.io
[openweathermap]: https://openweathermap.org
[openfaas docs]: https://docs.openfaas.com/reference/secrets/
