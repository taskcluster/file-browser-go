package browser;

import(
	"os";
	"container/list";
	"path/filepath";
	"testing";
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

func Exists (path string) bool {
	_, err := os.Stat(path);
	if err == nil {
		return true;
	}
	return os.IsExist(err);
}

func FailNotNil(err error, t *testing.T){
	if err != nil{
		t.Log(err);
		t.FailNow();
	}
}
