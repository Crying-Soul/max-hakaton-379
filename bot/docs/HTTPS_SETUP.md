# HTTPS setup for local development

This document explains how to enable HTTPS for the server and how to generate self-signed certificates for local testing.

1) Environment

- Edit the repository `.env` file (at project root) and set:

```
TLS_ENABLED=true
TLS_CERT_PATH=certs/server.crt
TLS_KEY_PATH=certs/server.key
```

2) Generate self-signed certs (Linux / macOS)

Run the provided script from the project root:

```bash
./certs/generate_certs.sh
```

This will create `certs/server.crt` and `certs/server.key`.

3) Trusting the certificate (optional)

- For browsers to accept the cert without warnings you can add the `server.crt` to your OS/browser trusted certificates. The exact steps depend on your OS.

4) Security notes

- Don't commit private keys to git. The `certs/.gitignore` file ignores `.key` and `.pem` files by default.
- For production use, obtain certificates from a trusted CA and never enable TLS with self-signed certs in public environments.

5) Integration notes

- The `.env` file contains `SERVER_HOST` and `SERVER_PORT`. When `TLS_ENABLED=true` start the server in HTTPS mode using the paths from `TLS_CERT_PATH` and `TLS_KEY_PATH`.
- If your server binary reads env vars directly, ensure it uses `TLS_ENABLED` and loads certs accordingly.
