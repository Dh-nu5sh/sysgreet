# Default Sysgreet Banner Snapshot

This reference output documents the constitutional gates enforced by the Sysgreet CLI.

## Constitution Gates Checklist

- ✅ Single-binary Go CLI (no external runtime dependencies)
- ✅ Cross-platform parity (Linux, macOS, Windows collectors tested)
- ✅ Startup < 50ms measured via `go test -bench Startup ./test/benchmarks`
- ✅ Colorful ASCII output with monochrome fallback
- ✅ Offline operation with embedded assets and configuration defaults

## Sample Output (mock data)

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

Use this snapshot for QA validation and regression testing until golden files are finalized.
