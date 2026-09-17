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

Note: nuts-admin keeps OIDC sessions in memory, so `replicaCount` above `1` together with `config.oidc.enabled: true` isn't supported — the chart refuses to render (`fail`) in that combination. There's no autoscaling support in this chart for the same reason.

## Installing nuts-admin

### From source

Execute the following command from the root of the chart folder. Replace `<NAME>` with the name you wish to give this Helm installation.

```shell
helm install <NAME> .
```

### From GitHub Container Registry

Chart releases are published as an OCI artifact to `ghcr.io` on every change to `charts/` on `main`. Install directly by version, optionally overriding values with your own `values.yaml`:

```shell
helm install -f values.yaml <NAME> oci://ghcr.io/nuts-foundation/helm-nuts-admin --version <VERSION>
```
