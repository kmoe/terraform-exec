package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestOutputCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		outputCmd, err := tf.outputCmd(context.Background())

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
		}, nil, outputCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		outputCmd, err := tf.outputCmd(context.Background(),
			State("teststate"))

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
			"-state=teststate",
		}, nil, outputCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		outputCmd, err := tf.outputCmd(context.Background(),
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"output",
			"-no-color",
			"-json",
		}, nil, outputCmd)
	})
}
