package tfexec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type showConfig struct {
	chdir        string
	planArg      string
	reattachInfo ReattachInfo
	stateArg     string
}

var defaultShowOptions = showConfig{}

type ShowOption interface {
	configureShow(*showConfig)
}

func (opt *ChdirOption) configureShow(conf *showConfig) {
	conf.chdir = opt.path
}

func (opt *PlanArgOption) configureShow(conf *showConfig) {
	conf.planArg = opt.path
}

func (opt *StateArgOption) configureShow(conf *showConfig) {
	conf.stateArg = opt.path
}

func (opt *ReattachOption) configureShow(conf *showConfig) {
	conf.reattachInfo = opt.info
}

// Show reads the default state path and outputs the state.
// To read a state or plan file, ShowState or ShowPlan must be used instead.
func (tf *Terraform) Show(ctx context.Context, opts ...ShowOption) (*tfjson.State, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd, err := tf.showCmd(ctx, true, mergeEnv)

	if err != nil {
		return nil, err
	}

	var ret tfjson.State
	ret.UseJSONNumber(true)
	err = tf.runTerraformCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowStateFile reads a given state file and outputs the state.
func (tf *Terraform) ShowStateFile(ctx context.Context, statePath string, opts ...ShowOption) (*tfjson.State, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	if statePath == "" {
		return nil, fmt.Errorf("statePath cannot be blank: use Show() if not passing statePath")
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd, err := tf.showCmd(ctx, true, mergeEnv, append(opts, StateArg(statePath))...)

	if err != nil {
		return nil, err
	}

	var ret tfjson.State
	ret.UseJSONNumber(true)
	err = tf.runTerraformCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowPlanFile reads a given plan file and outputs the plan.
func (tf *Terraform) ShowPlanFile(ctx context.Context, planPath string, opts ...ShowOption) (*tfjson.Plan, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	if planPath == "" {
		return nil, fmt.Errorf("planPath cannot be blank: use Show() if not passing planPath")
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd, err := tf.showCmd(ctx, true, mergeEnv, append(opts, PlanArg(planPath))...)

	if err != nil {
		return nil, err
	}

	var ret tfjson.Plan
	err = tf.runTerraformCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

// ShowPlanFileRaw reads a given plan file and outputs the plan in a
// human-friendly, opaque format.
func (tf *Terraform) ShowPlanFileRaw(ctx context.Context, planPath string, opts ...ShowOption) (string, error) {
	if planPath == "" {
		return "", fmt.Errorf("planPath cannot be blank: use Show() if not passing planPath")
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return "", err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd, err := tf.showCmd(ctx, false, mergeEnv, append(opts, PlanArg(planPath))...)

	if err != nil {
		return "", err
	}

	var ret bytes.Buffer
	showCmd.Stdout = &ret
	err = tf.runTerraformCmd(ctx, showCmd)
	if err != nil {
		return "", err
	}

	return ret.String(), nil

}

func (tf *Terraform) showCmd(ctx context.Context, jsonOutput bool, mergeEnv map[string]string, opts ...ShowOption) (*exec.Cmd, error) {
	c := defaultShowOptions

	for _, o := range opts {
		switch o.(type) {
		case *ChdirOption:
			err := tf.compatible(ctx, tf0_14_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-chdir was added in Terraform 0.14: %w", err)
			}
		}

		o.configureShow(&c)
	}

	var args []string

	// global opts
	if c.chdir != "" {
		args = append(args, "-chdir="+c.chdir)
	}

	args = append(args, "show")

	if jsonOutput {
		args = append(args, "-json")
	}

	args = append(args, "-no-color")

	if c.planArg != "" {
		args = append(args, c.planArg)
	}

	if c.stateArg != "" {
		args = append(args, c.stateArg)
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
