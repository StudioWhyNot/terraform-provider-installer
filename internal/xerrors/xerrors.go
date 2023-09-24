package xerrors

import (
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var ErrDoubleVersions = errors.New("version cannot be specified both in the name and explicitly")
var ErrVersionNotFound = errors.New("version not found")
var ErrNotInstalled = errors.New("not installed")

func ToDiags(err error) diag.Diagnostics {
	diags := diag.Diagnostics{}
	diags.AddError(err.Error(), errors.FlattenDetails(err))
	return diags
}

func AppendToDiagnostics(diags *diag.Diagnostics, err error) {
	diags.Append(ToDiags(err)...)
}
