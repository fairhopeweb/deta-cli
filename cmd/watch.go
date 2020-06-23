package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deta/deta-cli/runtime"
	"github.com/rjeczalik/notify"
	"github.com/spf13/cobra"
)

var (
	watchCmd = &cobra.Command{
		Use:   "watch [path]",
		Short: "Deploy changes in real time",
		RunE:  watch,
		Args:  cobra.MaximumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

func watch(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if len(args) != 0 {
		wd = args[0]
	}

	runtimeManager, err := runtime.NewManager(&wd, false)
	if err != nil {
		return err
	}

	isInitialized, err := runtimeManager.IsInitialized()
	if err != nil {
		return err
	}

	if !isInitialized {
		return fmt.Errorf(fmt.Sprintf("no deta micro initilialized in '%s'. see `deta new --help`", wd))
	}

	progInfo, err := runtimeManager.GetProgInfo()
	if err != nil {
		return err
	}

	c := make(chan notify.EventInfo, 1)

	// {dir}/... watch dir recursively
	if err := notify.Watch(filepath.Join(wd, "..."), c, notify.Write); err != nil {
		return err
	}

	fmt.Println("Watching changes")
	for {
		<-c
		err := deployChanges(runtimeManager, progInfo, true)
		if err != nil {
			return err
		}
	}
}
