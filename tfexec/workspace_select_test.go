package tfexec

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceSelectCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{
		// propagate for temp dirs
		"TMPDIR":      os.Getenv("TMPDIR"),
		"TMP":         os.Getenv("TMP"),
		"TEMP":        os.Getenv("TEMP"),
		"USERPROFILE": os.Getenv("USERPROFILE"),
	})

	t.Run("defaults", func(t *testing.T) {
		workspaceSelectCmd, err := tf.workspaceSelectCmd(context.Background(), "testworkspace")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "select",
			"-no-color",
			"testworkspace",
		}, nil, workspaceSelectCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		workspaceSelectCmd, err := tf.workspaceSelectCmd(context.Background(), "testworkspace", Chdir("testpath"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"workspace", "select",
			"-no-color",
			"testworkspace",
		}, nil, workspaceSelectCmd)
	})
}
