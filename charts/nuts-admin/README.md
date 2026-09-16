# Helm Chart for nuts-admin

This chart deploys [nuts-admin](https://github.com/nuts-foundation/nuts-admin), the management UI for a Nuts node, on a Kubernetes cluster.

It creates a Deployment, Service, ConfigMap (rendering `config.yaml`) and, optionally, an Ingress. It does not run or manage a Nuts node itself; point it at an existing one through `env.nutsNodeAddress`.

## Configuration

All configurable properties can be found in [./values.yaml](./values.yaml). The most relevant ones:

| Property             | Description                                                                                                            | Default                    |
|----------------------|--------------------------------------------------------------------------------------------------------------------------|-----------------------------|
| `env.nutsNodeAddress` | Internal API address of the Nuts node this instance manages, injected as `NUTS_NODE_ADDRESS`.                          | `http://localhost:8081`    |
| `config`              | Rendered into `config.yaml`, mounted into the pod and referenced via `NUTS_CONFIGFILE`. See the [nuts-admin README](../../README.md) and [deploy/config.yaml.example](../../deploy/config.yaml.example) for the options it accepts, e.g. `oidc` and `credentialprofiles`. | `{}`                        |
| `oidc.existingSecret` | Name of an existing Secret holding the OIDC client secret, injected as `NUTS_OIDC_CLIENT_SECRET`. Use this instead of setting `config.oidc.client.secret` directly, which would put it in the ConfigMap in plaintext. The rest of the OIDC config (`enabled`, `metadata`, `client.id`, `scope`) still goes through `config.oidc`. | `""` |
| `oidc.existingSecretKey` | Key within `oidc.existingSecret` holding the client secret.                                                         | `client-secret`             |
| `ingress.enabled`     | Expose nuts-admin through an Ingress.                                                                                  | `false`                    |

## Installing nuts-admin

### From source

Execute the following command from the root of the chart folder. Replace `<NAME>` with the name you wish to give this Helm installation.

```shell
helm install <NAME> .
```

### From the Nuts Helm repo

Add the repo:

```shell
helm repo add nuts-admin https://nuts-foundation.github.io/nuts-admin/
helm repo update
```

Then install it, optionally overriding values with your own `values.yaml`:

```shell
helm install -f values.yaml <NAME> nuts-admin/nuts-admin-chart
```
