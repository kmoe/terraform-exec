package tfexec

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestFormat(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("chdir", func(t *testing.T) {
		fmtCmd, err := tf.formatCmd(context.Background(),
			[]string{"testfile"},
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"fmt",
			"-no-color",
			"testfile",
		}, nil, fmtCmd)
	})

	t.Run("too old version", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, "0.7.6"))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		_, err = tf.formatCmd(context.Background(), []string{})
		if err == nil {
			t.Fatal("expected old version to fail")
		}

		var expectedErr *ErrVersionMismatch
		if !errors.As(err, &expectedErr) {
			t.Fatalf("error doesn't match: %#v", err)
		}
	})
}
