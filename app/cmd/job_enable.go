package cmd

import (
	"github.com/rockwang465/jenkins-cli/app/cmd/common"
	"github.com/rockwang465/jenkins-cli/app/i18n"
	"github.com/rockwang465/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobEnableOption is the job delete option
type JobEnableOption struct {
	common.BatchOption
	common.Option
}

var jobEnableOption JobEnableOption

func init() {
	jobCmd.AddCommand(jobEnabelCmd)
	jobEnableOption.SetFlag(jobEnabelCmd)
}

var jobEnabelCmd = &cobra.Command{
	Use:   "enable",
	Short: i18n.T("Enable a job in your Jenkins"),
	Long:  i18n.T("Enable a job in your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jobName := args[0]
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobEnableOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		err = jclient.EnableJob(jobName)
		return
	},
}
