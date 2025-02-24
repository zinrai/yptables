# yptables

yptables is a tool that generates iptables configurations from YAML files. It supports both shell script and iptables-restore formats.

## Features

- YAML-based configuration for iptables rules
- Supports filter and nat tables
- Handles both built-in and custom chains
- Generates either shell script or iptables-restore format
- Full support for iptables match modules

## Installation

```bash
$ go build -o yptables cmd/main.go
```

## Usage

```bash
yptables [options] <config.yaml>

Options:
  -format string
        Output format: 'script' or 'restore' (default "script")
  -output string
        Output file (default: stdout)
```

## Configuration Example

See [example.yaml](example.yaml) for a configuration example.

## Examples

Generate shell script:

```bash
$ yptables config.yaml > iptables.sh
```

Generate iptables-restore format:

```bash
$ yptables -format restore config.yaml > iptables.rules
```

Generated configurations can be validated using iptables-restore --test:

```bash
$ sudo iptables-restore --test iptables.rules
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
