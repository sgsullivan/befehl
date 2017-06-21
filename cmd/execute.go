package cmd

import (
	"fmt"
	"github.com/sgsullivan/befehl"
	"github.com/spf13/cobra"
)

var payload string
var routines int
var hostsList string

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
[auth]
privatekeyfile = "/home/ssullivan/alt/id_rsa"
sshuser = "eingeben"

`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, value := range []string{payload, hostsList} {
			if value == "" {
				panic("Missing payload, or hosts; see --help")
			}
			if routines == 0 {
				fmt.Printf("--routines not given, defaulting to 30..\n")
				routines = 30
			}
		}
		befehl.Fire(&hostsList, &payload, &routines, Config)
	},
}

func init() {
	RootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringVarP(&payload, "payload", "", "", "file location to the payload, which contains the commands to execute on the remote hosts")
	executeCmd.Flags().StringVarP(&hostsList, "hosts", "", "", "file location to hosts list, which contains all hosts (separated by newline) to run the payload on")
	executeCmd.Flags().IntVarP(&routines, "routines", "", 0, "maximum number of payloads that will run at once (defaults to 30)")
}
