# foodcli

Go CLI: login to foodora + show active order status.

Status: early prototype. API details reverse-engineered from an Android XAPK (`25.44.0`).

## Build

```sh
go test ./...
go build ./cmd/foodcli
```

## Configure country / base URL

Bundled presets (from the APK):

```sh
./foodcli countries
./foodcli config set --country HU
./foodcli config set --country AT
./foodcli config show
```

Manual:

```sh
./foodcli config set --base-url https://hu.fd-api.com/api/v5/ --global-entity-id NP_HU --target-iso HU
```

## Login

`oauth2/token` needs a `client_secret` (the app fetches it via remote config). `foodcli` auto-fetches it on first use and caches it locally.

Optional override (keeps secrets out of shell history):

```sh
export FOODORA_CLIENT_SECRET='...'
./foodcli login --email you@example.com --password-stdin
```

If MFA triggers and you're running in a TTY, `foodcli` prompts for the OTP code and retries automatically. Otherwise it stores the MFA token locally and prints a safe retry command (`--otp <CODE>`).

### Client headers

Some regions (e.g. Austria/mjam `mj.fd-api.com`) expect app-style headers like `X-FP-API-KEY` / `App-Name` / app `User-Agent`. `foodcli` uses an app-like header profile for `AT` by default.

For corporate flows, you can override the OAuth `client_id`:

```sh
./foodcli login --email you@example.com --client-id corp_android --password-stdin
```

### Cloudflare / bot protection

Some regions (e.g. Austria/mjam `mj.fd-api.com`) may return Cloudflare HTML (`HTTP 403`) for plain Go HTTP clients.

Use an interactive Playwright session (you solve the challenge in the opened browser window; no auto-bypass):

```sh
./foodcli login --email you@example.com --password-stdin --browser
```

Prereqs: `node` + `npx` available. First run may download Playwright + Chromium.

Tip: use a persistent profile to keep browser cookies/storage between runs (reduces re-challenges):

```sh
./foodcli login --email you@example.com --password-stdin --browser --browser-profile "$HOME/Library/Application Support/foodcli/browser-profile"
```

### Import cookies from Chrome (no browser run)

If you already solved bot protection / logged in in Chrome, you can import the cookies for the current `base_url` host:

```sh
./foodcli cookies chrome --profile "Default"
./foodcli orders
```

If the bot cookies live on the website domain (e.g. `https://www.foodora.at/`), import from there and store them for the API host:

```sh
./foodcli cookies chrome --url https://www.foodora.at/ --profile "Default"
```

If you have multiple profiles, try `--profile "Profile 1"` (or pass a profile path / Cookies DB via `--cookie-path`).

### Import session from Chrome (no password)

If you’re logged in on the website in Chrome, you can import `refresh_token` + `device_token` and then refresh to an API access token:

```sh
./foodcli session chrome --url https://www.foodora.at/ --profile "Default"
./foodcli session refresh --client-id android
./foodcli history
```
If `session refresh` errors with “refresh token … not found”, that site session isn’t valid for your configured `base_url` (common for some regions).

## Orders

```sh
./foodcli orders
./foodcli orders --watch
./foodcli history
./foodcli history --limit 50
./foodcli history show <orderCode>
./foodcli order <orderCode>
```

## Safety

This talks to private APIs. Use at your own risk; rate limits / bot protection may block requests.
