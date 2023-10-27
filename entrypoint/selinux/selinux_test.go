//go:build selinux

package selinux

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/testcase"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test:", func() {

	It("Start Up with no issues", func() {
		testcase.TestBuildCluster(GinkgoT())
	})

	It("Validate Nodes", func() {
		testcase.TestNodeStatus(
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validate Pods", func() {
		testcase.TestPodStatus(
			assert.PodAssertRestart(),
			assert.PodAssertReady(),
			assert.PodAssertStatus(),
		)
	})

	It("Validate selinux is enabled", func() {
		testcase.TestSelinuxEnabled(false)
	})

	It("Validate container, server and selinux version", func() {
		testcase.TestSelinux(false)
	})

	It("Validate container security", func() {
		testcase.TestSelinuxSpcT(false)
	})

	It("Validate context", func() {
		testcase.TestSelinuxContext(false)
	})

	It("Validate uninstall selinux policies", func() {
		testcase.TestUninstallPolicy(false)
	})

})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n", CurrentSpecReport().FullText())
	}
})
