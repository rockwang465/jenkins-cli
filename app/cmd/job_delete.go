package cmd

import (
	"fmt"
	"github.com/rockwang465/jenkins-cli/app/cmd/common"
	"github.com/rockwang465/jenkins-cli/app/i18n"
	"github.com/rockwang465/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobDeleteOption is the job delete option
type JobDeleteOption struct {
	common.BatchOption
	common.Option
}

var jobDeleteOption JobDeleteOption

func init() {
	jobCmd.AddCommand(jobDeleteCmd)
	jobDeleteOption.SetFlag(jobDeleteCmd)
	jobDeleteOption.BatchOption.Stdio = common.GetSystemStdio()
	jobDeleteOption.Option.Stdio = common.GetSystemStdio()
}

var jobDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: common.GetAliasesDel(),
	Short:   i18n.T("Delete a job in your Jenkins"),
	Long:    i18n.T("Delete a job in your Jenkins"),
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jobName := args[0]
		if !jobDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete job %s ?", jobName)) {
			return
		}

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobDeleteOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		err = jclient.Delete(jobName)
		return
	},
}
