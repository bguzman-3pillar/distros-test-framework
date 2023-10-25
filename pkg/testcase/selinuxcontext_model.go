package testcase

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

var conf = []configuration{
	{
		// Works correctly!
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
		// We are not able to execute this because our framework does not support the reboot part for this OS.
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
		// We are not able to execute this because our framework does not support the reboot part for this OS.
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
		// We are not able to execute this because our framework does not support the reboot part for this OS.
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
		// We are not able to execute this because our framework does not support the reboot part for this OS.
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
		// We are not able to execute this because our framework does not support the reboot part for this OS.
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
