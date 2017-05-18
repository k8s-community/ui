FROM alpine:3.5

# Service env parameters
ENV SERVICE_HOST 0.0.0.0
ENV SERVICE_PORT 8080

# additional env parameters
ENV GITHUB_CLIENT_ID f778...
ENV GITHUB_CLIENT_SECRET 807ff71...

CMD apk --no-cache add ca-certificates && update-ca-certificates

COPY vendor/github.com/k8s-community/k8s-community/static /static
COPY vendor/github.com/k8s-community/k8s-community/templates /templates
COPY oauth-proxy /

EXPOSE $SERVICE_PORT

CMD ["/oauth-proxy"]
