FROM busybox

# Service env parameters
ENV SERVICE_HOST 0.0.0.0
ENV SERVICE_PORT 8080

# Services-Dependencies
ENV USERMAN_BASE_URL https://services.k8s.community/user-manager

# DB parameters
ENV COCKROACHDB_PUBLIC_SERVICE_HOST localhost
ENV COCKROACHDB_PUBLIC_SERVICE_PORT 26257
ENV COCKROACHDB_USER k8scomm
ENV COCKROACHDB_PASSWORD k8scomm
ENV COCKROACHDB_NAME k8s_community

# additional env parameters
ENV GITHUB_CLIENT_ID f778...
ENV GITHUB_CLIENT_SECRET 807ff71...
ENV GITHUB_OAUTH_STATE just-a-very-secret-state
ENV K8S_GUEST_TOKEN Gfn5Kf0e1Fisg4b9Fmv6FdS8b5dSo6JC

COPY certs /etc/ssl/
COPY static /static
COPY templates /templates
COPY ui /

EXPOSE $SERVICE_PORT

CMD ["/ui"]
