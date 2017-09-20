FROM alpine:3.5

ENV USERMAN_SERVICE_PORT 8080

ENV K8S_HOST "https://master.k8s.community"
ENV K8S_TOKEN "Token is for access to k8s API"
ENV TLS_SECRET_NAME "tls-secret"
ENV DOCKER_REGISTRY_SECRET_NAME "registry-pull-secret"

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY user-manager /

EXPOSE $USERMAN_SERVICE_PORT

CMD ["/user-manager"]