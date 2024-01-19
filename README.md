# Secret Sharing

This web application afford yu to create and view secrets

## Running the App

1. `go build -o secret-app`

2. `DATA_FILE_PATH=./data.json ./secret-app`

Where `DATA_FILE_PATH` is a path for a file, which will store your secrets

## API
1. GET `/healthcheck` - checks the server status
2. Save your secret
```bash
curl -X POST http://localhost:8080 -d '{"plain_text":"My super secret"}'
```
It will create a file in `DATA_FILE_PATH`, store information there, and return the response containing your secret encrypted.

```json
{"id":"c0331ab6a4fad09a50e441644d2d676c"}
```

3. Read the secret

```bash
curl http://localhost:8080/c0331ab6a4fad09a50e441644d2d676c
```

It will respond a decrypted secret
```json
{"data":"My super secret"}
```

## A little dev description

The `filestore` package is responsible for reading and writing data to/from a file.
The file `handlers/secret_handler` writes a new secret to a file when the http method is `POST` and reads from the file,
when the http method is `GET`.
