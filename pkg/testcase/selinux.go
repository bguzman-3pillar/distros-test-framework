package testcase

import (
	"fmt"
	"log"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/distros-test-framework/factory"
	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/shared"
)

var (
	cmdPrefix  = "sudo ls -laZ"
	ignoreDir  = "-I .. -I ."
	rke2       = "/var/lib/rancher/rke2"
	k3s        = "/var/lib/rancher/k3s"
	systemD    = "/etc/systemd/system"
	usrBin     = "/usr/bin"
	usrLocal   = "/usr/local/bin"
	grepFilter = "| grep -v \"/\""
)

const (
	ctxUnitFile = "system_u:object_r:container_unit_file_t:s0"
	ctxExec     = "system_u:object_r:container_runtime_exec_t:s0"
	ctxVarLib   = "system_u:object_r:container_var_lib_t:s0"
	ctxFile     = "system_u:object_r:container_file_t:s0"
	ctxConfig   = "system_u:object_r:container_config_t:s0"
	ctxShare    = "system_u:object_r:container_share_t:s0"
	ctxRoFile   = "system_u:object_r:container_ro_file_t:s0"
	ctxLog      = "system_u:object_r:container_log_t:s0"
	ctxRunTmpfs = "system_u:object_r:container_var_run_t:s0"
	ctxTmpfs    = "system_u:object_r:container_runtime_tmpfs_t:s0"
	ctxTLS      = "system_u:object_r:rke2_tls_t:s0"
	ctxLock     = "system_u:object_r:k3s_lock_t:s0"
	ctxData     = "system_u:object_r:k3s_data_t:s0"
	ctxRoot     = "system_u:object_r:k3s_root_t:s0"
	ctxNone     = "<<none>>"
	ctxRke2TLS  = "system_u:object_r:rke2_tls_t:s0"
)

type cmdCtx map[string]string

type configuration struct {
	distroName string
	cmdCtx
}

// TestSelinuxEnabled Validates that containerd is running with selinux enabled in the config
func TestSelinuxEnabled() {
	product, err := shared.GetProduct()
	if err != nil {
		return
	}

	ips := shared.FetchNodeExternalIP()
	selinuxConfigAssert := "selinux: true"
	selinuxContainerdAssert := "enable_selinux = true"

	for _, ip := range ips {
		err := assert.CheckComponentCmdNode("cat /etc/rancher/"+
			product+"/config.yaml", ip, selinuxConfigAssert)
		Expect(err).NotTo(HaveOccurred())
		errCont := assert.CheckComponentCmdNode("sudo cat /var/lib/rancher/"+
			product+"/agent/etc/containerd/config.toml", ip, selinuxContainerdAssert)
		Expect(errCont).NotTo(HaveOccurred())
	}
}

