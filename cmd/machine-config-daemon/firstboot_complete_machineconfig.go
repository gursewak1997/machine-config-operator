package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	daemon "github.com/openshift/machine-config-operator/pkg/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var firstbootCompleteMachineconfig = &cobra.Command{
	Use:                   "firstboot-complete-machineconfig",
	DisableFlagsInUseLine: true,
	Short:                 "Complete the host's initial boot into a MachineConfig",
	Args:                  cobra.MaximumNArgs(0),
	Run:                   executeFirstbootCompleteMachineConfig,
}

// init executes upon import
func init() {
	rootCmd.AddCommand(firstbootCompleteMachineconfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func runFirstBootCompleteMachineConfig(_ *cobra.Command, _ []string) error {
	flag.Set("logtostderr", "true")
	flag.Parse()

	exitCh := make(chan error)
	defer close(exitCh)

	dn, err := daemon.New(daemon.NewNodeUpdaterClient(), exitCh)
	if err != nil {
		return err
	}

	return dn.RunFirstbootCompleteMachineconfig()
}

func executeFirstbootCompleteMachineConfig(cmd *cobra.Command, args []string) {
	for {
		err := runFirstBootCompleteMachineConfig(cmd, args)
		if err != nil {
			glog.Warningf("error: %v\n", err)
			glog.Info("Sleeping 1 minute for retry")
			time.Sleep(time.Minute)
		} else {
			break
		}
	}
}
