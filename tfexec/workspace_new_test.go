package tfexec

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceNewCmd(t *testing.T) {
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
		workspaceNewCmd, err := tf.workspaceNewCmd(context.Background(), "workspace-name")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "new",
			"-no-color",
			"workspace-name",
		}, nil, workspaceNewCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		workspaceNewCmd, err := tf.workspaceNewCmd(context.Background(), "workspace-name", LockTimeout("200s"), CopyState("teststate"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "new",
			"-no-color",
			"-lock-timeout=200s",
			"-lock=false",
			"-state=teststate",
			"workspace-name",
		}, nil, workspaceNewCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		workspaceNewCmd, err := tf.workspaceNewCmd(context.Background(), "workspace-name", Chdir("testpath"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"workspace", "new",
			"-no-color",
			"workspace-name",
		}, nil, workspaceNewCmd)
	})
}
