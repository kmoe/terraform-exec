package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type workspaceSelectConfig struct {
	chdir string
}

var defaultWorkspaceSelectOptions = workspaceSelectConfig{}

// WorkspaceSelectCmdOption represents options that are applicable to the WorkspaceSelect method.
type WorkspaceSelectCmdOption interface {
	configureWorkspaceSelect(*workspaceSelectConfig)
}

func (opt *ChdirOption) configureWorkspaceSelect(conf *workspaceSelectConfig) {
	conf.chdir = opt.path
}

// WorkspaceSelect represents the workspace select subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceSelect(ctx context.Context, workspace string, opts ...WorkspaceSelectCmdOption) error {
	// TODO: [DIR] param option

	cmd, err := tf.workspaceSelectCmd(ctx, workspace, opts...)

	if err != nil {
		return err
	}

	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) workspaceSelectCmd(ctx context.Context, workspace string, opts ...WorkspaceSelectCmdOption) (*exec.Cmd, error) {
	// TODO: [DIR] param option

	c := defaultWorkspaceSelectOptions

	for _, o := range opts {
		switch o.(type) {
		case *ChdirOption:
			err := tf.compatible(ctx, tf0_14_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-chdir was added in Terraform 0.14: %w", err)
			}
		}

		o.configureWorkspaceSelect(&c)
	}

	var args []string

	// global opts
	if c.chdir != "" {
		args = append(args, "-chdir="+c.chdir)
	}

	args = append(args, []string{"workspace", "select", "-no-color", workspace}...)

	cmd := tf.buildTerraformCmd(ctx, nil, args...)

	return cmd, nil
}
