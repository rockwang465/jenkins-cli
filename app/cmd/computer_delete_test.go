package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/rockwang465/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rockwang465/jenkins-cli/mock/mhttp"
)

var _ = Describe("create delete command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          io.Writer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		computerDeleteOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		var (
			err error
		)

		BeforeEach(func() {
			var data []byte
			data, err = GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("should success", func() {
			name := "fake-name"

			client.PrepareForComputerDeleteRequest(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", name)

			rootCmd.SetArgs([]string{"computer", "delete", name})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
