package e2etest

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	validateMinVersion = version.Must(version.NewVersion("0.12.0"))
)

func TestValidate(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(validateMinVersion) {
			t.Skip("terraform validate -json was added in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		validation, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if !validation.Valid {
			t.Fatalf("expected valid, got %#v", validation)
		}
	})

	runTest(t, "invalid", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(validateMinVersion) {
			t.Skip("terraform validate -json was added in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Logf("error initializing: %s", err)

			// allow for invalid config errors only here
			// 0.13 will return this, 0.12 will not
			// unsure why 0.12 terraform init does not have a non-zero exit code for syntax problems
			var confErr *tfexec.ErrConfigInvalid
			if !errors.As(err, &confErr) {
				t.Fatalf("expected err ErrConfigInvalid, got %T: %s", err, err)
			}
		}

		actual, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// reset byte locations in actual as CRLF issues render them off between operating systems
		cleanActual := []tfexec.Diagnostic{}
		for _, diag := range actual.Diagnostics {
			diag.Range.Start.Byte = 0
			diag.Range.End.Byte = 0
			cleanActual = append(cleanActual, diag)
		}

		assert.Equal(t, []tfexec.Diagnostic{
			{
				Severity: "error",
				Summary:  "Unsupported block type",
				Detail:   "Blocks of type \"bad_block\" are not expected here.",
				Range: &tfexec.Range{
					Filename: "main.tf",
					Start: tfexec.Pos{
						Line:   1,
						Column: 1,
					},
					End: tfexec.Pos{
						Line:   1,
						Column: 10,
					},
				},
			},
			{
				Severity: "error",
				Summary:  "Unsupported argument",
				Detail:   "An argument named \"bad_attribute\" is not expected here.",
				Range: &tfexec.Range{
					Filename: "main.tf",
					Start: tfexec.Pos{
						Line:   5,
						Column: 5,
					},
					End: tfexec.Pos{
						Line:   5,
						Column: 18,
					},
				},
			},
		}, cleanActual)
	})
}
