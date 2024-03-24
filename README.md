# nuts-admin
Application which shows how to integrate with the Nuts node to administer identities.

## Warning on authentication

This application does not support user authentication. Make sure to restrict access in any other case than local development. 
The application proxies REST API calls to the configured Nuts node, so leaving it unsecured will allow anyone to access the proxied Nuts node REST APIss.

## Running
Example running the application, connecting to a Nuts node running on `http://nutsnode:8081`:

```shell
$ docker run -p 1305:1305 -e NUTS_NODE_ADDRESS=http://nutsnode:8081 nutsfoundation/nuts-admin:latest
```

When running in Docker without a config file mounted at `/app/config.yaml` it will use the default configuration.

The application can be configured through `/app/config.yaml` or environment variables.
It supports the following configuration options:

- `port` or `PORT`: overrides the default HTTP port (`1305`) the application listens on. 
- `node.address` or `NUTS_NODE_ADDRESS`: points to the internal API of the Nuts node, e.g. `http://nutsnode:8081`.

The following properties should be used if API authentication is enabled on the Nuts node:
- `node.auth.keyfile` or `NUTS_NODE_AUTH_KEYFILE`: points to a PEM encoded private key file. The corresponding public key should be configured on the Nuts node in SSH authorized keys format.
- `node.auth.user` or `NUTS_NODE_AUTH_USER`: must match the user in the SSH authorized keys file.
- `node.auth.audience` or `NUTS_NODE_AUTH_AUDIENCE` must match the configured audience.

## Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
make dev
```

The API and domain types are generated from the `api/api.yaml`.
```shell
make gen-api
```

### Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com
