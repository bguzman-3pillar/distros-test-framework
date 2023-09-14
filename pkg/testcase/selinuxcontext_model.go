package testcase

var (
	cmdPrefix = "sudo ls -laZ"
	ignoreDir = "-I .. -I ."
	rke2      = "/var/lib/rancher/rke2"
	systemD   = "/etc/systemd/system"
	usrBin    = "/usr/bin"
	usrLocal  = "/usr/local/bin"
)

const (
	ctxUnitFile = "system_u:object_r:container_unit_file_t"
	ctxExec     = "system_u:object_r:container_runtime_exec_t"
	ctxVarLib   = "system_u:object_r:container_var_lib_t"
	ctxFile     = "system_u:object_r:container_file_t"
	ctxConfig   = "system_u:object_r:container_config_t"
	ctxShare    = "system_u:object_r:container_share_t"
	ctxLog      = "system_u:object_r:container_log_t"
	ctxRunTmpfs = "system_u:object_r:container_var_run_t"
	ctxTmpfs    = "system_u:object_r:container_runtime_tmpfs_t"
	ctxTLS      = "system_u:object_r:rke2_tls_t"
	ctxLock     = "system_u:object_r:k3s_lock_t"
	ctxData     = "system_u:object_r:k3s_data_t"
	ctxRoot     = "system_u:object_r:k3s_root_t"
	ctxNone     = "<<none>>"
)

type cmdCtx map[string]string

type configuration struct {
	distroName string
	cmdCtx
}

