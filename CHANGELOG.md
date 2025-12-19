# Changelog

## 0.1.0 (Unreleased)

- Initial CLI (`login`, `orders`, `order`, `config`, `countries`)
- Fetch/store OAuth `client_secret` from Firebase Remote Config (`secret fetch`)
- OAuth token flow with refresh + MFA detection (`mfa_triggered`)
- Order tracking endpoints (`tracking/active-orders`, `tracking/orders/{orderCode}`)
