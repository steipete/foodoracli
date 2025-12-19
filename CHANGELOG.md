# Changelog

## 0.1.0 (Unreleased)

- Initial CLI (`login`, `orders`, `order`, `config`, `countries`)
- Past orders (`history` via `orders/order_history`)
- Auto-fetch/cache OAuth `client_secret` from Firebase Remote Config
- OAuth token flow with refresh + MFA detection (`mfa_triggered`)
- Interactive OTP prompt + retry (TTY)
- Order tracking endpoints (`tracking/active-orders`, `tracking/orders/{orderCode}`)
- Optional Playwright interactive login (`--browser`) + Cloudflare cookie capture (e.g. Austria/mjam)
- `--config` flag works (use separate config files for testing)
- OAuth `--client-id` override (e.g. `corp_android`)
