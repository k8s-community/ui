# handlers

This library represents a set of useful middleware and handlers for [takama's router](https://github.com/takama/router).

## List of existed handlers (middleware)

### Info

Info handler shows useful information about the service.

        r := router.New()

        version := "1.0.0"
        repo := "handlers"
        commit := "019cc819f8af4e2f7533fb3760f21387a4ef0cce"

        r.GET("/info", info.Handler(version, repo, commit))
        r.Listen(":3000")

### Bearer Token

Bearer token middleware implements [the OAuth 2.0 Authorization Framework: Bearer Token Usage](https://tools.ietf.org/html/rfc6750).

        r := router.New()

        token := "s-fdF8-mF_9.B-4.1Cfd"
        h := func(ctx *router.Control) {}

        r.GET("/test", bearer.WithToken(token, h))
        r.Listen(":3000")

## Contributing

Contributors are welcome! Please, follow the [Contributing Guidelines](CONTRIBUTING.md).

If you have any questions, feel free to [create an issue](https://github.com/openprovider/handlers/issues/new).

Contributors (unsorted):

- [Elena Grahovac](https://github.com/rumyantseva)

## Current version

0.1.1
