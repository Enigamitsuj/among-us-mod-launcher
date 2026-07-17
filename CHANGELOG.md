# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Unified publisher branding to Enigamitsuj across app metadata and docs
- Documented SmartScreen guidance and SHA256 verification for unsigned Windows releases
- CI now runs `go test ./...`
- Detect protected install folders early and surface a clear Program Files / UAC message
- Make icon rebuild script portable; keep committed icon assets as the build source of truth

## [1.0.0] - TBD

### Added

- Windows single-exe launcher built with Go + Wails v3
- Town of Us: Mira as the first supported mod
- Among Us detection for Steam, Epic, Itch, and Xbox
- Release version picker from GitHub Releases
- Install / update / uninstall flow with launcher ownership markers
- Atomic install promotion so failed updates keep the previous mod copy
- Dark frameless UI with install progress and running-game status

### Security / notes

- Release builds are unsigned; Windows SmartScreen may warn on first run
- Original Among Us install is never modified; mods install beside the game
