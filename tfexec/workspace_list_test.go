package tfexec

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestParseWorkspaceList(t *testing.T) {
	for i, c := range []struct {
		expected        []string
		expectedCurrent string
		stdout          string
	}{
		{
			[]string{"default"},
			"default",
			`* default

`,
		},
		{
			[]string{"default", "foo", "bar"},
			"foo",
			`  default
* foo
  bar

`,
		},

		// linux new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\n* foo\n\n",
		},
		// windows new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\r\n* foo\r\n\r\n",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actualList, actualCurrent := parseWorkspaceList(c.stdout)

			if actualCurrent != c.expectedCurrent {
				t.Fatalf("expected selected %q, got %q", c.expectedCurrent, actualCurrent)
			}

			if !reflect.DeepEqual(c.expected, actualList) {
				t.Fatalf("expected %#v, got %#v", c.expected, actualList)
			}
		})
	}
}

func TestWorkspaceListCmd(t *testing.T) {
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
		workspaceListCmd, err := tf.workspaceListCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "list",
			"-no-color",
		}, nil, workspaceListCmd)
	})

	t.Run("chdir", func(t *testing.T) {
		workspaceListCmd, err := tf.workspaceListCmd(context.Background(), Chdir("testpath"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"-chdir=testpath",
			"workspace", "list",
			"-no-color",
		}, nil, workspaceListCmd)
	})
}
