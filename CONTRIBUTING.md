# Contributing

Thanks for helping improve Among Us Mod Launcher.

## Workflow

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Open a pull request against `main`

`main` is protected: **pull requests must be reviewed and approved before merge**.

## Guidelines

- Prefer small, focused PRs
- Keep backend packages independent of the React UI
- Friendly user-facing errors only (no stack traces in the UI)
- Never modify the original Among Us install directory
- Default install location stays next to the launcher executable

## Local development

```bash
cd frontend && npm install && cd ..
wails3 generate bindings -ts -i ./...
wails3 dev
```

## Code of conduct

Be respectful. This is a community tool for players and modders.
