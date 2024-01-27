# Secret Sharing

This web application allows you to create and view your secrets for a specific file.

## Running the App

1. `go build -o secret-app cmd/secret-share-web/main.go`

2. `DATA_FILE_PATH=./data.json ./secret-app`

Where `DATA_FILE_PATH` is a path for a file that will store your secrets

if needed to check the project's version, you can use flag `-v`
```bash
./secret-app -v
```

## API
1. `curl http://localhost:8080/healthcheck` - checks the server status
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

It will respond to a decrypted secret
```json
{"data":"My super secret"}
```
