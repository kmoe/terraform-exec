package tfinstall

import (
	"context"

	"github.com/hashicorp/go-version"
)

type ExactVersionOption struct {
	tfVersion  string
	installDir string

	UserAgent string
}

var _ ExecPathFinder = &ExactVersionOption{}

// ExactVersion returns a pointer to an ExactVersionOption object specifying a terraform version and installation directory passed to the function
// Since ExactVersionOption is exported, it can also be instantiated manually
// Manual instantiation of ExactVersionOption will be necessary if a custom UserAgent is required
func ExactVersion(tfVersion string, installDir string) *ExactVersionOption {
	opt := &ExactVersionOption{
		tfVersion:  tfVersion,
		installDir: installDir,
	}

	return opt
}

// ExecPath downloads a given Terraform binary at a given path
// The binary is downloaded from the official terraform release platform (currently https://releases.hashicorp.com/terraform)
// The version of the binary must be specified via the *ExactVersionOption object that this method is called on
// The path where the binary is downloaded must be specified via the *ExactVersionOption object that this method is called on
// It verifies the integrity of the binary via its checksum
// It returns the path where the binary was downloadeds
func (opt *ExactVersionOption) ExecPath(ctx context.Context) (string, error) {
	// validate version
	_, err := version.NewVersion(opt.tfVersion)
	if err != nil {
		return "", err
	}

	return downloadWithVerification(ctx, opt.tfVersion, opt.installDir, opt.UserAgent)
}
