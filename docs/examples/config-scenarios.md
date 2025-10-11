# Configuration Scenarios

## Minimal Configuration (disable resources)

```yaml
display:
  memory: false
  disk: false
  load: false

ascii:
  font: "standard"
  color: "cyan"
```

## Focus on Networking Only

```yaml
display:
  hostname: true
  os: true
  ip_addresses: true
  remote_ip: true
  uptime: false
  user: false
  memory: false
  disk: false
  load: false
  datetime: true
  last_login: false

network:
  show_interface_names: true
  max_interfaces: 5
```

## Troubleshooting Metric Discrepancies

If resource values appear inconsistent with native tools:

1. Run `go test -run TestResourceCollectorMatchesSystemStats ./test/integration/...` on the target OS.
2. Verify the user running `sysgreet` has permission to read filesystem metadata for the home directory.
3. Ensure no virtualization layers hide physical interfaces (VPN, container bridges). Adjust `network.max_interfaces` or disable specific sections if needed.
4. For Windows hosts, the CPU usage calculation relies on `cpu.PercentWithContext`. When the banner runs during login scripts, the first sampling window may be noisy; run the banner twice or lower the sampling interval via configuration if necessary.

## Non-interactive Bootstrap Policies

- `CI=1 bin/sysgreet --config-policy=keep` ensures automation never writes a config file (ideal for hosts that manage configs externally).
- `CI=1 SYSGREET_CONFIG_POLICY=overwrite bin/sysgreet` regenerates the default config on every run without prompting and keeps the latest backup beside the active file.
- If both a flag and environment variable are provided, the flag wins so one-off jobs can override fleet defaults.
