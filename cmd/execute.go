package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/sgsullivan/befehl"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute the given payload against the given hosts list",
	Long: `Executes the given payload on each host in the hosts list. Hosts in the hosts
list should be separated by a new line. You can control how many payloads run concurrently by
passing the routines flag. 

By default befehl will use the private key in $HOME/.ssh/id_rsa. This can be overrode by
specifying auth.privatekeyfile in ~/.befehl.[toml|json|yaml].

By default befehl will write the output of each payload for each host in $HOME/befehl/logs. This
can be overrode by specifying general.logdir in ~/.befehl.[toml|json|yaml].

By default befehl will attempt to ssh as root. This can be overrode by specifying auth.sshuser
in ~/.befehl.[toml|json|yaml].

Heres an example specifying all of the above mentioned options:

[general]
logdir = "/home/ssullivan/log-special"
[ssh]
privatekeyfile = "/home/ssullivan/alt/id_rsa"
user = "eingeben"
knownhostspath = "/home/asullivan/alt/.ssh/known_hosts"
hostkeyverificationenabled = true

`,
	Run: func(cmd *cobra.Command, args []string) {
		hostsFile, _ := cmd.Flags().GetString("hosts")
		payload, _ := cmd.Flags().GetString("payload")
		routines, _ := cmd.Flags().GetInt("routines")

		if routines == 0 {
			color.Yellow("--routines not given, defaulting to 30..\n")
			routines = 30
		}

		instance := befehl.New(&befehl.Options{
			PrivateKeyFile: Config.GetString("ssh.privatekeyfile"),
			SshUser:        Config.GetString("ssh.sshuser"),
			LogDir:         Config.GetString("general.logdir"),
			SshHostKeyConfig: befehl.SshHostKeyConfig{
				Enabled:        Config.GetBool("ssh.hostkeyverificationenabled"),
				KnownHostsPath: Config.GetString("ssh.knownhostspath"),
			},
		})

		if err := instance.Execute(hostsFile, payload, routines); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(executeCmd)

	executeCmd.Flags().String("payload", "", "file location to the payload, which contains the commands to execute on the remote hosts")
	executeCmd.Flags().String("hosts", "", "file location to hosts list, which contains all hosts (separated by newline) to run the payload on")
	executeCmd.Flags().Int("routines", 0, "maximum number of payloads that will run at once (defaults to 30)")

	executeCmd.MarkFlagRequired("payload")
	executeCmd.MarkFlagRequired("hosts")
}
