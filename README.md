# nuts-admin
Application which shows how to integrate with the Nuts node to administer identities.

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

The following properties configure OIDC user authorization in Nuts admin:
- `oidc.enabled` or `NUTS_OIDC_ENABLED`: set to `true` to enable OIDC user authentication.
- `oidc.metadata` or `NUTS_OIDC_METADATA`: points to the OIDC metadata endpoint, e.g. `https://auth.example.com/.well-known/openid-configuration`.
- `oidc.client.id` or `NUTS_OIDC_CLIENT_ID`: the client ID to use for OIDC authentication.
- `oidc.client.secret` or `NUTS_OIDC_CLIENT_SECRET`: the client secret to use for OIDC authentication.
- `oidc.scope` or `NUTS_OIDC_SCOPE`: the scope(s) to use for OIDC authentication, defaults to `openid`, `profile`, and `email`.

The following properties should be used if API authentication is enabled on the Nuts node:
- `node.auth.keyfile` or `NUTS_NODE_AUTH_KEYFILE`: points to a PEM encoded private key file. The corresponding public key should be configured on the Nuts node in SSH authorized keys format.
- `node.auth.user` or `NUTS_NODE_AUTH_USER`: must match the user in the SSH authorized keys file.
- `node.auth.audience` or `NUTS_NODE_AUTH_AUDIENCE`: must match the configured audience.

## User Authentication

This application does support OIDC user authentication. This has only been tested with Azure Entra ID, but it should work with any OIDC provider.

### Preliminary warning

However, if OIDC user authentication is not enabled, make sure to restrict access in any other case than local development.
The application proxies REST API calls to the configured Nuts node, so leaving it unsecured will allow anyone to access the proxied Nuts node REST APIs.

### Configuration on Azure
You can have users logged in with their Azure Entra ID account. Nuts Admin will authenticate to Azure Entra ID using `client_id` and `client_secret`.
The configuration (as environment variables) could look as below.
It assumes Nuts Admin runs on `https://example.com/admin/`. Make sure to fill in the correct `tenantId` and `client_id` and `client_secret`.

```shell
NUTS_URL=https://example.com/admin/
NUTS_OIDC_ENABLED=true
NUTS_OIDC_METADATA=https://login.microsoftonline.com/<tenantId>/v2.0/.well-known/openid-configuration
NUTS_OIDC_CLIENT_ID=<client_id>
NUTS_OIDC_CLIENT_SECRET=<client_secret>
```

Typically, you create an enterprise application in Azure Entra ID and configure the redirect URI to `https://example.com/admin/auth/openid-connect/callback` (given the example base URL).
Then, you make sure the users or groups that you want to give access are assigned to the enterprise application.
To let Nuts Admin access Azure Entra ID, you can register a user-assigned managed identity in Azure and generate the `client_id` and `client_secret`.

For more information regarding OIDC on Azure, see https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols-oidc.

## Requesting Credentials

This application supports requesting credentials using the Nuts node through OpenID4VCI.

You need to configure the credential issuers using the following properties:

* `NUTS_CREDENTIALPROFILES_<ID>_TYPE`: the type of the credential to be issued, e.g. `DegreeCredential`.
* `NUTS_CREDENTIALPROFILES_<ID>_ISSUER`: the URL of the issuer, e.g. `https://issuer.example.com`.

When there is at least one profile configured, the option will become avaiable in the subject/wallet view.

You'll also need to enable the authorization endpoint on the Nuts node for OpenID4VCI to work using `NUTS_AUTH_AUTHORIZATIONENDPOINT_ENABLED`.

## Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
make dev
```

You can access the website at `http://localhost:1305/` by default

The API and domain types are generated from the `api/api.yaml`.
```shell
make gen-api
```

### Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com
