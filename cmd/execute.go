package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/sgsullivan/befehl"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute the given payload(s) against the given host(s) from configuration",
	Long: `Execute the given payload(s) against the given host(s) from configuration specified
by the runconfig flag. Below is an example runconfig:

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
	  }
	]
  }

You can control how many payloads are executed concurrently by passing the routines flag.

By default befehl will use the private key in $HOME/.ssh/id_rsa. This can be overrode by
specifying auth.privatekeyfile in ~/.befehl.[toml|json|yaml].

By default befehl will write the output of each payload for each host in $HOME/befehl/logs. This
can be overrode by specifying general.logdir in ~/.befehl.[toml|json|yaml].

Heres an example specifying all supported options:

[general]
logdir = "/home/ssullivan/log-special"
[ssh]
privatekeyfile = "/home/ssullivan/alt/id_rsa"
knownhostspath = "/home/asullivan/alt/.ssh/known_hosts"
hostkeyverificationenabled = true

`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfig, _ := cmd.Flags().GetString("runconfig")
		routines, _ := cmd.Flags().GetInt("routines")

		if routines == 0 {
			color.Yellow("--routines not given, defaulting to 30..\n")
			routines = 30
		}

		instance, err := befehl.New(&befehl.Options{
			PrivateKeyFile: Config.GetString("ssh.privatekeyfile"),
			LogDir:         Config.GetString("general.logdir"),
			SshHostKeyConfig: befehl.SshHostKeyConfig{
				Enabled:        Config.GetBool("ssh.hostkeyverificationenabled"),
				KnownHostsPath: Config.GetString("ssh.knownhostspath"),
			},
			RunConfigPath: runConfig,
		})
		if err != nil {
			panic(err)
		}

		if err := instance.Execute(routines); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(executeCmd)

	executeCmd.Flags().String("runconfig", "", "file location to the runtime configuration")
	executeCmd.Flags().Int("routines", 0, "maximum number of payloads that will run at once (defaults to 30)")

	executeCmd.MarkFlagRequired("runconfig")
}
