# pam-k8s-sa

PAM module for authenticating using Kubernetes Service Account tokens.

This was created to authenticate with PostgreSQL, and the default parameters reflect that.

## Example usage with PostgreSQL

When run inside a pod in the cluster, all default parameters should work out of the box in a standard cluster.

The following example assumes the PostgreSQL server is running outside the cluster and Kubernetes API server can be access at `https://127.0.0.1:6443`.

In `/etc/pam.d/psql-k8s-sa`:

```
auth required /path/to/pam_k8s_sa.so \
    server_url=https://127.0.0.1:6443 \
    ca_file=/etc/postgresql/kubernetes_ca.crt
account required /path/to/pam_k8s_sa.so
```

> Here we're assuming we have manually written the Kubernetes cluster CA to `/etc/postgresql/kubernetes_ca.crt`. As we're not providing any `token_file` parameter we need to have the token available at `/var/run/secrets/kubernetes.io/serviceaccount/token` or allow anonymous access to the OpenID Connect Discovery endpoint.

In `pg_hba.conf` add the following line (adjust the actual CIDR):

```
host all +k8s_sa 0.0.0.0/0 pam pamservice=psql-k8s-sa
```

> This will instruct PostgreSQL to use PAM authentication using PAM service `psql-k8s-sa` for every role that has been granted the `k8s_sa` role.

In PostgreSQL you can then run the following:

```sql
CREATE ROLE k8s_sa;
CREATE ROLE app_sa$app_ns WITH LOGIN;
GRANT k8s_sa TO app_sa$app_ns;
```

After that you can connect to PostgreSQL as user `app_sa$app_ns` using the service account token for service account `app-sa` in namespace `app-ns` as password.

## Parameters

#### server_url

Default: `https://kubernetes.default.svc.cluster.local`.

The URL to Kubernetes API server. This will be used to discover OpenID Connect configuration using `{{server_url}}/.well-known/openid-configuration`.

#### issuer

Default: `https://kubernetes.default.svc.cluster.local`.

The expected value of the `iss` field in the Service Account's JWT token.

#### audience

Default: `https://kubernetes.default.svc.cluster.local`.

The expected value in the `aud` field in the Service Account's JWT token.

#### username_template

Default: `{{.Name | replace "-" "_"}}${{.Namespace | replace "-" "_"}}`.

A template that is rendered before comparing its output value to the username that is provided to the PAM module. The template is a good default for using PAM authentication in PostgreSQL where dashes are not allowed and `$` is the most natural separator between name and namespace.

`{{.Name}}` is the name of the service account taken from the JWT token and `{{.Namespace}}` is the Kubernetes namespace of the service account, also taken from the JWT token.

There is a single function available called `replace`, which is a direct call to [strings.ReplaceAll](https://pkg.go.dev/strings#ReplaceAll).

#### token_file

Default: `/var/run/secrets/kubernetes.io/serviceaccount/token`.

A path to a file with an authentication token in order to be able to use the OpenID Connect discovery URL of the Kubernetes API Server. If the file does not exist it will skip authentication, in which case [anonymous access has to be enabled](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#service-account-issuer-discovery).

#### ca_file

Default: `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`.

A path to CA certificates to use for TLS verification in requests to the API Server OpenID Connect discovery.

#### verify_tls

Default: `true`.

Can be used to disable TLS verification in requests to the API Server OpenID Connect discovery URL.
