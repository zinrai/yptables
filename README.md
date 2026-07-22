# yptables

yptables is a tool that generates iptables-restore format configurations from YAML files.

## Features

- YAML-based configuration for iptables rules
- Supports filter and nat tables
- Handles both built-in and custom chains
- Generates iptables-restore format for atomic rule application
- Full support for iptables match modules

## Usage

```bash
$ yptables <config.yaml>
```

## Configuration Example

See [examples](./examples) directory for a configuration example.

## Examples

Generate iptables-restore format:

```bash
$ yptables config.yaml > iptables.rules
```

Validate generated configuration:

```bash
$ sudo iptables-restore --test iptables.rules
```

Apply configuration:

```bash
$ sudo iptables-restore iptables.rules
```

## License

This project is licensed under the [MIT License](./LICENSE).
