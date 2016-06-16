mackerel-plugin-darwin-battery
=====================

Darwin battery custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-darwin-battery [-metric-key-prefix=battery-capacity]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.battery]
command = "/path/to/mackerel-plugin-darwin-battery"
```
