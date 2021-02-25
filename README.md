
# keycloak

[![ci](https://github.com/zemirco/keycloak/workflows/ci/badge.svg)](https://github.com/zemirco/keycloak/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/zemirco/keycloak.svg)](https://pkg.go.dev/github.com/zemirco/keycloak)

keycloak is a Go client library for accessing the [Keycloak API](https://www.keycloak.org/docs-api/12.0/rest-api/index.html).

## Installation

```bash
go get github.com/zemirco/keycloak
```

## Usage

```go
import "github.com/zemirco/keycloak"

// create your oauth configuration
config := oauth2.Config{
    ClientID: "admin-cli",
    Endpoint: oauth2.Endpoint{
        TokenURL: "http://localhost:8080/auth/realms/master/protocol/openid-connect/token",
    },
}

// get a valid token from keycloak
ctx := context.Background()
token, err := config.PasswordCredentialsToken(ctx, "admin", "admin")
if err != nil {
    panic(err)
}

// create a new http client that uses the token on every request
client := config.Client(ctx, token)

// create a new keycloak instance and provide the http client
k := keycloak.NewKeycloak(client)

// start using the library and, for example, create a new realm
realm := &keycloak.Realm{
    Enabled: keycloak.Bool(true),
    ID:      keycloak.String("myrealm"),
    Realm:   keycloak.String("myrealm"),
}

res, err := k.Realms.Create(ctx, realm)
if err != nil {
    panic(err)
}
```

## Development

Use `docker-compose` to start Keycloak locally.

```bash
docker-compose up -d
```

Keycloak is running at http://localhost:8080/. The admin credentials are `admin` (username) and `admin` (password). If you want to change them simply edit the `docker-compose.yml`.

Keycloak uses PostgreSQL and all data is kept across restarts.

Use `down` if you want to stop the Keycloak server.

```bash
docker-compose down
```

## Architecture

The main entry point is `keycloak.go`. This is where the Keycloak instance is created. It all starts in this file.

Within Keycloak we also have the concept of clients. They are the ones that connect to Keycloak for authentication and authorization purposes, e.g. our frontend and backend apps. That is why this library simply uses the `keycloak` instance of type `Keycloak` and not a `client` instance like [go-github](https://github.com/google/go-github). Although technically this library is a Keycloak client library for Go. However this distinction should make it clear what is meant when we talk about a client in our context.

## Testing

You need to have Keycloak running on your local machine to execute the tests. Simply use `docker-compose` to start it.

All tests are independent from each other. Before each test we create a realm and after each test we delete it. You don't have to worry about it since the helper function `createRealm` does that automatically for you. Inside this realm you can do whatever you want. You don't have to clean up after yourself since everything is deleted automatically when the realm is deleted.

Run all tests.

```bash
go test -race -v ./...
```

Create code coverage.

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

We have also provided a simple `Makefile` that run both jobs automatically.

```bash
make
```

Open `coverage.html` with your browser.

## Design goals

1. Zero dependencies

    It's just the Go standard library.

1. Idiomatic Go

    Modelled after [go-github](https://github.com/google/go-github) and [go-jira](https://github.com/andygrunwald/go-jira).

1. Keep authentication outside this library

    This is the major difference to most of the other Go Keycloak libraries.

    We leverage the brilliant [oauth2](https://github.com/golang/oauth2) package to deal with authentication. We have provided multiple examples to show you the workflow. It basically means we do not provide any methods to call the `/token` endpoint.

1. Return struct and HTTP response

    Whenever the Keycloak API returns JSON content you'll get a proper struct as well as the HTTP response.

    ```go
    func (s *ClientsService) Get(ctx context.Context, realm, id string) (*Client, *http.Response, error)
    ```

## Related work

- https://github.com/Nerzal/gocloak
- https://github.com/PhilippHeuer/go-keycloak
- https://github.com/coreos/go-oidc
- https://github.com/keycloak/kcinit
- https://github.com/pulumi/pulumi-keycloak/tree/master/sdk/go/keycloak
- https://github.com/airmap/go-keycloak
- https://github.com/cloudtrust/keycloak-client
- https://github.com/myENA/go-keycloak
- https://github.com/threez/go-keycloak

## License

MIT
