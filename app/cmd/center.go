package cmd

import (
	"fmt"
	"github.com/rockwang465/jenkins-cli/app/cmd/common"
	"net/http"

	appCfg "github.com/rockwang465/jenkins-cli/app/config"
	"github.com/rockwang465/jenkins-cli/app/i18n"

	"github.com/rockwang465/jenkins-cli/app/helper"

	"github.com/rockwang465/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterOption is the center cmd option
type CenterOption struct {
	common.WatchOption

	RoundTripper http.RoundTripper
	CenterStatus string
}

var centerOption CenterOption

func init() {
	rootCmd.AddCommand(centerCmd)
}

var centerCmd = &cobra.Command{
	Use:   "center",
	Short: i18n.T("Manage your update center"),
	Long:  i18n.T("Manage your update center"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins, cmd, centerOption.RoundTripper)

		_, err = printUpdateCenter(jenkins, cmd, centerOption.RoundTripper)
		return
	},
}

func printUpdateCenter(jenkins *appCfg.JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) (
	status *client.UpdateCenter, err error) {
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

	var centerStatus string
	if status, err = jclient.Status(); err == nil {
		centerStatus += fmt.Sprintf("RestartRequiredForCompletion: %v\n", status.RestartRequiredForCompletion)
		if status.Jobs != nil {
			for i, job := range status.Jobs {
				if job.Type == "InstallationJob" {
					centerStatus += fmt.Sprintf("%d, %s, %s, %v, %s\n", i, job.Type, job.Name, job.Status, job.ErrorMessage)
				} else {
					centerStatus += fmt.Sprintf("%d, %s, %s\n", i, job.Type, job.ErrorMessage)
				}
			}
		}

		if centerOption.CenterStatus != centerStatus {
			centerOption.CenterStatus = centerStatus

			cmd.Printf("%s", centerStatus)
		}
	}
	return
}

func printJenkinsStatus(jenkins *appCfg.JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) {
	jclient := &client.JenkinsStatusClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth

	status, err := jclient.Get()
	if err == nil {
		cmd.Println("Jenkins Version:", status.Version)
	}
	helper.CheckErr(cmd, err)
}
