package terraformutils

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/shared"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/ssh"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/winrm"
)

type RemoteConnectionInfo struct {
	Type           types.String `tfsdk:"type"`
	User           types.String `tfsdk:"user"`
	Password       types.String `tfsdk:"password"`
	PrivateKey     types.String `tfsdk:"private_key"`
	Certificate    types.String `tfsdk:"certificate"`
	Host           types.String `tfsdk:"host"`
	HostKey        types.String `tfsdk:"host_key"`
	Port           types.Int64  `tfsdk:"port"`
	Agent          types.Bool   `tfsdk:"agent"`
	ScriptPath     types.String `tfsdk:"script_path"`
	TargetPlatform types.String `tfsdk:"target_platform"`
	Timeout        types.String `tfsdk:"timeout"`

	ProxyScheme       types.String `tfsdk:"proxy_scheme"`
	ProxyHost         types.String `tfsdk:"proxy_host"`
	ProxyPort         types.Int64  `tfsdk:"proxy_port"`
	ProxyUserName     types.String `tfsdk:"proxy_user_name"`
	ProxyUserPassword types.String `tfsdk:"proxy_user_password"`

	BastionUser        types.String `tfsdk:"bastion_user"`
	BastionPassword    types.String `tfsdk:"bastion_password"`
	BastionPrivateKey  types.String `tfsdk:"bastion_private_key"`
	BastionCertificate types.String `tfsdk:"bastion_certificate"`
	BastionHost        types.String `tfsdk:"bastion_host"`
	BastionHostKey     types.String `tfsdk:"bastion_host_key"`
	BastionPort        types.Int64  `tfsdk:"bastion_port"`

	AgentIdentity types.String `tfsdk:"agent_identity"`

	HTTPS    types.Bool   `tfsdk:"https"`
	Insecure types.Bool   `tfsdk:"insecure"`
	NTLM     types.Bool   `tfsdk:"use_ntlm"`
	CACert   types.String `tfsdk:"cacert"`
}

const userConnectionSeperator = "@"

func (r *RemoteConnectionInfo) GetConnectionName() string {
	if r == nil || r.User.IsNull() || r.Host.IsNull() {
		return ""
	}
	return r.User.ValueString() + userConnectionSeperator + r.Host.ValueString()
}

func (r *RemoteConnectionInfo) WaitForHost() error {
	if r == nil || r.User.IsNull() || r.Host.IsNull() {
		return nil
	}
	port, timeout, err := r.GetPortAndTimeout()
	if err != nil {
		return err
	}
	host := shared.IpFormat(r.Host.ValueString())
	host = fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (r *RemoteConnectionInfo) GetPortAndTimeout() (int64, time.Duration, error) {
	var err error
	port, timeout := r.GetDefaultPortAndTimeout()
	if !r.Port.IsNull() {
		port = r.Port.ValueInt64()
	}
	if !r.Timeout.IsNull() {
		timeout, err = time.ParseDuration(r.Timeout.ValueString())
	}
	return port, timeout, err
}

func (r *RemoteConnectionInfo) GetDefaultPortAndTimeout() (int64, time.Duration) {
	switch r.Type.ValueString() {
	case "winrm":
		return winrm.DefaultPort, winrm.DefaultTimeout
	case "ssh":
		fallthrough
	default:
		return ssh.DefaultPort, ssh.DefaultTimeout
	}
}

func MakeCommunicator(info *RemoteConnectionInfo) (communicator.Communicator, error) {
	if info == nil || info.Host.ValueString() == "" {
		return nil, nil
	}
	valMap := StructToCtyValueMap(info)
	return communicator.New(*valMap)
}
