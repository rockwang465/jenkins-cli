package cmd

import (
	"encoding/json"
	"github.com/rockwang465/jenkins-cli/app/cmd/common"
	"github.com/rockwang465/jenkins-cli/app/i18n"

	"github.com/rockwang465/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterIdentityOption option for upgrade Jenkins
type CenterIdentityOption struct {
	common.Option
}

var centerIdentityOption CenterIdentityOption

func init() {
	centerCmd.AddCommand(centerIdentityCmd)
}

var centerIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: i18n.T("Print the identity of current Jenkins"),
	Long:  i18n.T("Print the identity of current Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.CoreClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerIdentityOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var identity client.JenkinsIdentity
		var data []byte
		if identity, err = jClient.GetIdentity(); err == nil {
			if data, err = json.MarshalIndent(identity, "", " "); err == nil {
				cmd.Println(string(data))
			}
		}
		return
	},
}
