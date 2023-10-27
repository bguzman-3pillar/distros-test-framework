//go:build upgradeselinux

package selinux

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/customflag"
	"github.com/rancher/distros-test-framework/pkg/testcase"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test:", func() {

	It("Start Up with no issues", func() {
		testcase.TestBuildCluster(GinkgoT())
	})

	It("Validate Nodes pre upgrade", func() {
		testcase.TestNodeStatus(
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validate Pods pre upgrade", func() {
		testcase.TestPodStatus(
			assert.PodAssertRestart(),
			assert.PodAssertReady(),
			assert.PodAssertStatus(),
		)
	})

	It("Validate selinux is enabled pre upgrade", func() {
		testcase.TestSelinuxEnabled(false)
	})

	It("Validate container, server and selinux version pre upgrade", func() {
		testcase.TestSelinux(false)
	})

	It("Validate container security pre upgrade", func() {
		testcase.TestSelinuxSpcT(false)
	})

	It("Validate context", func() {
		testcase.TestSelinuxContext(false)
	})

	It("Upgrade manual", func() {
		_ = testcase.TestUpgradeClusterManually(customflag.ServiceFlag.InstallMode.String())
	})

	It("Validate Nodes Post upgrade", func() {
		testcase.TestNodeStatus(
			assert.NodeAssertReadyStatus(),
			assert.NodeAssertVersionTypeUpgrade(customflag.ServiceFlag),
		)
	})

	It("Validate Pods Post upgrade", func() {
		testcase.TestPodStatus(
			assert.PodAssertRestart(),
			assert.PodAssertReady(),
			assert.PodAssertStatus(),
		)
	})

	It("Validate selinux is enabled Post upgrade", func() {
		testcase.TestSelinuxEnabled(false)
	})

	It("Validate container, server and selinux version Post upgrade", func() {
		testcase.TestSelinux(false)
	})

	It("Validate container security Post upgrade", func() {
		testcase.TestSelinuxSpcT(false)
	})

	It("Validate context Post Upgrade", func() {
		testcase.TestSelinuxContext(false)
	})

	It("Validate uninstall selinux policies Post Upgrade", func() {
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
