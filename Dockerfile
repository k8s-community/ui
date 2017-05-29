FROM alpine:3.5

# Service env parameters
ENV SERVICE_HOST 0.0.0.0
ENV SERVICE_PORT 8080

# Services-Dependencies
ENV USERMAN_BASE_URL https://services.k8s.community/user-manager

# DB parameters
ENV COCKROACHDB_SERVICE_HOST localhost
ENV COCKROACHDB_SERVICE_PORT 26257
ENV COCKROACHDB_USER k8scomm
ENV COCKROACHDB_PASSWORD k8scomm
ENV COCKROACHDB_NAME k8s_community

# additional env parameters
ENV GITHUB_CLIENT_ID f778...
ENV GITHUB_CLIENT_SECRET 807ff71...
ENV K8S_GUEST_TOKEN Gfn5Kf0e1Fisg4b9Fmv6FdS8b5dSo6JC

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY static /static
COPY templates /templates
COPY k8s-community /

EXPOSE $SERVICE_PORT

CMD ["/k8s-community"]
