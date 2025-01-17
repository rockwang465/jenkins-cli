package cmd

import (
	"github.com/rockwang465/jenkins-cli/util"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("shell command", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""

		tempFile, err := ioutil.TempFile("", "example")
		Expect(err).To(BeNil())

		rootOptions.ConfigFile = tempFile.Name()
		shellOptions.ExecContext = util.FakeExecCommandSuccess
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	//Context("basic test", func() {
	//	It("should success", func() {
	//		data, err := GenerateSampleConfig()
	//		Expect(err).To(BeNil())
	//		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
	//		Expect(err).To(BeNil())
	//
	//		rootCmd.SetArgs([]string{"shell", "yourServer"})
	//
	//		buf := new(bytes.Buffer)
	//		rootCmd.SetOutput(buf)
	//		_, err = rootCmd.ExecuteC()
	//		Expect(err).To(BeNil())
	//
	//		Expect(buf.String()).To(ContainSubstring("testing: warning: no tests to run\nPASS\n"))
	//	})
	//})
})
