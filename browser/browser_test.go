package browser;

import (
	"testing";
)

func TestList(t *testing.T) {
	res := List("/home");
	if len(res.GetDirs()) == 0 {
		t.Fail();
	}
	for _, f := range res.GetDirs() {
		t.Logf(f);
	}
	res = List("Directory/Not/Valid");
	if len(res.GetDirs()) != 0 {
		t.Fail();
	}
}

func TestRun(t *testing.T) {
	cmds := []Command{
		{ Cmd:"ls", Args:[]string{"/home"}, },
		{ Cmd:"ls", Args:[]string{"/var/"}, },
		{ Cmd:"exit", Args:[]string{}, },
	}

	for _, cmd := range cmds {
		res := Run(cmd);
		t.Log(res);
	}
}
