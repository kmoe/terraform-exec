package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestShowCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	// defaults
	showCmd, err := tf.showCmd(context.Background(), true, nil)

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
	}, nil, showCmd)
}

func TestShowCmdChdir(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	// defaults
	showCmd, err := tf.showCmd(context.Background(), true, nil, Chdir("testpath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"-chdir=testpath",
		"show",
		"-json",
		"-no-color",
	}, nil, showCmd)
}

func TestShowStateFileCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), true, nil, StateArg("statefilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"statefilepath",
	}, nil, showCmd)
}

func TestShowStateFileCmdChdir(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), true, nil, Chdir("testpath"), StateArg("statefilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"-chdir=testpath",
		"show",
		"-json",
		"-no-color",
		"statefilepath",
	}, nil, showCmd)
}

func TestShowPlanFileCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), true, nil, PlanArg("planfilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}

func TestShowPlanFileCmdChdir(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), true, nil, Chdir("testpath"), PlanArg("planfilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"-chdir=testpath",
		"show",
		"-json",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}

func TestShowPlanFileRawCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), false, nil, PlanArg("planfilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"show",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}

func TestShowPlanFileRawCmdChdir(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd, err := tf.showCmd(context.Background(), false, nil, Chdir("testpath"), PlanArg("planfilepath"))

	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"-chdir=testpath",
		"show",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}
