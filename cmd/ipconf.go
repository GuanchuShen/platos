package cmd

import (
	"github.com/shenguanchu/platos/ipconf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand()
}

var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: IpConfHandle,
}

func IpConfHandle(cmd *cobra.Command, args []string) {
	ipconf.RunMain()
}
