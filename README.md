# Among Us Mod Launcher

A modern Windows desktop launcher for Among Us community mods.

**First supported mod:** [Town of Us: Mira](https://github.com/AU-Avengers/TOU-Mira)

Built with Go + [Wails v3](https://v3.wails.io/), React, TypeScript, Tailwind CSS, and Framer Motion.

Distributed as a **single Windows executable** — no Node, .NET, or Electron runtime required for end users.

---

## Download & install

1. Open the [latest GitHub Release](https://github.com/Enigamitsuj/among-us-mod-launcher/releases/latest).
2. Download `among-us-mod-launcher-<version>-windows-x64.exe` (and the matching `.sha256` file if you want to verify).
3. Optional integrity check in PowerShell:

```powershell
Get-FileHash .\among-us-mod-launcher-1.0.0-windows-x64.exe -Algorithm SHA256
Get-Content .\among-us-mod-launcher-1.0.0-windows-x64.exe.sha256
```

4. Run the executable.

### SmartScreen warning

Release builds are currently **unsigned**. The first launch may show **Windows protected your PC**.

Choose **More info** → **Run anyway**. Prefer downloads from this GitHub repository only, and verify the SHA256 when possible.

### Requirements

- Windows 10/11 x64
- [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/) (usually already installed on modern Windows)
- A local Among Us install (Steam, Epic, Itch, or Xbox / Microsoft Store)

---

## Philosophy

This should feel like an indie-studio launcher, not a bare zip installer:

- Launch → Install → Play
- Zero configuration for most users
- Original Among Us install is **never** modified
- Mods install **beside** the detected game in a separate folder (for example `Among Us - TOU Mira`)

Maintained by **Enigamitsuj**.

---

## Features (v1.0)

- Dark, frameless launcher UI with custom title bar
- Mod registry starting with Town of Us: Mira
- Automatic Among Us detection across Steam, Epic, Itch, and Xbox
- Version / compatibility checks against selected releases
- GitHub Releases version picker (latest / stable / beta labels)
- Download → extract → install progress with atomic updates
- Install / Play / Uninstall based on the selected vs installed version
- Launcher-owned install markers (uninstall only removes launcher-managed folders)

---

## Development

### Prerequisites

- Go 1.25+
- Node.js 22+
- Wails CLI v3 pinned to `v3.0.0-alpha2.117` (match CI)
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

### Test

```bash
go test ./...
```

### Build

```bash
wails3 task build
```

The executable lands under `bin/`.

Optional local metadata stamp:

```bash
node frontend/scripts/stamp-windows.mjs bin/among-us-mod-launcher.exe 1.0.0 build/windows/icon.ico
```

---

## Project layout

```
internal/
  appservice/     # Wails-bound API for the UI
  mods/           # Mod registry + shared types
  game/           # Among Us detection + compatibility
  githubapi/      # GitHub Releases client
  installer/      # Download / extract / atomic install
  installrecord/  # Ownership markers for managed installs
  fsutil/         # Path helpers
frontend/         # React + Tailwind UI
.github/workflows # CI + tagged releases
```

---

## Releasing

Maintainers publish by pushing a SemVer tag such as `v1.0.0`.

The release workflow:

1. Stamps Windows version metadata from the tag
2. Builds the Windows executable
3. Publishes `among-us-mod-launcher-<version>-windows-x64.exe` plus a `.sha256` sidecar

The git tag is the source of truth for release versioning. Local defaults in `frontend/package.json`, `build/config.yml`, and `build/windows/info.json` should stay aligned between releases.

Code signing is not wired yet; SmartScreen warnings are expected until an Authenticode certificate is added.

---

## Contributing

Public contributions are welcome via pull requests.

- Open a PR against `main`
- **PRs require review/approval before merge**
- Keep packages small and logic independent of the UI

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Security

See [SECURITY.md](SECURITY.md) for how to report vulnerabilities.

---

## Disclaimer & trademarks

This project is an unofficial community tool and is **not affiliated with, endorsed by, or associated with Innersloth LLC**.

Among Us and related marks are © Innersloth LLC. This launcher never redistributes the base game. You must own a legitimate copy of Among Us to use it.

Town of Us: Mira is maintained by [AU-Avengers](https://github.com/AU-Avengers/TOU-Mira) and is licensed separately under GPL-3.0. Mod assets and trademarks belong to their respective owners.

---

## License

This project's source code is licensed under the [GNU General Public License v3.0](LICENSE).
