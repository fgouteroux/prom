# prom

prom is a prometheus tooling cli. 
It contains theses 2 projects in one and it could be extended more easily.

- https://github.com/fgouteroux/prom-push-cli
- https://github.com/fgouteroux/promtoolfmt

## Main Usage

```
NAME:
   prom - Prometheus tooling

USAGE:
   prom [global options] command [command options] [arguments...]
   
AUTHOR:
   - François Gouteroux <francois.gouteroux@gmail.com>
   
COMMANDS:
   metrics  metrics operations
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -D    show debug output (default: false) [$PROM_DEBUG, $DEBUG]
   --help, -h     show help
   --version, -v  print the version
   
COPYRIGHT:
   (c) 2023 François Gouteroux
```

## Metrics usage

```
NAME:
   prom metrics - metrics operations

USAGE:
   prom metrics command [command options] [arguments...]

COMMANDS:
   push     push prometheus metrics to a remote write.
   check    check prometheus metrics.
   fmt      format prometheus metrics.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

