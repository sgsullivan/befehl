# befehl
## Run arbitrary commands over ssh in mass

- run the given payload(s) in PWD on host(s) in config.json.. up to 2000 at a time.

`./befehl execute --runconfig config.json --routines 2000`

Output of each payload run for every node will be in the log directory (by default, its `$HOME/befehl/logs`) in a file named after the machine it ran on.

An example runconfig (`config.json` shown above) is shown below:

```json
{
  "payload": "integration_tests/examples/payload",
  "user": "root",
  "hosts": [{
      "host": "127.0.0.1",
      "port": 1000
    },
    {
      "host": "127.0.0.1",
      "port": 1001,
      "user": "snowflake",
      "payload": "integration_tests/examples/payload-override"
    },
    {
      "host": "127.0.0.1",
      "port": 1002
    },
    {
      "host": "127.0.0.1",
      "port": 1003
    },
    {
      "host": "127.0.0.1",
      "port": 1004
    }
  ]
}
```

As you can see, you can override the `payload` from the default, as `127.0.0.1:1001` is doing in the example above.
## Configuration

You can configure befehl with a config file (~/.befehl.[toml|json|yaml]) any serialization format that upstream viper supports befehl supports for the config file. Valid configuration options:

```toml
[general]
logdir = "/home/ssullivan/log-special"
[ssh]
privatekeyfile = "/home/ssullivan/alt/id_rsa"
knownhostspath = "/home/asullivan/alt/.ssh/known_hosts"
hostkeyverificationenabled = true
```

Unless enabled as shown above, ssh known host verification is disabled.

## Obtaining prebuilt binaries

Head on over to the [releases page](https://github.com/sgsullivan/befehl/releases) to get prebuilt binaries for your platform.

## Building

Once you have your Go environment setup, it should be as simple as cloning this git repo and running `make`. The resulting binary will be located at `_exe/befehl`.
