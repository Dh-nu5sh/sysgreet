# Hostinfo

> Beautiful, low-latency system context for every terminal login.

Hostinfo keeps you oriented the moment a shell prompt appears. It prints the
hostname in ASCII art alongside a curated snapshot of operating system,
network, and resource telemetry—so you always know **which** machine you are on
and **whether** it is healthy. Built for managing home labs and fleets alike, it
remains lightweight, offline-friendly, and cross-platform across Linux, macOS,
and Windows.

---

## Why Hostinfo exists

I created Hostinfo while operating a growing home lab and juggling
multiple SSH sessions. I wanted a professional banner—_not_ a novelty—that
instantly answered three questions:

1. **Where am I logged in?** (Hostname, OS, architecture, remote source)
2. **Is this host behaving?** (Uptime, memory, disk, CPU trends)
3. **What network path am I on?** (Primary route, relevant secondary interfaces)

Hostinfo delivers those answers in under 50 ms without calling out to the
network or depending on external runtimes.

---

## Highlights

- **Single static binary** – Go 1.22+, no CGO, no daemons, no service
  dependencies.
- **Cross-platform parity** – Linux/macOS show load averages; Windows surfaces
  CPU usage. Interface filtering avoids noisy virtual adapters everywhere.
- **Configurable yet optional** – YAML or TOML profiles toggle sections, pick
  fonts/colors, set layout order, and cap the interface list. Defaults “just
  work” with zero files.
- **Graceful degradation** – Missing metrics or SSH metadata simply fall back;
  the banner keeps rendering.
- **Performance-guarded** – Startup benchmark (<50 ms median, <80 ms p95) runs in
  CI; process RSS stays <15 MB.
- **Professional aesthetics** – Embedded FIGlet fonts, ANSI color with automatic
  monochrome fallback, 80-column mindful layout.

---

## Quick start

### Install the binary

```bash
# Via Go (requires Go 1.22+)
go install github.com/veteranbv/hostinfo/cmd/hostinfo@latest

# Or download a release artifact (Linux/macOS/Windows, amd64 & arm64)
# https://github.com/veteranbv/hostinfo/releases
```

> _Tip:_ The binary runs entirely offline. Copy it between hosts without
> worrying about external assets.

### Wire into your shell

| Shell            | Snippet                                                                                       |
|------------------|------------------------------------------------------------------------------------------------|
| Bash / Zsh       | `echo 'hostinfo' >> ~/.bashrc` (or `~/.zshrc`)                                                 |
| Fish             | `echo 'hostinfo' >> ~/.config/fish/config.fish`                                               |
| PowerShell       | `Add-Content $PROFILE 'hostinfo'`                                                             |
| Windows Terminal | Add `hostinfo` to your profile script so it runs after each session attaches                  |
| SSH `ForceCommand` | `ForceCommand /usr/local/bin/hostinfo && /bin/bash` (keeps banner even when no profile runs) |

Need to silence the banner temporarily? Use `hostinfo --disable` in CI jobs or
scripts that call the shell non-interactively.

---

## Configuration (optional)

Hostinfo looks for configuration in this order:

1. `HOSTINFO_CONFIG` environment variable (absolute or `~/` paths)
2. `~/.config/hostinfo/config.yaml` (or `.yml`, `.toml`)
3. `~/.hostinfo.yaml` / `.toml`

Example YAML:

```yaml
# ~/.config/hostinfo/config.yaml
ascii:
  font: "slant"
  color: "cyan"
  monochrome: false

display:
  hostname: true
  os: true
  ip_addresses: true
  remote_ip: true
  uptime: true
  user: true
  memory: true
  disk: true
  load: true
  datetime: true
  last_login: true

layout:
  compact: false
  sections: ["header", "network", "system", "resources"]

network:
  show_interface_names: true
  max_interfaces: 4
```

Environment variables override everything (e.g.
`HOSTINFO_DISPLAY_MEMORY=false`, `HOSTINFO_ASCII_FONT=standard`). See
[`configs/example.yaml`](configs/example.yaml) and
[`configs/example.toml`](configs/example.toml) for full references.

---

## What the banner shows

```ascii
 _   _           _   _        __ _
| | | | ___  ___| |_(_) ___  / _(_) __ _ _ __ ___
| |_| |/ _ \/ __| __| |/ __|| |_| |/ _` | '_ ` _ \
|  _  |  __/\__ \ |_| | (__ |  _| | (_| | | | | | |
|_| |_|\___||___/\__|_|\___||_| |_|\__,_|_| |_| |_|

Linux 6.8.0 (x86_64)

System
  Uptime: 4d 12h 33m
  User: alice /home/alice
  Time: Fri, 10 Oct 2025 09:45:00 PDT

Network
  Primary: 192.168.1.42 (en0)
  Secondary: 10.8.0.2 (utun2)
  Remote: 203.0.113.5

Resources
  Mem: 12.3GB free / 16.0GB (23% used)
  Disk: 210.0GB used / 512.0GB (41% used)
  CPU Load: 0.45 0.52 0.60
```

- **System** – Hostname (ASCII art), OS name/version, architecture, uptime,
  active user + home, current time, last login when available.
- **Network** – Primary outbound interface based on routing table, filtered list
  of secondary physical interfaces, SSH remote IP (from `SSH_CONNECTION` or
  `SSH_CLIENT`). Loopback, link-local, Docker/VM, and down interfaces stay out of
  view by default.
- **Resources** – Memory, disk, and CPU metrics with highlight thresholds (≥75% in
  yellow, ≥90% in red). Windows surfaces realtime CPU usage; Unix hosts show load
  averages.

---

## Performance guarantees

- **Startup** – `< 50 ms` median, `< 80 ms` p95 (validated by
  `go test -bench Startup ./test/benchmarks`)
- **Binary footprint** – `< 10 MB` for all release targets (GoReleaser checks)
- **Runtime memory** – `< 15 MB` RSS for default banner
- **No network activity** – All data collected locally, offline-safe

Enable `HOSTINFO_DEBUG=1` to log collector errors without interrupting output.

---

## Development & contribution

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development guidelines, code standards, and workflow.

**Quick start:**

```bash
git clone https://github.com/veteranbv/hostinfo.git
cd hostinfo
go mod tidy
make test
make bench
```

**Common tasks:**

```bash
make fmt            # Format code
make lint           # Run linters (requires golangci-lint)
make test-coverage  # Run tests with coverage report
make build          # Build the binary
```

PRs are welcome—please open an issue describing new collectors, layout ideas, or
platform-specific improvements before diving in.

---

## Release process

- CI (`.github/workflows/ci.yml`) runs `golangci-lint`, unit tests with race
  detection, integration tests, and validates startup performance (<80ms p95).
- Releases use GoReleaser (`.goreleaser.yml`) to ship signed binaries for
  Linux/macOS (amd64/arm64) and Windows (amd64), plus checksums.
- `go install github.com/veteranbv/hostinfo@VERSION` is validated during the
  release workflow.

---

## Roadmap

- Optional JSON output for scripting in CI/CD pipelines
- Extended GPU/storage telemetry for workstation profiles
- Pluggable section framework (e.g., Kubernetes context, vault status)
- Prebuilt Windows installer for enterprise onboarding

Ideas welcome—open a discussion if a feature would make Hostinfo more useful for
your fleet.

---

## License

Hostinfo is licensed under the [Apache License 2.0](LICENSE).

Copyright © 2025 Henry Sowell
