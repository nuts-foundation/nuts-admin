# nuts-admin
Application which shows how to integrate with the Nuts node to administer identities.

## Warning on authentication

This application does not support authentication. Make sure to restrict access in any other case than local development. 
The application proxies REST API calls to the configured Nuts node, so leaving it unsecured will allow anyone to access the proxied Nuts node REST APIss.

## Building and running

### Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
make dev
```

The API and domain types are generated from the `api/api.yaml`.
```shell
make gen-api
```

### Docker
```shell
$ docker run -p 1305:1305 nutsfoundation/nuts-admin
```

## Configuration
When running in Docker without a config file mounted at `/app/config.yaml` it will use the default configuration.

The `node.auth.keyfile` config parameter should point to a PEM encoded private key file. The corresponding public key should be configured on the Nuts node in SSH authorized keys format.
`node.auth.user` Is required when using Nuts node API token security. It must match the user in the SSH authorized keys file.

## Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com
