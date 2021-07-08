package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestRefreshCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
		}, nil, refreshCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("refreshdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"refreshdir",
		}, nil, refreshCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		refreshCmd, err := tf.refreshCmd(context.Background(),
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"refresh",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
		}, nil, refreshCmd)
	})
}
