# befehl
## Run arbitrary commands over ssh in mass

- run payload.sh in PWD on hosts in targets file in PWD.. up to 2000 at a time.

`./befehl execute --hosts targets --payload payload.sh -routines 2000`

Output of each payload run will be in the log directory (by default, its `$HOME/befehl/logs`) in a file named after the machine it ran on.
The targets file should be a plain text file containing all hosts to run the payload on, separated by a new line.

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
