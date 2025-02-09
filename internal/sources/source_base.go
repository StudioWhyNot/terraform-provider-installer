package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

type SourceData interface {
	Initialize(ctx context.Context) bool
	CopyFromTypedInstalledProgramInfo(programInfo *models.TypedInstalledProgramInfo)
	GetRemoteConnectionInfo() *terraformutils.RemoteConnectionInfo
}

type SourceBase[T any] struct {
	Installer      installers.Installer[T]
	Communicator   communicator.Communicator
	ConnectionInfo *terraformutils.RemoteConnectionInfo
}

func NewSourceBase[T any](installer installers.Installer[T]) *SourceBase[T] {
	return &SourceBase[T]{
		Installer:    installer,
		Communicator: nil,
	}
}

func (s *SourceBase[T]) GetCommunicator() communicator.Communicator {
	return s.Communicator
}

func (s *SourceBase[T]) TryConnect(context context.Context) error {
	if s.Communicator == nil {
		return nil
	}
	err := s.ConnectionInfo.WaitForHost()
	if err != nil {
		return err
	}
	return s.Communicator.Connect(system.NewDefaultLogger(context))
}

func (s *SourceBase[T]) TryDisconnect() error {
	if s.Communicator == nil {
		return nil
	}
	return s.Communicator.Disconnect()
}

func GetIDFromName(name string, installerType enums.InstallerType) types.String {
	return types.StringValue(installerType.GetIDFromName(name))
}

func GetIDFromNameAndVersion(seperator string, name types.String, version types.String, installerType enums.InstallerType) types.String {
	return types.StringValue(models.GetIDFromNameAndVersion(seperator, name.ValueString(), version.ValueString(), installerType))
}

type TerraformDataProvider interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}
