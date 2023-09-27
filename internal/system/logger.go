package system

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DefaultLogger struct {
	Context context.Context
}

func (o DefaultLogger) Output(val string) {
	tflog.Info(o.Context, val)
}

func NewDefaultLogger(context context.Context) DefaultLogger {
	return DefaultLogger{
		Context: context,
	}
}
