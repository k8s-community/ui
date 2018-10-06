FROM scratch

ENV USERMAN_LOCAL_PORT 8080

ENV K8S_BASE_URL "https://master.k8s.community"
ENV K8S_TOKEN "Token is for access to k8s API"
ENV TLS_SECRET_NAME "tls-secret"
ENV DOCKER_REGISTRY_SECRET_NAME "registry-pull-secret"

COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/user-manager /

EXPOSE $USERMAN_LOCAL_PORT

CMD ["/user-manager"]