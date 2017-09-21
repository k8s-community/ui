# k8s-community

Install our application on GitHub and we'll deliver your services to Kubernetes.

You can use this service to create users and necessary environment.
Right now it supports only GitHub.

## How to run the service

To run the service you need to define these environment variables:

| Variable | Description | Example |
|---|---|---|
| SERVICE_HOST | Host listen by the service | 0.0.0.0 |
| SERVICE_PORT | Port listen by the service| 80 |
| GITHUB_CLIENT_ID | [ClientID](https://github.com/settings/developers) of your application | f778... |
| GITHUB_CLIENT_SECRET | [ClientSecret](https://github.com/settings/developers) of your application  | 807ff71... |

For example, you can run service using `make run` (not for production, only for experiment!):


    env SERVICE_HOST=0.0.0.0 SERVICE_PORT=80 \
    GITHUB_CLIENT_ID=f778... GITHUB_CLIENT_SECRET=807ff71... \
    COCKROACHDB_PUBLIC_SERVICE_HOST=localhost COCKROACHDB_PUBLIC_SERVICE_PORT=26257 \
    COCKROACHDB_USER=k8scomm COCKROACHDB_PASSWORD=k8scomm COCKROACHDB_NAME=k8s_community \
    K8S_GUEST_TOKEN=12345 \
    USERMAN_BASE_URL=https://services.k8s.community/user-manager \
    go run k8s-community.go


**TODO:** Add link to chart with configuration description.


## Useful commands

Run migration:

    kubectl run -it --rm cockroach-client --image=cockroachdb/cockroach --restart=Never --command -- ./cockroach sql --host cockroachdb-public --insecure --database=k8s_community < 000.sql
