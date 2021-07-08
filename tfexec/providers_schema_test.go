package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersSchemaCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		schemaCmd, err := tf.providersSchemaCmd(context.Background())

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"providers",
			"schema",
			"-json",
			"-no-color",
		}, nil, schemaCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		schemaCmd, err := tf.providersSchemaCmd(context.Background(),
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"providers",
			"schema",
			"-json",
			"-no-color",
		}, nil, schemaCmd)
	})
}
