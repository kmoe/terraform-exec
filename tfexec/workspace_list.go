package tfexec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type workspaceListConfig struct {
	chdir string
}

var defaultWorkspaceListOptions = workspaceListConfig{}

// WorkspaceListCmdOption represents options that are applicable to the WorkspaceList method.
type WorkspaceListCmdOption interface {
	configureWorkspaceList(*workspaceListConfig)
}

func (opt *ChdirOption) configureWorkspaceList(conf *workspaceListConfig) {
	conf.chdir = opt.path
}

// WorkspaceList represents the workspace list subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceList(ctx context.Context, opts ...WorkspaceListCmdOption) ([]string, string, error) {
	// TODO: [DIR] param option
	wlCmd, err := tf.workspaceListCmd(ctx, opts...)

	if err != nil {
		return nil, "", err
	}

	var outBuf bytes.Buffer
	wlCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(ctx, wlCmd)
	if err != nil {
		return nil, "", err
	}

	ws, current := parseWorkspaceList(outBuf.String())

	return ws, current, nil
}

func (tf *Terraform) workspaceListCmd(ctx context.Context, opts ...WorkspaceListCmdOption) (*exec.Cmd, error) {
	// TODO: [DIR] param option

	c := defaultWorkspaceListOptions

	for _, o := range opts {
		switch o.(type) {
		case *ChdirOption:
			err := tf.compatible(ctx, tf0_14_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-chdir was added in Terraform 0.14: %w", err)
			}
		}

		o.configureWorkspaceList(&c)
	}

	var args []string

	// global opts
	if c.chdir != "" {
		args = append(args, "-chdir="+c.chdir)
	}

	args = append(args, []string{"workspace", "list", "-no-color"}...)

	cmd := tf.buildTerraformCmd(ctx, nil, args...)

	return cmd, nil
}

const currentWorkspacePrefix = "* "

func parseWorkspaceList(stdout string) ([]string, string) {
	lines := strings.Split(stdout, "\n")

	current := ""
	workspaces := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, currentWorkspacePrefix) {
			line = strings.TrimPrefix(line, currentWorkspacePrefix)
			current = line
		}
		workspaces = append(workspaces, line)
	}

	return workspaces, current
}