// TestSelinux Validates container-selinux version, rke2-selinux version and rke2-selinux version
func TestSelinux() {
	cluster := factory.AddCluster(GinkgoT())
	product, err := shared.GetProduct()
	if err != nil {
		return
	}

	var serverCmd string
	var serverAsserts []string
	agentAsserts := []string{"container-selinux", product + "-selinux"}

	switch product {
	case "k3s":
		serverCmd = "rpm -qa container-selinux k3s-selinux"
		serverAsserts = []string{"container-selinux", "k3s-selinux"}
	default:
		serverCmd = "rpm -qa container-selinux rke2-server rke2-selinux"
		serverAsserts = []string{"container-selinux", "rke2-selinux", "rke2-server"}
	}

	if cluster.NumServers > 0 {
		for _, serverIP := range cluster.ServerIPs {
			err := assert.CheckComponentCmdNode(serverCmd, serverIP, serverAsserts...)
			Expect(err).NotTo(HaveOccurred())
		}
	}

	if cluster.NumAgents > 0 {
		for _, agentIP := range cluster.AgentIPs {
			err := assert.CheckComponentCmdNode("rpm -qa container-selinux "+product+"-selinux", agentIP, agentAsserts...)
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

// https://github.com/k3s-io/k3s/blob/master/install.sh
// https://github.com/rancher/rke2/blob/master/install.sh
// Based on this info, this is the way to validate the correct context

// TestSelinuxContext Validates directories to ensure they have the correct selinux contexts created
func TestSelinuxContext() {
	cluster := factory.AddCluster(GinkgoT())
	product, err := shared.GetProduct()
	if err != nil {
		log.Println(err)
	}

	if cluster.NumServers > 0 {
		for _, ip := range cluster.ServerIPs {
			var context map[string]string
			context, err := getContext(product, ip)
			Expect(err).NotTo(HaveOccurred())

			fmt.Print("\nThese are the whole commands to use in this context validation\n")
			fmt.Print("Command to run || Context expected\n")
			for cmdsToRun, contExpected := range context {
				fmt.Println(cmdsToRun + " || " + contExpected)
			}

			for cmd, expectedContext := range context {
				res, err := shared.RunCommandOnNode(cmd, ip)
				fmt.Println("\nRunning cmd:", cmd, "\nExpected context:", expectedContext)
				fmt.Println("Result: \n", res)
				if res != "" {
					Expect(res).Should(ContainSubstring(expectedContext), "Error on cmd %v \n Context %v \nnot found on ", cmd, expectedContext, res)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		}
	}
}

func getVersion(osRelease, ip string) string {
	if strings.Contains(osRelease, "VERSION_ID") {
		res, err := shared.RunCommandOnNode("cat /etc/os-release | grep 'VERSION_ID'", ip)
		Expect(err).NotTo(HaveOccurred())
		parts := strings.Split(res, "=")
		if len(parts) == 2 {
			// Get version
			version := strings.Trim(parts[1], "\"")
			// if dot exist get the first number
			if dotIndex := strings.Index(version, "."); dotIndex != -1 {
				version = version[:dotIndex]
			}
			return version
		}
	}
	return ""
}

var osPolicy string

func getContext(product, ip string) (cmdCtx, error) {
	res, err := shared.RunCommandOnNode("cat /etc/os-release", ip)
	if err != nil {
		return nil, err
	}

	fmt.Println("OS Release: \n", res)
	policyMapping := map[string]string{
		"ID_LIKE='suse' VARIANT_ID='sle-micro'": "sle_micro",
		"ID_LIKE='suse'":                        "micro_os",
		"ID_LIKE='coreos'":                      "coreos",
		"VARIANT_ID='coreos'":                   "coreos",
	}

	for k, v := range policyMapping {
		if strings.Contains(res, k) {
			return selectSelinuxPolicy(product, v), nil
		}
	}

	version := getVersion(res, ip)
	versionMapping := map[string]string{
		"7": "centos7",
		"8": "centos8",
		"9": "centos9",
	}

	if policy, ok := versionMapping[version]; ok {
		return selectSelinuxPolicy(product, policy), nil
	}

	return nil, fmt.Errorf("unable to determine policy for %s on os: %s", ip, res)
}

func selectSelinuxPolicy(product, osType string) cmdCtx {
	key := fmt.Sprintf("%s_%s", product, osType)

	for _, config := range conf {
		if config.distroName == key {
			fmt.Printf("\nUsing '%s' policy for this %s cluster.\n", osType, product)
			osPolicy = osType
			return config.cmdCtx
		}
	}

	fmt.Printf("Configuration for %s not found!\n", key)
	return nil
}

// TestSelinuxSpcT Validate that containers don't run with spc_t
func TestSelinuxSpcT() {
	cluster := factory.AddCluster(GinkgoT())

	for _, serverIP := range cluster.ServerIPs {
		res, err := shared.RunCommandOnNode("ps auxZ | grep metrics | grep -v grep", serverIP)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).ShouldNot(ContainSubstring("spc_t"))
	}
}

// TestUninstallPolicy Validate that un-installation will remove the rke2-selinux or k3s-selinux policy
func TestUninstallPolicy() {
	product, err := shared.GetProduct()
	//product, err := shared.GetProduct()
	if err != nil {
		log.Println(err)
	}
	cluster := factory.AddCluster(GinkgoT())
	var serverUninstallCmd string
	var agentUninstallCmd string
	var serverCmd string

	switch product {
	case "k3s":
		serverUninstallCmd = "k3s-uninstall.sh"
		agentUninstallCmd = "k3s-agent-uninstall.sh"
		serverCmd = "rpm -qa container-selinux k3s-selinux"

	default:
		serverUninstallCmd = "sudo rke2-uninstall.sh"
		agentUninstallCmd = "sudo rke2-uninstall.sh"
		serverCmd = "rpm -qa container-selinux rke2-server rke2-selinux"
	}

	for _, serverIP := range cluster.ServerIPs {
		fmt.Println("Uninstalling "+product+" on server: ", serverIP)

		_, err := shared.RunCommandOnNode(serverUninstallCmd, serverIP)
		Expect(err).NotTo(HaveOccurred())

		res, errSel := shared.RunCommandOnNode(serverCmd, serverIP)
		Expect(errSel).NotTo(HaveOccurred())

		if strings.Contains(osPolicy, "centos7") {
			Expect(res).Should(ContainSubstring("container-selinux"))
			Expect(res).ShouldNot(ContainSubstring(product + "-selinux"))
		} else {
			Expect(res).Should(BeEmpty())
		}

	}

	for _, agentIP := range cluster.AgentIPs {
		fmt.Println("Uninstalling "+product+" on agent: ", agentIP)

		_, err := shared.RunCommandOnNode(agentUninstallCmd, agentIP)
		Expect(err).NotTo(HaveOccurred())

		res, errSel := shared.RunCommandOnNode("rpm -qa container-selinux "+product+"-selinux", agentIP)
		Expect(errSel).NotTo(HaveOccurred())

		if osPolicy == "centos7" {
			Expect(res).Should(ContainSubstring("container-selinux"))
			Expect(res).ShouldNot(ContainSubstring(product + "-selinux"))
		} else {
			Expect(res).Should(BeEmpty())
		}
	}
}

var conf = []configuration{
	{
		distroName: "rke2_centos7",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/rke2*":                                                            ctxUnitFile,
			cmdPrefix + " " + "/lib" + systemD + "/rke2*":                                                   ctxUnitFile,
			cmdPrefix + " " + usrLocal + "/lib" + systemD + "/rke2*":                                        ctxUnitFile,
			cmdPrefix + " " + usrBin + "/rke2":                                                              ctxExec,
			cmdPrefix + " " + usrLocal + "/rke2":                                                            ctxExec,
			cmdPrefix + " " + "/var/lib/cni " + ignoreDir:                                                   ctxVarLib,
			cmdPrefix + " " + "/var/lib/cni/* " + ignoreDir:                                                 ctxVarLib,
			cmdPrefix + " " + "/opt/cni " + ignoreDir:                                                       ctxFile,
			cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                                                     ctxFile,
			cmdPrefix + " " + "/var/lib/kubelet/pods " + ignoreDir:                                          ctxFile,
			cmdPrefix + " " + "/var/lib/kubelet/pods/* " + ignoreDir:                                        ctxFile,
			cmdPrefix + " " + rke2 + " " + ignoreDir:                                                        ctxVarLib,
			cmdPrefix + " " + rke2 + "/* " + ignoreDir:                                                      ctxVarLib,
			cmdPrefix + " " + rke2 + "/data":                                                                ctxExec,
			cmdPrefix + " " + rke2 + "/data/*":                                                              ctxExec,
			cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir + " " + grepFilter:                       ctxConfig,
			cmdPrefix + " " + rke2 + "/data/*/charts/* " + ignoreDir + " " + grepFilter:                     ctxConfig,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + " " + grepFilter:               ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                                            ctxLog,
			cmdPrefix + " " + rke2 + "/server/logs/ " + ignoreDir:                                           ctxLog,
			cmdPrefix + " " + "/var/run/flannel " + ignoreDir:                                               ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/flannel/* " + ignoreDir:                                             ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                   ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                 ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:   ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter: ctxTmpfs,
			cmdPrefix + " " + "/var/log/containers " + ignoreDir:                                            ctxLog,
			cmdPrefix + " " + "/var/log/containers/* " + ignoreDir:                                          ctxLog,
			cmdPrefix + " " + "/var/log/pods " + ignoreDir:                                                  ctxLog,
			cmdPrefix + " " + "/var/log/pods/* " + ignoreDir:                                                ctxLog,
			cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                                             ctxTLS,
			cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                                           ctxTLS,
		},
	},
	{
		// Works partially, has a bug related
		distroName: "rke2_centos8",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/rke2*": ctxUnitFile,
			// TODO: issue related to UnitFile https://github.com/rancher/rke2/issues/4741
			//cmdPrefix + " " + "/lib/systemd/system/rke2*":                                              ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/rke2*":                                    ctxUnitFile,
			cmdPrefix + " " + usrBin + "/rke2":                                                         ctxExec,
			cmdPrefix + " " + usrLocal + "/rke2":                                                       ctxExec,
			cmdPrefix + " " + "/opt/cni " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + rke2 + " " + ignoreDir:                                                   ctxVarLib,
			cmdPrefix + " " + rke2 + "/* " + ignoreDir:                                                 ctxVarLib,
			cmdPrefix + " " + rke2 + "/data " + ignoreDir:                                              ctxExec,
			cmdPrefix + " " + rke2 + "/data/* " + ignoreDir:                                            ctxExec,
			cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir + " " + grepFilter:                  ctxConfig,
			cmdPrefix + " " + rke2 + "/data/*/charts/* " + ignoreDir + " " + grepFilter:                ctxConfig,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:   ctxFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter: ctxFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + " " + grepFilter:          ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:   ctxRoFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter: ctxRoFile,
			cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                                       ctxLog,
			cmdPrefix + " " + rke2 + "/server/logs/* " + ignoreDir:                                     ctxLog,
			cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                                        ctxTLS,
			cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                                      ctxTLS,
		},
	},
	{
		// Works partially, has a bug related
		distroName: "rke2_centos9",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/rke2*": ctxUnitFile,
			// TODO: issue related to UnitFile https://github.com/rancher/rke2/issues/4741
			//cmdPrefix + " " + "/lib/systemd/system/rke2*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/rke2*":                                       ctxUnitFile,
			cmdPrefix + " " + usrBin + "/rke2":                                                            ctxExec,
			cmdPrefix + " " + usrLocal + "/rke2":                                                          ctxExec,
			cmdPrefix + " " + "/opt/cni " + ignoreDir:                                                     ctxFile,
			cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                                                   ctxFile,
			cmdPrefix + " " + rke2 + " " + ignoreDir:                                                      ctxVarLib,
			cmdPrefix + " " + rke2 + "/* " + ignoreDir:                                                    ctxVarLib,
			cmdPrefix + " " + rke2 + "/data " + ignoreDir:                                                 ctxExec,
			cmdPrefix + " " + rke2 + "/data/* " + ignoreDir:                                               ctxExec,
			cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir + " " + grepFilter:                     ctxConfig,
			cmdPrefix + " " + rke2 + "/data/*/charts/* " + ignoreDir + " " + grepFilter:                   ctxConfig,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:      ctxFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/ " + ignoreDir + " " + grepFilter:     ctxFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:      ctxRoFile,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:    ctxRoFile,
			cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                                          ctxLog,
			cmdPrefix + " " + rke2 + "/server/logs/* " + ignoreDir:                                        ctxLog,
			cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                                           ctxTLS,
			cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                                         ctxTLS,
		},
	},
	{
		// TODO: We are not able to execute this because our framework does not support the reboot part for this OS.
		distroName: "rke2_micro_os",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/rke2*":                                                          ctxUnitFile,
			cmdPrefix + " " + "/lib/systemd/system/rke2*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/rke2*":                                       ctxUnitFile,
			cmdPrefix + " " + usrBin + "/rke2":                                                            ctxExec,
			cmdPrefix + " " + usrLocal + "/rke2":                                                          ctxExec,
			cmdPrefix + " " + "/opt/cni " + ignoreDir:                                                     ctxFile,
			cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                                                   ctxFile,
			cmdPrefix + " " + rke2 + " " + ignoreDir:                                                      ctxVarLib,
			cmdPrefix + " " + rke2 + "/* " + ignoreDir:                                                    ctxVarLib,
			cmdPrefix + " " + rke2 + "/data " + ignoreDir:                                                 ctxExec,
			cmdPrefix + " " + rke2 + "/data/* " + ignoreDir:                                               ctxExec,
			cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir + " " + grepFilter:                     ctxConfig,
			cmdPrefix + " " + rke2 + "/data/*/charts/* " + ignoreDir + " " + grepFilter:                   ctxConfig,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/ " + ignoreDir + " " + grepFilter:     ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:    ctxShare,
			cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                                          ctxLog,
			cmdPrefix + " " + rke2 + "/server/logs/* " + ignoreDir:                                        ctxLog,
			cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                                           ctxRke2TLS,
			cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                                         ctxRke2TLS,
		},
	},
	{
		// TODO: We are not able to execute this because our framework does not support the reboot part for this OS.
		distroName: "rke2_sle_micro",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/rke2*":                                                          ctxUnitFile,
			cmdPrefix + " " + "/lib/systemd/system/rke2*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/rke2.*":                                      ctxUnitFile,
			cmdPrefix + " " + usrBin + "/rke2":                                                            ctxExec,
			cmdPrefix + " " + usrLocal + "/rke2":                                                          ctxExec,
			cmdPrefix + " " + "/opt/rke2/bin/rke2":                                                        ctxExec,
			cmdPrefix + " " + "/opt/cni " + ignoreDir:                                                     ctxFile,
			cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                                                   ctxFile,
			cmdPrefix + " " + rke2 + " " + ignoreDir:                                                      ctxVarLib,
			cmdPrefix + " " + rke2 + "/* " + ignoreDir:                                                    ctxVarLib,
			cmdPrefix + " " + rke2 + "/data " + ignoreDir:                                                 ctxExec,
			cmdPrefix + " " + rke2 + "/data/*" + ignoreDir:                                                ctxExec,
			cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir + " " + grepFilter:                     ctxConfig,
			cmdPrefix + " " + rke2 + "/data/*/charts/* " + ignoreDir + " " + grepFilter:                   ctxConfig,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:    ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:      ctxShare,
			cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:    ctxShare,
			cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                                          ctxLog,
			cmdPrefix + " " + rke2 + "/server/logs/* " + ignoreDir:                                        ctxLog,
			cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                                           ctxTLS,
			cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                                         ctxTLS,
		},
	},
	{
		// Works partially, has a bug related and some different outputs
		distroName: "k3s_centos7",
		cmdCtx: cmdCtx{
			// TODO: issue related to UnitFile  https://github.com/k3s-io/k3s/issues/8317
			//cmdPrefix + " " + systemD + "/k3s*":                      ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":         ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*":   ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                       ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":                 ctxExec,
			cmdPrefix + " " + "/var/lib/cni " + ignoreDir:            ctxVarLib,
			cmdPrefix + " " + "/var/lib/cni/* " + ignoreDir:          ctxVarLib,
			cmdPrefix + " " + "/var/lib/kubelet/pods " + ignoreDir:   ctxFile,
			cmdPrefix + " " + "/var/lib/kubelet/pods/* " + ignoreDir: ctxFile,
			/* TODO: Here the expected output is "system_u:object_r:container_var_lib_t:s0"
			and is showing this "unconfined_u:object_r:container_var_lib_t:s0" (user part is not the expected)*/
			//cmdPrefix + " " + k3s + " " + ignoreDir:                                                          ctxVarLib,
			cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                        ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter:     ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/data " + ignoreDir:                                                     ctxData,
			cmdPrefix + " " + k3s + "/data/* " + ignoreDir:                                                   ctxData,
			cmdPrefix + " " + k3s + "/data/.lock":                                                            ctxLock,
			cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter:                            ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter:                            ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter:                          ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/log/containers " + ignoreDir:                                             ctxLog,
			cmdPrefix + " " + "/var/log/containers/* " + ignoreDir:                                           ctxLog,
			cmdPrefix + " " + "/var/log/pods " + ignoreDir:                                                   ctxLog,
			cmdPrefix + " " + "/var/log/pods/* " + ignoreDir:                                                 ctxLog,
			cmdPrefix + " " + "/var/run/flannel " + ignoreDir:                                                ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/flannel/* " + ignoreDir:                                              ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
	{
		distroName: "k3s_centos8",
		cmdCtx: cmdCtx{
			// TODO: issue related to UnitFile  https://github.com/k3s-io/k3s/issues/8317
			//cmdPrefix + " " + systemD + "/k3s*":                                                          ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":       ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*": ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                     ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":               ctxExec,
			/* TODO: Expected context "system_u:object_r:container_var_lib_t:s0" and is showing "unconfined_u:object_r:container_var_lib_t:s0" */
			//cmdPrefix + " " + k3s + " " + ignoreDir:                                                      ctxVarLib,
			//cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                    ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:      ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:    ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:      ctxRoFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:    ctxRoFile,

			/* TODO: Expected context "system_u:object_r:k3s_data_t:s0" and is showing "unconfined_u:object_r:k3s_lock_t:s0"*/
			//cmdPrefix + " " + k3s + "/data " + ignoreDir:   ctxData,
			//cmdPrefix + " " + k3s + "/data/* " + ignoreDir: ctxData,

			/* TODO: Expected context is "system_u:object_r:k3s_lock_t:s0" and is showing "unconfined_u:object_r:k3s_lock_t:s0" */
			//cmdPrefix + " " + k3s + "/data/.lock":                                 ctxLock,

			/* TODO: For these directories output shows "unconfined_u:object_r:k3s_root_t:s0"	and the expected one is "system_u:object_r:k3s_root_t:s0"*/
			//cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter: ctxRoot,
			//cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,

			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter + " | grep -v 'total 0'":   ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter + " | grep -v 'total 0'": ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
	{

		distroName: "k3s_centos9",
		cmdCtx: cmdCtx{
			// TODO: issue related to UnitFile  https://github.com/k3s-io/k3s/issues/8317
			//cmdPrefix + " " + systemD + "/k3s*":                                                          ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":       ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*": ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                     ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":               ctxExec,
			// TODO: Output: unconfined_u Expected: system_u
			// cmdPrefix + " " + k3s + " " + ignoreDir:                                                      ctxVarLib,
			// cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                    ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:      ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:    ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter: ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:      ctxRoFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:    ctxRoFile,
			/* TODO: Expected "system_u:object_r:k3s_data_t:s0" and is showing "unconfined_u:object_r:k3s_data_t:s0" */
			//cmdPrefix + " " + k3s + "/data " + ignoreDir:                                                 ctxData,
			//cmdPrefix + " " + k3s + "/data/* " + ignoreDir:                                               ctxData,

			/* TODO: Expected "system_u:object_r:k3s_lock_t:s0 " and is showing "unconfined_u:object_r:k3s_lock_t:s0"*/
			//cmdPrefix + " " + k3s + "/data/.lock":                                                            ctxLock,

			/* TODO: Expected "system_u:object_r:k3s_root_t:s0 " and is showing "unconfined_u:object_r:k3s_root_t:s0" */
			//cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter:                            ctxRoot,
			//cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,

			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter + " | grep -v 'total 0'":   ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter + " | grep -v 'total 0'": ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
	{
		// TODO: We are not able to execute this because our framework does not support the reboot part for this OS.
		distroName: "k3s_coreos",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/k3s*":                                                              ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*":                                           ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                                                               ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":                                                         ctxExec,
			cmdPrefix + " " + k3s + " " + ignoreDir:                                                          ctxVarLib,
			cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                        ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:          ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:        ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter:     ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/data " + ignoreDir:                                                     ctxData,
			cmdPrefix + " " + k3s + "/data/* " + ignoreDir:                                                   ctxData,
			cmdPrefix + " " + k3s + "/data/.lock":                                                            ctxLock,
			cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter:                            ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter:                            ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter:                          ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
	{
		// TODO: We are not able to execute this because our framework does not support the reboot part for this OS.
		distroName: "k3s_micro_os",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/k3s*":                                                              ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*":                                           ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                                                               ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":                                                         ctxExec,
			cmdPrefix + " " + k3s + " " + ignoreDir:                                                          ctxVarLib,
			cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                        ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:          ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:        ctxFile,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter:     ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/data " + ignoreDir:                                                     ctxData,
			cmdPrefix + " " + k3s + "/data/* " + ignoreDir:                                                   ctxData,
			cmdPrefix + " " + k3s + "/data/.lock":                                                            ctxLock,
			cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter:                            ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter:                            ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter:                          ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
	{
		// TODO: We are not able to execute this because our framework does not support the reboot part for this OS.
		distroName: "k3s_sle_micro",
		cmdCtx: cmdCtx{
			cmdPrefix + " " + systemD + "/k3s*":                                                              ctxUnitFile,
			cmdPrefix + " " + "/usr/lib/systemd/system/k3s*":                                                 ctxUnitFile,
			cmdPrefix + " " + "/usr/local/lib/systemd/system/k3s*":                                           ctxUnitFile,
			cmdPrefix + " " + "/usr/s?bin/k3s":                                                               ctxExec,
			cmdPrefix + " " + "/usr/local/s?bin/k3s":                                                         ctxExec,
			cmdPrefix + " " + k3s + " " + ignoreDir:                                                          ctxVarLib,
			cmdPrefix + " " + k3s + "/* " + ignoreDir:                                                        ctxVarLib,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/snapshots/*/.* " + ignoreDir + " " + grepFilter:     ctxNone,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes " + ignoreDir + " " + grepFilter:          ctxShare,
			cmdPrefix + " " + k3s + "/agent/containerd/*/sandboxes/* " + ignoreDir + " " + grepFilter:        ctxShare,
			cmdPrefix + " " + k3s + "/data " + ignoreDir:                                                     ctxData,
			cmdPrefix + " " + k3s + "/data/* " + ignoreDir:                                                   ctxData,
			cmdPrefix + " " + k3s + "/data/.lock":                                                            ctxLock,
			cmdPrefix + " " + k3s + "/data/*/bin " + ignoreDir + " " + grepFilter:                            ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/* " + ignoreDir + " " + grepFilter:                          ctxRoot,
			cmdPrefix + " " + k3s + "/data/*/bin/.*links " + ignoreDir + " " + grepFilter:                    ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/.*sha256sums " + ignoreDir + " " + grepFilter:               ctxData,
			cmdPrefix + " " + k3s + "/data/*/bin/cni " + ignoreDir + " " + grepFilter:                        ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd " + ignoreDir + " " + grepFilter:                 ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim " + ignoreDir + " " + grepFilter:            ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/containerd-shim-runc-v[12] " + ignoreDir + " " + grepFilter: ctxExec,
			cmdPrefix + " " + k3s + "/data/*/bin/runc " + ignoreDir + " " + grepFilter:                       ctxExec,
			cmdPrefix + " " + k3s + "/data/*/etc " + ignoreDir + " " + grepFilter:                            ctxConfig,
			cmdPrefix + " " + k3s + "/data/*/etc/* " + ignoreDir + " " + grepFilter:                          ctxConfig,
			cmdPrefix + " " + k3s + "/storage " + ignoreDir:                                                  ctxFile,
			cmdPrefix + " " + k3s + "/storage/* " + ignoreDir:                                                ctxFile,
			cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                                                    ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                                                  ctxRunTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm " + ignoreDir + " " + grepFilter:    ctxTmpfs,
			cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/* " + ignoreDir + " " + grepFilter:  ctxTmpfs,
		},
	},
}
