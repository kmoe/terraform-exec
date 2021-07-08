package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestImportCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		importCmd, err := tf.importCmd(context.Background(), "my-addr", "my-id")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"import",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"my-addr",
			"my-id",
		}, nil, importCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		importCmd, err := tf.importCmd(context.Background(), "my-addr2", "my-id2",
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("testvarfile"),
			Lock(false),
			Var("var1=foo"),
			Var("var2=bar"),
			AllowMissingConfig(true),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"import",
			"-no-color",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-allow-missing-config",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"my-addr2",
			"my-id2",
		}, nil, importCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		importCmd, err := tf.importCmd(context.Background(),
			"my-addr",
			"my-id",
			Chdir("testpath"),
		)

		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"import",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"my-addr",
			"my-id",
		}, nil, importCmd)
	})
}
