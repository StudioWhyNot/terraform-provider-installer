package terraformutils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
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

func MakeCommunicator(info *RemoteConnectionInfo) (communicator.Communicator, error) {
	valMap := StructToCtyValueMap(info)
	if valMap == nil {
		return nil, nil
	}
	return communicator.New(*valMap)
}
