# foodoracli

Go CLI: login to foodora + show active order status.

Status: early prototype. API details reverse-engineered from an Android XAPK (`25.44.0`).

## Build

```sh
go test ./...
go build ./cmd/foodoracli
```

## Configure country / base URL

Bundled presets (from the APK):

```sh
./foodoracli countries
./foodoracli config set --country HU
./foodoracli config show
```

Manual:

```sh
./foodoracli config set --base-url https://hu.fd-api.com/api/v5/ --global-entity-id NP_HU --target-iso HU
```

## Login

`oauth2/token` needs a `client_secret` (the app fetches it via remote config). Provide it as env var:

```sh
export FOODORA_CLIENT_SECRET='...'
./foodoracli login --email you@example.com --password-stdin
```

Or fetch it via Firebase Remote Config (APK-derived config) and store in your local config:

```sh
./foodoracli secret fetch --store
```

If MFA triggers, rerun with the printed `--mfa-token` and pass `--otp <CODE>`.

## Orders

```sh
./foodoracli orders
./foodoracli orders --watch
./foodoracli order <orderCode>
```

## Safety

This talks to private APIs. Use at your own risk; rate limits / bot protection may block requests.
