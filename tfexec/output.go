package tfexec

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type outputConfig struct {
	chdir string
	state string
	json  bool
}

var defaultOutputOptions = outputConfig{}

// OutputOption represents options used in the Output method.
type OutputOption interface {
	configureOutput(*outputConfig)
}

func (opt *ChdirOption) configureOutput(conf *outputConfig) {
	conf.chdir = opt.path
}

func (opt *StateOption) configureOutput(conf *outputConfig) {
	conf.state = opt.path
}

// OutputMeta represents the JSON output of 'terraform output -json',
// which resembles state format version 3 due to a historical accident.
// Please see hashicorp/terraform/command/output.go.
// TODO KEM: Should this type be in terraform-json?
type OutputMeta struct {
	Sensitive bool            `json:"sensitive"`
	Type      json.RawMessage `json:"type"`
	Value     json.RawMessage `json:"value"`
}

// Output represents the terraform output subcommand.
func (tf *Terraform) Output(ctx context.Context, opts ...OutputOption) (map[string]OutputMeta, error) {
	outputCmd, err := tf.outputCmd(ctx, opts...)

	if err != nil {
		return nil, err
	}

	outputs := map[string]OutputMeta{}
	err = tf.runTerraformCmdJSON(ctx, outputCmd, &outputs)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (tf *Terraform) outputCmd(ctx context.Context, opts ...OutputOption) (*exec.Cmd, error) {
	c := defaultOutputOptions

	for _, o := range opts {
		switch o.(type) {
		case *ChdirOption:
			err := tf.compatible(ctx, tf0_14_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-chdir was added in Terraform 0.14: %w", err)
			}
		}

		o.configureOutput(&c)
	}

	var args []string

	// global opts
	if c.chdir != "" {
		args = append(args, "-chdir="+c.chdir)
	}

	args = append(args, []string{"output", "-no-color", "-json"}...)

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
