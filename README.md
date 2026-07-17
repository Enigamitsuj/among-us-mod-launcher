# Among Us Mod Launcher

A modern desktop launcher for Among Us community mods.

**First supported mod:** [Town of Us: Mira](https://github.com/AU-Avengers/TOU-Mira)

Built with Go + [Wails v3](https://v3.wails.io/), React, TypeScript, Tailwind CSS, and Framer Motion.

Distributed as a **single Windows executable** — no Node, .NET, or Electron runtime required for end users.

---

## Philosophy

This is not a bare installer. It should feel like an official indie-studio launcher:

- Launch → Install → Play
- Zero configuration for most users
- Original Among Us install is **never** modified
- Mods install next to the launcher executable by default

Installer created by **FBI OpenUp**.

---

## Features (v0.1)

- Dark, frameless launcher UI with custom title bar
- Mod registry (Town of Us: Mira first; more mods later)
- Automatic Steam / Among Us detection
- GitHub Releases version picker (latest / stable / beta labels)
- Download → extract → install progress
- Optional desktop shortcut
- Optional launch after install
- Reinstall confirmation dialog

---

## Development

### Prerequisites

- Go 1.25+
- Node.js 20+
- Wails CLI v3 (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- WebView2 (included with modern Windows)

### Setup

```bash
cd frontend
npm install
cd ..
wails3 generate bindings -ts -i ./...
```

### Run (dev)

```bash
wails3 dev
```

### Build

```bash
wails3 build
```

The executable lands under `bin/`.

---

## Project layout

```
internal/
  appservice/   # Wails-bound API for the UI
  mods/         # Mod registry + shared types
  steam/        # Steam + Among Us detection
  githubapi/    # GitHub Releases client
  installer/    # Download / extract / install
  shortcut/     # Desktop shortcut creation
  fsutil/       # Path helpers
frontend/       # React + Tailwind UI
```

---

## Contributing

Public contributions are welcome via pull requests.

- Open a PR against `main`
- **PRs require review/approval before merge**
- Keep packages small and logic independent of the UI

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

Community project. Mod assets and trademarks belong to their respective owners.
Among Us is © Innersloth. Town of Us: Mira is maintained by AU-Avengers.
