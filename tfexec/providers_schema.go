package tfexec

import (
	"context"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type providersSchemaConfig struct {
	chdir string
}

var defaultProvidersSchemaOptions = providersSchemaConfig{}

type ProvidersSchemaOption interface {
	configureProvidersSchema(*providersSchemaConfig)
}

func (opt *ChdirOption) configureProvidersSchema(conf *providersSchemaConfig) {
	conf.chdir = opt.path
}

// ProvidersSchema represents the terraform providers schema -json subcommand.
func (tf *Terraform) ProvidersSchema(ctx context.Context) (*tfjson.ProviderSchemas, error) {
	schemaCmd, err := tf.providersSchemaCmd(ctx)

	if err != nil {
		return nil, err
	}

	var ret tfjson.ProviderSchemas
	err = tf.runTerraformCmdJSON(ctx, schemaCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Terraform) providersSchemaCmd(ctx context.Context, opts ...ProvidersSchemaOption) (*exec.Cmd, error) {
	c := defaultProvidersSchemaOptions

	for _, o := range opts {
		switch o.(type) {
		case *ChdirOption:
			err := tf.compatible(ctx, tf0_14_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-chdir was added in Terraform 0.14: %w", err)
			}
		}

		o.configureProvidersSchema(&c)
	}

	var args []string

	// global opts
	if c.chdir != "" {
		args = append(args, "-chdir="+c.chdir)
	}

	args = append(args, []string{"providers", "schema", "-json", "-no-color"}...)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
