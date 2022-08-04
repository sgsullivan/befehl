# befehl
## Run arbitrary commands over ssh in mass

- run payload.sh in PWD on hosts in targets file in PWD.. up to 2000 at a time.

`./befehl execute --hosts targets --payload payload.sh -routines 2000`

Output of each payload run for every node will be in the log directory (by default, its `$HOME/befehl/logs`) in a file named after the machine it ran on.

The targets file should be a plain text file (shown below) containing all hosts to run the payload on, separated by a new line. If the host has an alternate ssh port (aka not port 22) then specify the alternate port like `192.168.0.2:2222`. An example host list is shown below specifying alternate ssh port:

```
192.168.0.2
192.168.0.3:1000
192.168.0.4:22
```

In this example, the connection attempt to 192.16.0.2 will be attempted on port 22 because the port wasn't specified.

## Configuration

You can configure befehl with a config file (~/.befehl.[toml|json|yaml]) any serialization format that upstream viper supports befehl supports for the config file. Valid configuration options:

```toml
[general]
logdir = "/home/ssullivan/log-special"
[auth]
privatekeyfile = "/home/ssullivan/alt/.ssh/id_rsa"
sshuser = "nonrootuser"
```

These options should be self explanatory so I wont describe what each does here.

## Obtaining prebuilt binaries

Head on over to the [releases page](https://github.com/sgsullivan/befehl/releases) to get prebuilt binaries for your platform.

## Building

Once you have your Go environment setup, it should be as simple as cloning this git repo and running `make`. The resulting binary will be located at `_exe/befehl`.
