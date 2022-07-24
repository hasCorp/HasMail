hasMail
=======

Centralized email automation service for hasCorp.

## Run locally

A complete example:
```
go run cmd/mailservice/main.go -bypass -port=2564
```

### Authentication

#### SendGrid
Store your SendGrid API key in a `credentials.json` at the root of the repo 
(this is in the `.gitignore`). Take a look at `./credentials.json.template` 
as an example:
```json
{
    "SENDGRID_API_KEY": "Your API key goes here",
    "FROM_NAME": "Friendly Name",
    "FROM_ADDR": "some-email@fake.dev"
}
```

If you don't want to store the credentials and configurations in a JSON file in
the project repo, you can expose environment variables with the same key names.
```bash
# for Linux/Unix systems
export SENDGRID_API_KEY="Your API key goes here",
export FROM_NAME="Friendly Name",
export FROM_ADDR="some-email@fake.dev"
```
```powershell
# for Windows systems
set SENDGRID_API_KEY="Your API key goes here",
set FROM_NAME="Friendly Name",
set FROM_ADDR="some-email@fake.dev"
```

#### Bypassing client authentication
For local development, it makes sense to do some testing without requiring
a hard dependency on the authentication service to verify incoming requests.
When running locally, pass in the `-bypass` flag to ignore client auth verification
```
go run cmd/mailservice/main.go -bypass
```

#### Listening port
By default, the HTTP server listens on port `8000`. This can be changed when
running via the `-port` flag:
```
go run cmd/mailservice/main.go -bypass -port=2564
```

### Building
You can build the project locally simply by running:
```bash
go build .
```

Or you can use the `Dockerfile` at the root of the repo to build an image.
```bash
docker build -t hascorp/hasmail -f Dockerfile .
```

### Running
You can run the project locally with `go`:
```bash
go run ./cmd/mailservice
```

Or you can use the built Docker image to run a container:
```bash
docker run -it -p 8000:8000 hascorp/hasmail
```

### Testing
Ping the server with cURL or your preferred client:
```bash
# ping healthcheck endpoint
curl localhost:8000/

# verify routes work with no-op endpoint
curl -d '{"a": "b"}' -H 'Content-Type: application/json' localhost:8000/mail/noop

# send a sample mail
curl -d '{"name": "Hank Pecker", "vars": {"foo": "bar"}, "recipient": "hank@hascorp.dev"}' -H 'Content-Type: application/json' localhost:8000/mail/sample
```

## Releasing

### Pipeline
TODO: this

### Dockerize
Build with the production Dockerfile:
```bash
docker build -t hascorp/hasmail-prod -f Dockerfile.production .
```

This can be tested locally like the regular Dockerfile:
```bash
docker run -it -p 8000:8000 hascorp/hasmail-prod
```

## Testing

### Unit testing
Run unit tests locally:
```bash
go test -v ./...
```

### Integration testing
TBD
