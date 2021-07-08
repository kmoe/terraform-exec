package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestApplyCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("foo.tfvars"),
			VarFile("bar.tfvars"),
			Lock(false),
			Parallelism(99),
			Refresh(false),
			Target("target1"),
			Target("target2"),
			Var("var1=foo"),
			Var("var2=bar"),
			DirOrPlan("testfile"),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=foo.tfvars",
			"-var-file=bar.tfvars",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"testfile",
		}, nil, applyCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
		}, nil, applyCmd)
	})

	t.Run("plan", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			PlanArg("testplan"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"testplan",
		}, nil, applyCmd)
	})
}