var conf = configuration{
	distroName: "rke2_centos7",
	cmdCtx: cmdCtx{
		cmdPrefix + " " + systemD + "/rke2*":                                   ctxUnitFile,
		cmdPrefix + " " + "/lib" + systemD + "/rke2*":                          ctxUnitFile,
		cmdPrefix + " " + usrLocal + "/lib" + systemD + "/rke2*":               ctxUnitFile,
		cmdPrefix + " " + usrBin + "/rke2":                                     ctxExec,
		cmdPrefix + " " + usrLocal + "/rke2":                                   ctxExec,
		cmdPrefix + " " + "/var/lib/cni " + ignoreDir:                          ctxVarLib,
		cmdPrefix + " " + "/var/lib/cni/* " + ignoreDir:                        ctxVarLib,
		cmdPrefix + " " + "/opt/cni " + ignoreDir:                              ctxFile,
		cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                            ctxFile,
		cmdPrefix + " " + "/var/lib/kubelet/pods " + ignoreDir:                 ctxFile,
		cmdPrefix + " " + "/var/lib/kubelet/pods/* " + ignoreDir:               ctxFile,
		cmdPrefix + " " + rke2 + " " + ignoreDir:                               ctxVarLib,
		cmdPrefix + " " + rke2 + "/* " + ignoreDir:                             ctxVarLib,
		cmdPrefix + " " + rke2 + "/data(/.*)?":                                 ctxExec,
		cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir:                 ctxConfig,
		cmdPrefix + " " + rke2 + "/data/*/charts/*":                            ctxConfig,
		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir:  ctxShare,
		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/ " + ignoreDir: ctxShare,
		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/[^/]*/.*":      ctxNone,
		cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir:  ctxShare,
		cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/ " + ignoreDir: ctxShare,
		cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                   ctxLog,
		cmdPrefix + " " + rke2 + "/server/logs/ " + ignoreDir:                  ctxLog,
		cmdPrefix + " " + "/var/run/flannel " + ignoreDir:                      ctxRunTmpfs,
		cmdPrefix + " " + "/var/run/flannel/* " + ignoreDir:                    ctxRunTmpfs,
		cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                          ctxRunTmpfs,
		cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                        ctxRunTmpfs,
		cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm":          ctxTmpfs,
		cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/*":        ctxTmpfs,
		cmdPrefix + " " + "/var/log/containers " + ignoreDir:                   ctxLog,
		cmdPrefix + " " + "/var/log/containers/* " + ignoreDir:                 ctxLog,
		cmdPrefix + " " + "/var/log/pods " + ignoreDir:                         ctxLog,
		cmdPrefix + " " + "/var/log/pods/* " + ignoreDir:                       ctxLog,
		cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                    ctxTLS,
		cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                  ctxTLS,
	},
}

// var data = map[string]map[string]string{
// 	"rke2_centos7": {
// 		cmdPrefix + " " + systemD + "/rke2*":                                   ctxUnitFile,
// 		cmdPrefix + " " + "/lib" + systemD + "/rke2*":                          ctxUnitFile,
// 		cmdPrefix + " " + usrLocal + "/lib" + systemD + "/rke2*":               ctxUnitFile,
// 		cmdPrefix + " " + usrBin + "/rke2":                                     ctxExec,
// 		cmdPrefix + " " + usrLocal + "/rke2":                                   ctxExec,
// 		cmdPrefix + " " + "/var/lib/cni " + ignoreDir:                          ctxVarLib,
// 		cmdPrefix + " " + "/var/lib/cni/* " + ignoreDir:                        ctxVarLib,
// 		cmdPrefix + " " + "/opt/cni " + ignoreDir:                              ctxFile,
// 		cmdPrefix + " " + "/opt/cni/* " + ignoreDir:                            ctxFile,
// 		cmdPrefix + " " + "/var/lib/kubelet/pods " + ignoreDir:                 ctxFile,
// 		cmdPrefix + " " + "/var/lib/kubelet/pods/* " + ignoreDir:               ctxFile,
// 		cmdPrefix + " " + rke2 + " " + ignoreDir:                               ctxVarLib,
// 		cmdPrefix + " " + rke2 + "/* " + ignoreDir:                             ctxVarLib,
// 		cmdPrefix + " " + rke2 + "/data(/.*)?":                                 ctxExec,
// 		cmdPrefix + " " + rke2 + "/data/*/charts " + ignoreDir:                 ctxConfig,
// 		cmdPrefix + " " + rke2 + "/data/*/charts/*":                            ctxConfig,
// 		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots " + ignoreDir:  ctxShare,
// 		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/ " + ignoreDir: ctxShare,
// 		cmdPrefix + " " + rke2 + "/agent/containerd/*/snapshots/[^/]*/.*":      ctxNone,
// 		cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes " + ignoreDir:  ctxShare,
// 		cmdPrefix + " " + rke2 + "/agent/containerd/*/sandboxes/ " + ignoreDir: ctxShare,
// 		cmdPrefix + " " + rke2 + "/server/logs " + ignoreDir:                   ctxLog,
// 		cmdPrefix + " " + rke2 + "/server/logs/ " + ignoreDir:                  ctxLog,
// 		cmdPrefix + " " + "/var/run/flannel " + ignoreDir:                      ctxRunTmpfs,
// 		cmdPrefix + " " + "/var/run/flannel/* " + ignoreDir:                    ctxRunTmpfs,
// 		cmdPrefix + " " + "/var/run/k3s " + ignoreDir:                          ctxRunTmpfs,
// 		cmdPrefix + " " + "/var/run/k3s/* " + ignoreDir:                        ctxRunTmpfs,
// 		cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm":          ctxTmpfs,
// 		cmdPrefix + " " + "/var/run/k3s/containerd/*/sandboxes/*/shm/*":        ctxTmpfs,
// 		cmdPrefix + " " + "/var/log/containers " + ignoreDir:                   ctxLog,
// 		cmdPrefix + " " + "/var/log/containers/* " + ignoreDir:                 ctxLog,
// 		cmdPrefix + " " + "/var/log/pods " + ignoreDir:                         ctxLog,
// 		cmdPrefix + " " + "/var/log/pods/* " + ignoreDir:                       ctxLog,
// 		cmdPrefix + " " + rke2 + "/server/tls " + ignoreDir:                    ctxTLS,
// 		cmdPrefix + " " + rke2 + "/server/tls/* " + ignoreDir:                  ctxTLS,
// 	},
// 	"rke2_centos8":   {},
// 	"rke2_centos9":   {},
// 	"rke2_micro_os":  {},
// 	"rke2_sle_micro": {},
// 	"k3s_centos7":    {},
// 	"k3s_centos8":    {},
// 	"k3s_centos9":    {},
// 	"k3s_micro_os":   {},
// 	"k3s_sle_micro":  {},
// 	"k3s_coreos":     {},
// }

//
// // https://github.com/k3s-io/k3s/blob/master/install.sh
// // https://github.com/rancher/rke2/blob/master/install.sh
// // Based on this info, this is the way to validate the correct context
// func getasaContext(product string, ip string) map[string]string {
// 	rke2_centos7 := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /lib/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/rke2*":                            "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/bin/rke2":                                                  "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/bin/rke2":                                            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/cni -I .. -I .":                                        "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/cni/* -I .. -I .":                                      "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /opt/cni -I .. -I .":                                            "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /opt/cni/* -I . -I ..":                                          "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/kubelet/pods -I .. -I .":                               "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/kubelet/pods/* -I .. -I .":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2 -I .. -I .":                               "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/* -I .. -I .":                             "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data(/.*)?":                               "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/*/charts -I .. -I .":                 "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/*/charts/*":                          "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/snapshots -I .. -I .":  "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/snapshots/ -I .. -I .": "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/snapshots/[^/]*/.*":    "",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/sandboxes -I .. -I .":  "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/sandboxes/ -I .. -I .": "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs -I .. -I .":                   "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs/ -I .. -I .":                  "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/run/flannel -I .. -I .":                                    "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/flannel/* -I .. -I .":                                  "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s -I .. -I .":                                        "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/* -I .. -I .":                                      "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/*/sandboxes/*/shm":                      "system_u:object_r:container_runtime_tmpfs_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/*/sandboxes/*/shm/*":                    "system_u:object_r:container_runtime_tmpfs_t",
// 		"sudo ls -laZ /var/log/containers -I .. -I .":                                 "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/log/containers/* -I .. -I .":                               "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/log/pods -I .. -I .":                                       "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/log/pods/* -I .. -I .":                                     "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls -I .. -I .":                    "system_u:object_r:rke2_tls_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls/* -I .. -I .":                  "system_u:object_r:rke2_tls_t",
// 	}
// 	rke2_centos8 := map[string]string{
// 		// https://github.com/rancher/rke2-selinux/blob/master/policy/centos8/rke2.fc
// 		"sudo ls -laZ /etc/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /lib/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/rke2*":                            "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/bin/rke2":                                                  "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/bin/rke2":                                            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /opt/cni -I .. -I .":                                            "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /opt/cni/* -I .. -I .":                                          "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2(/.*)?":                                    "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data(/.*)?":                               "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/[^/]/charts(/.*)?":                   "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots":         "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]*":   "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/sandboxes(/.*)?":   "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs(/.*)?":                        "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls(/.*)?":                         "system_u:object_r:rke2_tls_t",
// 	}
// 	rke2_centos9 := map[string]string{
// 		// https://github.com/rancher/rke2-selinux/blob/master/policy/centos9/rke2.fc
// 		"sudo ls -laZ /etc/systemd/system/rke2*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /lib/systemd/system/rke2*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/rke2*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/bin/rke2":                                                   "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/bin/rke2":                                             "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /opt/cni -I .. -I .":                                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /opt/cni/* -I . -I ..":                                           "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2 -I .. -I .":                                "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/* -I .. -I .":                              "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data -I .. -I .":                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/* -I .. -I .":                         "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/*/charts -I .. -I .":                  "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/*/charts/*":                           "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/snapshots -I . -I ..":   "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/snapshots/ -I . -I ..":  "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]*/.*": "",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/sandboxes -I .. -I .":   "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/*/sandboxes/ -I .. -I .":  "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs -I .. -I .":                    "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls -I .. -I .":                     "system_u:object_r:rke2_tls_t",
// 	}
// 	rke2_micro_os := map[string]string{
// 		// https://github.com/rancher/rke2-selinux/blob/master/policy/microos/rke2.fc
// 		"sudo ls -laZ /etc/systemd/system/rke2*":                                      "	system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /lib/systemd/system/rke2*":                                      "	system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/rke2.*":                           "	system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/bin/rke2":                                                  "	system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/bin/rke2":                                            "	system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /opt/cni(/.*)?":                                                 "	system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2(/.*)?":                                    "	system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data(/.*)?":                               "	system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/[^/]/charts(/.*)?":                   "	system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots":         "	system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]":    "	system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]/.*": "	<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/sandboxes(/.)?":    "	system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs(/.*)?":                        "	system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls(/.*)?":                         "	system_u:object_r:rke2_tls_t",
// 	}
// 	rke2_sle_micro := map[string]string{
// 		// https://github.com/rancher/rke2-selinux/blob/master/policy/slemicro/rke2.fc
// 		"sudo ls -laZ /etc/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /lib/systemd/system/rke2*":                                      "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/rke2.*":                           "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/bin/rke2	":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/bin/rke2":                                            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /opt/rke2/bin/rke2":                                             "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /opt/cni(/.*)?":                                                 "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2(/.*)?":                                    "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data(/.*)?":                               "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/data/[^/]/charts(/.*)?":                   "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots":         "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/snapshots/[^/]/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/rke2/agent/containerd/[^/]*/sandboxes(/.)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/logs(/.*)?":                        "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/lib/rancher/rke2/server/tls(/.*)?":                         "system_u:object_r:rke2_tls_t",
// 	}
// 	k3s_centos7 := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/s?bin/k3s":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/s?bin/k3s":                                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/cni(/.*)?":                                             "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/kubelet/pods(/.*)?":                                    "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/log/containers(/.*)?":                                      "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/log/pods(/.*)?":                                            "system_u:object_r:container_log_t",
// 		"sudo ls -laZ /var/run/flannel(/.*)?":                                         "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
// 	k3s_centos8 := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/bin/k3s":                                             "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
// 	k3s_centos9 := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/s?bin/k3s":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/s?bin/k3s":                                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
// 	k3s_coreos := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/s?bin/k3s":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/s?bin/k3s":                                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
// 	k3s_micro_os := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/s?bin/k3s":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/s?bin/k3s":                                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
// 	k3s_sle_micro := map[string]string{
// 		"sudo ls -laZ /etc/systemd/system/k3s*":                                       "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/lib/systemd/system/k3s*":                                   "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/local/lib/systemd/system/k3s*":                             "system_u:object_r:container_unit_file_t",
// 		"sudo ls -laZ /usr/s?bin/k3s":                                                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /usr/local/s?bin/k3s":                                           "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s(/.*)?":                                     "system_u:object_r:container_var_lib_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots":          "-d system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*":    "-d system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/snapshots/[^/]*/.*": "<<none>>",
// 		"sudo ls -laZ /var/lib/rancher/k3s/agent/containerd/[^/]*/sandboxes(/.*)?":    "system_u:object_r:container_share_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data(/.*)?":                                "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/.lock":                                "system_u:object_r:k3s_lock_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin(/.*)?":                      "system_u:object_r:k3s_root_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]links":                   "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/[.]sha256sums":              "system_u:object_r:k3s_data_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/cni":                        "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd":                 "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim":            "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/containerd-shim-runc-v[12]": "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/bin/runc":                       "system_u:object_r:container_runtime_exec_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/data/[^/]*/etc(/.*)?":                      "system_u:object_r:container_config_t",
// 		"sudo ls -laZ /var/lib/rancher/k3s/storage(/.*)?":                             "system_u:object_r:container_file_t",
// 		"sudo ls -laZ /var/run/k3s(/.*)?":                                             "system_u:object_r:container_var_run_t",
// 		"sudo ls -laZ /var/run/k3s/containerd/[^/]*/sandboxes/[^/]*/shm(/.*)?":        "system_u:object_r:container_runtime_tmpfs_t",
// 	}
//
// 	res, err := shared.RunCommandOnNode("cat /etc/os-release", ip)
// 	Expect(err).NotTo(HaveOccurred())
//
// 	if strings.Contains(res, "ID_LIKE='suse'") {
// 		if strings.Contains(res, "VARIANT_ID='sle-micro'") {
// 			if product == "k3s" {
// 				fmt.Println("Using 'slemicro' policy for this K3S cluster.")
// 				return k3s_sle_micro
// 			} else {
// 				fmt.Println("Using 'slemicro' policy for this RKE2 cluster.")
// 				return rke2_sle_micro
// 			}
// 		}
// 		if product == "k3s" {
// 			fmt.Println("Using 'microos' policy for this K3S cluster")
// 			return k3s_micro_os
// 		} else {
// 			fmt.Println("Using 'microos' policy for this RKE2 cluster.")
// 			return rke2_micro_os
// 		}
// 	}
// 	if strings.Contains(res, "ID_LIKE='coreos'") || strings.Contains(res, "VARIANT_ID='coreos'") {
// 		fmt.Println("Using 'coreos' policy for this k3s cluster")
// 		return k3s_coreos
// 	}
// 	if strings.Contains(res, "VERSION_ID") {
// 		res, err := shared.RunCommandOnNode("cat /etc/os-release | grep 'VERSION_ID'", ip)
// 		Expect(err).NotTo(HaveOccurred())
//
// 		parts := strings.Split(res, "=")
//
// 		if len(parts) == 2 {
// 			version := strings.Trim(parts[1], "\"")
// 			if strings.HasPrefix(version, "7") {
// 				if product == "k3s" {
// 					fmt.Println("Using 'centos7' policy for this K3S cluster")
// 					return k3s_centos7
// 				} else {
// 					fmt.Println("Using 'centos7' policy for this RKE2 cluster.")
// 					return rke2_centos7
// 				}
// 			}
// 			if strings.HasPrefix(version, "8") {
// 				if product == "k3s" {
// 					fmt.Println("Using 'centos8' policy for this K3S cluster")
// 					return k3s_centos8
// 				} else {
// 					fmt.Println("Using 'centos8' policy for this RKE2 cluster")
// 					return rke2_centos8
// 				}
// 			}
// 			if strings.HasPrefix(version, "9") {
// 				if product == "k3s" {
// 					fmt.Println("Using 'centos9' policy for this K3S cluster")
// 					return k3s_centos9
// 				} else {
// 					fmt.Println("Using 'centos9' policy for this RKE2 cluster")
// 					return rke2_centos9
// 				}
// 			}
// 		}
// 	}
//
// 	return rke2_micro_os
// }
