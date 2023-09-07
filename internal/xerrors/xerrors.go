package xerrors

import (
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var ErrVersionNotFound = errors.New("version not found")
var ErrNotInstalled = errors.New("not installed")

func ToDiags(err error) diag.Diagnostics {
	diags := diag.Diagnostics{}
	diags.AddError(err.Error(), errors.FlattenDetails(err))
	return diags
}
