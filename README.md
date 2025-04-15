# Prometheus Exporter Custom

This project contains a custom Prometheus exporter.

## 🚀 Getting Started

You can develop in this project using **Nix** or **Dev Containers** (e.g. with GitHub Codespaces or VS Code Remote Containers).

### 🧪 Using `nix develop`

If you have Nix with flakes enabled, run:

```bash
nix develop
```

This drops you into a development shell with `go`, `go-task`, and any other tools preconfigured. The prompt will show `[nix develop]` to let you know you're in the environment.

### 🐳 Using Dev Containers

If you're using VS Code and have the "Dev Containers" extension:

1. Open the command palette (`Cmd+Shift+P` / `Ctrl+Shift+P`)
2. Select **"Dev Containers: Reopen in Container"**
3. VS Code will automatically build and enter the dev environment

This will set up a container with all necessary tools (e.g., Go and Task) as defined in `.devcontainer`.

## 📦 Dependencies

- [Go](https://golang.org/)
- [Task](https://taskfile.dev/)

These are automatically provided via `nix develop` or Dev Containers.

## 📜 License

MIT License
