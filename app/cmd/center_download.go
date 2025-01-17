package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rockwang465/jenkins-cli/app/i18n"

	"github.com/rockwang465/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterDownloadOption as the options of download command
type CenterDownloadOption struct {
	LTS     bool
	Mirror  string
	Version string

	Output       string
	ShowProgress bool

	Formula string

	Thread       int
	RoundTripper http.RoundTripper
}

var centerDownloadOption CenterDownloadOption

func init() {
	centerCmd.AddCommand(centerDownloadCmd)
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Version, "war-version", "", "",
		i18n.T("Version of the Jenkins which you want to download"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Mirror, "mirror", "m", "default",
		i18n.T("The mirror site of Jenkins"))
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.ShowProgress, "progress", "p", true,
		i18n.T("If you want to show the download progress"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Output, "output", "o", "",
		i18n.T("The file of output"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Formula, "formula", "", "",
		i18n.T("The formula of jenkins.war, only support zh currently"))
	centerDownloadCmd.Flags().IntVarP(&centerDownloadOption.Thread, "thread", "t", 0,
		"Using multi-thread to download jenkins.war")
}

var centerDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: i18n.T("Download jenkins.war"),
	Long: i18n.T(`Download jenkins.war from a mirror site.
You can get more mirror sites from https://jenkins-zh.cn/tutorial/management/mirror/
If you want to download different formulas of jenkins.war, please visit the following project
https://github.com/jenkins-zh/jenkins-formulas`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		if err = centerDownloadOption.DownloadJenkins(); err == nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "download success, %s\n", centerDownloadOption.Output)
		}
		return err
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if centerDownloadOption.Output != "" {
			return
		}

		var userHome string
		if userHome, err = homedir.Dir(); err != nil {
			return
		}

		centerDownloadOption.Output = fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, centerDownloadOption.Version)
		return
	},
}

// DownloadJenkins download the Jenkins
func (c *CenterDownloadOption) DownloadJenkins() (err error) {
	parentDir := filepath.Dir(c.Output)
	if err = os.MkdirAll(parentDir, os.FileMode(0755)); err != nil {
		return
	}

	mirrorSite := getMirror(c.Mirror)
	if mirrorSite == "" {
		err = fmt.Errorf("cannot found Jenkins mirror by: %s", c.Mirror)
		return
	}

	jClient := &client.UpdateCenterManager{
		MirrorSite: mirrorSite,
		JenkinsCore: client.JenkinsCore{
			RoundTripper: c.RoundTripper,
		},
		Formula:      c.Formula,
		LTS:          c.LTS,
		Version:      c.Version,
		Output:       c.Output,
		ShowProgress: c.ShowProgress,
		Thread:       c.Thread,
	}

	jenkinsWarURL := jClient.GetJenkinsWarURL()
	logger.Info("prepare to download jenkins.war", zap.String("URL", jenkinsWarURL))

	err = jClient.DownloadJenkins()
	return
}
