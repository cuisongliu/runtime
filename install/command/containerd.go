package command

import (
	"bytes"
	"fmt"
	sealos "github.com/fanux/sealos/install"
	"github.com/wonderivan/logger"
	"text/template"
)

type Containerd struct{}

func NewContainerd() stepInterface {
	var stepInterface stepInterface
	stepInterface = &Containerd{}
	return stepInterface
}

const containerdFileName = "containerd.tgz"

func (d *Containerd) lib() string {
	if Lib == "" {
		return "/var/lib/containerd"
	} else {
		return Lib
	}
}
func (d *Containerd) SendPackage(host string) {
	sendPackage(host, PkgUrl, containerdFileName)
}

func (d *Containerd) Tar(host string) {
	cmd := fmt.Sprintf("tar --strip-components=1 -xvzf /root/%s -C /usr/local/bin", containerdFileName)
	sealos.Cmd(host, cmd)
}
func (d *Containerd) Config(host string) {
	cmd := "mkdir -p " + d.lib()
	sealos.Cmd(host, cmd)
	cmd = "mkdir -p /etc/containerd"
	sealos.Cmd(host, cmd)
	//cmd = "containerd config default > /etc/containerd/config.toml"
	cmd = "echo \"" + string(d.configFile()) + "\" > /etc/containerd/config.toml"
	sealos.Cmd(host, cmd)
}

func (d *Containerd) Enable(host string) {
	cmd := "echo \"" + string(d.serviceFile()) + "\" > /usr/lib/systemd/system/containerd.service"
	sealos.Cmd(host, cmd)
	cmd = "systemctl enable  containerd.service && systemctl restart  containerd.service"
	sealos.Cmd(host, cmd)
}

func (d *Containerd) Version(host string) {
	cmd := "containerd --version"
	sealos.Cmd(host, cmd)
	logger.Warn("pull docker hub command. ex: ctr images pull docker.io/library/alpine:3.8")
	logger.Warn("pull http registry command. ex:  ctr images pull 10.0.45.222/library/alpine:3.8 --plain-http")
}

func (d *Containerd) Uninstall(host string) {
	cmd := "systemctl stop  containerd.service && systemctl disable containerd.service"
	sealos.Cmd(host, cmd)
	cmd = "rm -rf /usr/local/bin/ctr && rm -rf /usr/local/bin/containerd* "
	sealos.Cmd(host, cmd)
	cmd = "rm -rf /var/lib/containerd && rm -rf /etc/containerd/* "
	sealos.Cmd(host, cmd)
	if d.lib() != "" {
		cmd = "rm -rf " + d.lib()
		sealos.Cmd(host, cmd)
	}
}

func (d *Containerd) serviceFile() []byte {
	var templateText = string(`[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target
  
[Service]
ExecStart=/usr/local/bin/containerd
Restart=always
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
  
[Install]
WantedBy=multi-user.target
`)
	return []byte(templateText)
}
func (d *Containerd) configFile() []byte {
	var templateText = string(`root = \"{{.CONTAINERD_LIB}}\"
state = \"/run/containerd\"
oom_score = 0

[grpc]
  address = \"/run/containerd/containerd.sock\"
  uid = 0
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216

[debug]
  address = \"\"
  uid = 0
  gid = 0
  level = \"\"

[metrics]
  address = \"\"
  grpc_histogram = false

[cgroup]
  path = \"\"

[plugins]
  [plugins.cgroups]
    no_prometheus = false
  [plugins.cri]
    stream_server_address = \"127.0.0.1\"
    stream_server_port = \"0\"
    enable_selinux = false
    sandbox_image = \"k8s.gcr.io/pause:3.1\"
    stats_collect_period = 10
    systemd_cgroup = false
    enable_tls_streaming = false
    max_container_log_line_size = 16384
    [plugins.cri.containerd]
      snapshotter = \"overlayfs\"
      no_pivot = false
      [plugins.cri.containerd.default_runtime]
        runtime_type = \"io.containerd.runtime.v1.linux\"
        runtime_engine = \"\"
        runtime_root = \"\"
      [plugins.cri.containerd.untrusted_workload_runtime]
        runtime_type = \"\"
        runtime_engine = \"\"
        runtime_root = \"\"
    [plugins.cri.cni]
      bin_dir = \"/opt/cni/bin\"
      conf_dir = \"/etc/cni/net.d\"
      conf_template = \"\"
    [plugins.cri.registry]
      [plugins.cri.registry.mirrors]
        [plugins.cri.registry.mirrors.\"docker.io\"]
          endpoint = [\"https://registry-1.docker.io\"]
        {{range .CONTAINERD_REGISTRY -}}[plugins.cri.registry.mirrors.\"{{.}}\"]
          endpoint = [\"{{.}}\"]
    {{end -}}
    [plugins.cri.x509_key_pair_streaming]
      tls_cert_file = \"\"
      tls_key_file = \"\"
  [plugins.diff-service]
    default = [\"walking\"]
  [plugins.linux]
    shim = \"containerd-shim\"
    runtime = \"runc\"
    runtime_root = \"\"
    no_shim = false
    shim_debug = false
  [plugins.opt]
    path = \"/opt/containerd\"
  [plugins.restart]
    interval = \"10s\"
  [plugins.scheduler]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = \"0s\"
    startup_delay = \"100ms\"
`)
	tmpl, err := template.New("text").Parse(templateText)
	if err != nil {
		logger.Error("template parse failed:", err)
		panic(1)
	}
	var envMap = make(map[string]interface{})
	envMap["CONTAINERD_REGISTRY"] = RegistryArr
	envMap["CONTAINERD_LIB"] = d.lib()
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}

func (d *Containerd) Print() {
	urlPrefix := "https://github.com/containerd/containerd/releases/download/v%s/containerd-%s.linux-amd64.tar.gz"
	versions := []string{
		"1.1.0",
		"1.1.1",
		"1.1.2",
		"1.1.3",
		"1.1.4",
		"1.1.5",
		"1.1.6",
		"1.1.7",

		"1.2.0",
		"1.2.1",
		"1.2.2",
		"1.2.3",
		"1.2.4",
		"1.2.5",
		"1.2.6",
		"1.2.7",
	}

	for _, v := range versions {
		println(fmt.Sprintf(urlPrefix, v, v))
	}
}