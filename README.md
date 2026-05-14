# configvault

A minimal secrets manager that syncs environment configs from Vault or AWS SSM into local dotenv files.

---

## Installation

```bash
go install github.com/yourusername/configvault@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/configvault.git
cd configvault
go build -o configvault .
```

---

## Usage

Pull secrets from **HashiCorp Vault** into a local `.env` file:

```bash
configvault pull --source vault --path secret/myapp --out .env
```

Pull secrets from **AWS SSM Parameter Store**:

```bash
configvault pull --source ssm --path /myapp/production --out .env
```

Push a local `.env` file back to your secrets backend:

```bash
configvault push --source vault --path secret/myapp --in .env
```

### Example `.env` output

```env
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
API_KEY=supersecretkey
DEBUG=false
```

### Configuration

configvault reads backend credentials from environment variables or a `~/.configvault.yaml` config file. See [docs/configuration.md](docs/configuration.md) for full details.

---

## Supported Backends

| Backend         | Status |
|-----------------|--------|
| HashiCorp Vault | ✅ Stable |
| AWS SSM         | ✅ Stable |
| AWS Secrets Manager | 🚧 Coming soon |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)