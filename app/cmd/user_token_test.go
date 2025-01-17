package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/rockwang465/jenkins-cli/client"
	"github.com/rockwang465/jenkins-cli/mock/mhttp"
)

var _ = Describe("user token command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		err          error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		userTokenOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		var data []byte
		data, err = GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("with http requests", func() {
		It("lack of arguments", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			userTokenCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})

			rootCmd.SetArgs([]string{"user", "token"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("help"))
		})

		It("should success", func() {
			targetUser := "target-user"

			tokenName := "fakename"
			client.PrepareCreateToken(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", tokenName, targetUser)

			rootCmd.SetArgs([]string{"user", "token", "-g", "-n", tokenName, "--target-user", targetUser})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("{\n  \"Status\": \"ok\",\n  \"Data\": {\n    \"TokenName\": \"\",\n    \"TokenUUID\": \"\",\n    \"TokenValue\": \"\"\n  }\n}\n"))
		})

		It("with status code 500", func() {
			targetUser := "target-user"

			tokenName := "fakename"
			response := client.PrepareCreateToken(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", tokenName, targetUser)
			response.StatusCode = 500

			rootCmd.SetArgs([]string{"user", "token", "-g", "-n", tokenName, "--target-user", targetUser})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())

			Expect(buf.String()).To(ContainSubstring("unexpected status code: 500"))
		})
	})
})
