# oauth-proxy

OAuth Proxy Service.

You can use this service to create users and necessary environment.
Right now it supports only GitHub.

## Getting Started

You can store page templates in `/templates` and static files (images, css etc.) in `/static` directories.

For example, you can checkout our [Landing Page](https://github.com/k8s-community/k8s-community).
In this case please follow this instruction:
 
1. Get static and template files: `go get -u github.com/k8s-community/k8s-community`.

2. Get this service: `go get -u github.com/k8s-community/oauth-proxy`.

3. Make symbolic links to `static` and `template` directories:

        ln -s $GOPATH/src/github.com/k8s-community/k8s-community/static ./static
        ln -s $GOPATH/src/github.com/k8s-community/k8s-community/templates ./templates
    
Otherwise, you can just create empty `static` and `template directories` and define necessary files yourself.

## How to run the service

To run the service you need to define these environment variables:

| Variable | Description | Example |
|---|---|---|
| SERVICE_HOST | Host listen by the service | 0.0.0.0 |
| SERVICE_PORT | Port listen by the service| 80 |
| GITHUB_CLIENT_ID | [ClientID](https://github.com/settings/developers) of your application | f778... |
| GITHUB_CLIENT_SECRET | [ClientSecret](https://github.com/settings/developers) of your application  | 807ff71... |

For example, you can run service using `go run` (not for production, only for experiment!):


    env SERVICE_HOST=0.0.0.0 SERVICE_PORT=80 GITHUB_CLIENT_ID=f778... GITHUB_CLIENT_SECRET=807ff71... go run oauth-proxy.go


**TODO:** Add link to chart with configuration description.
