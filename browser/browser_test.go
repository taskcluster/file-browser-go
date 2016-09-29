package browser;

import (
	"testing";
	"os";
	"path/filepath";
	"container/list";
)

// Utility Functions
type Pair struct {
	First, Second string;
}
func CompareDirectory(root1, root2 string) bool {
	q := list.New();
	q.PushBack(Pair{root1, root2});
	m, n := len(root1), len(root2);
	for q.Len() > 0 {
		temp := q.Front().Value.(Pair);
		q.Remove(q.Front());
		p1, p2 := temp.First, temp.Second;
		if p1[m:] != p2[n:] {
			return false;
		}
		if IsDir(p1) != IsDir(p2) {
			return false;
		}
		if IsDir(p1) {
			f1, err1 := os.Open(p1);
			f2, err2 := os.Open(p2);
			if err1 != nil || err2 != nil {
				return false;
			}
			names1, err1 := f1.Readdirnames(-1);
			names2, err2 := f2.Readdirnames(-1);
			if len(names1) != len(names2) {
				return false;
			}
			for i, _ := range names1 {
				q.PushBack(Pair{ filepath.Join(p1, names1[i]), filepath.Join(p2, names2[i]) });
			}
		}
	}
	return true;
}

func TestList(t *testing.T) {
	res := List("/home");
	if len(res.(*ResultSet).GetDirs()) == 0 {
		t.Fail();
	}
	t.Log(res);
	res = List("Directory/Not/Valid");
	if len(res.(*ResultSet).GetDirs()) != 0 {
		t.Fail();
	}
}

func TestMakeDirectoryAndRemove(t *testing.T) {
	home := os.Getenv("HOME");
	paths := []string{"/test_folder","/test_folder/sub_folder"};

	defer func() {
		r := Remove(filepath.Join(home, paths[0])).(*ResultSet);
		if r.Err != "" {
			t.Log(r.Err);
			t.Fail();
		}
	}();

	for _, p := range paths {
		_ = MakeDirectory(filepath.Join(home, p));
	}
	res := List(paths[0]).(*ResultSet);
	var flag bool = false;
	res = List(home).(*ResultSet);
	t.Log(res);
	for _, p := range res.GetDirs() {
		if p == "test_folder" {
			flag = true;
			break;
		}
	}
	if !flag {
		t.Logf("Folder /test_folder not created.");
		t.Fail();
	}
	res = List(filepath.Join(home,"/test_folder")).(*ResultSet);
	if len(res.GetDirs()) == 0 {
		t.Fail();
	}
}

func TestCopy (t *testing.T) {
	// Create a directory for copying
	home := os.Getenv("HOME");
	dir := []string{"copy_folder", "copy_folder/sub1",
	"copy_folder/sub2", "copy_folder/sub1/sub3", 
	"copy_to"};
	for _, p := range dir {
		res := MakeDirectory(filepath.Join(home,p));
		if err := res.(*ResultSet).Err; err != "" {
			t.Log(err);
			t.FailNow();
		}
	}

	d1 := filepath.Join(home,"copy_folder");
	d2 := filepath.Join(home,"copy_to/copy_folder");

	defer func() {
		_ = Remove(d1);
		_ = Remove(filepath.Join(home, "copy_to"));
	}();

	_ = Copy(d1, d2);		
	WaitForOperationsToComplete();
	if CompareDirectory(d1,d2) == false {
		t.Fail();
	}
}

func TestRun(t *testing.T) {
	cmds := []Command{
		{ Cmd:"List", Args:[]string{"/home"}, },
		{ Cmd:"List", Args:[]string{"/var/"}, },
		{ Cmd:"Exit", Args:[]string{}, },
	}

	for _, cmd := range cmds {
		ret := Run(cmd);
		if ret == nil {
			break;
		}
		res := ret.(*ResultSet);
		if res.Err != "" {
			t.Logf(res.Err);
			t.Fail();
		}
		t.Log(res);
	}
}
