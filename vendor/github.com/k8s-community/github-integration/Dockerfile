FROM scratch

ENV GITHUBINT_LOCAL_PORT 8080
ENV GITHUBINT_BRANCH "release/"

ENV GITHUBINT_TOKEN "Webhook secret is in integration settings on Github"
ENV GITHUBINT_PRIV_KEY "Private key is in integration settings on Github"
ENV GITHUBINT_INTEGRATION_ID "Integration ID is in it's settings on Github"

# DB parameters
ENV GITHUBDB_USER githubint
ENV GITHUBDB_PASSWORD githubint
ENV GITHUBDB_NAME github_integration

ENV CICD_BASE_URL http://k8s-build-01:8080
ENV USERMAN_BASE_URL https://services.k8s.community/user-manager

COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/github-integration /

CMD ["/github-integration"]

EXPOSE $GITHUBINT_LOCAL_PORT
