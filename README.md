# patchwork-deploy

Minimal deployment orchestrator that applies ordered shell patches to remote hosts over SSH.

---

## Installation

```bash
go install github.com/yourname/patchwork-deploy@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/patchwork-deploy.git
cd patchwork-deploy && go build -o patchwork-deploy .
```

---

## Usage

Define your patches as numbered shell scripts in a directory:

```
patches/
  001_install_deps.sh
  002_migrate_db.sh
  003_restart_services.sh
```

Run against a remote host:

```bash
patchwork-deploy --host user@192.168.1.10 --patches ./patches
```

patchwork-deploy will connect over SSH and execute each script in order, stopping on the first failure.

### Options

| Flag | Description |
|------|-------------|
| `--host` | Remote host in `user@host` format |
| `--patches` | Path to directory containing patch scripts |
| `--key` | Path to SSH private key (default: `~/.ssh/id_rsa`) |
| `--dry-run` | Print patches that would be applied without executing |

---

## How It Works

1. Scans the patches directory and sorts scripts lexicographically.
2. Opens an SSH session to the target host.
3. Executes each script sequentially, streaming output to stdout.
4. Halts immediately if any script exits with a non-zero status.

---

## License

MIT © yourname