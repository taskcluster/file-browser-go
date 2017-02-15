package browser

import (
	// "gopkg.in/vmihailenco/msgpack.v2"
	// "io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

var outChan chan interface{}

func TestMain(m *testing.M) {
	outChan = make(chan interface{})
	os.Exit(m.Run())
}

func TestList(t *testing.T) {
	var res *ResultSet = nil
	temp, err := ioutil.TempDir("", "list")
	FailNotNil(err, t)
	defer os.Remove(temp)
	// Make two directories in temp
	err = os.Mkdir(filepath.Join(temp, "a"), 0777)
	FailNotNil(err, t)
	err = os.Mkdir(filepath.Join(temp, "b"), 0777)
	FailNotNil(err, t)
	go List("test", outChan, temp)
	res = (<-outChan).(*ResultSet)
	if res.Err != "" {
		t.Fail()
	}
	if len(res.Files) != 2 {
		t.Fail()
	}
}

func TestListNotExist(t *testing.T) {
	var res *ResultSet = nil
	go List("test", outChan, "does/not/exist")
	res = (<-outChan).(*ResultSet)
	if res.Err == "" {
		t.Log("Error should occur.")
		t.Fail()
	}
	t.Log("Error: ", res.Err)
}

func TestMakeDirectoryAndRemove(t *testing.T) {
	home, err := ioutil.TempDir("", "MakeAndRemove")
	FailNotNil(err, t)
	paths := []string{"/test_folder", "/test_folder/sub_folder"}
	var res *ResultSet = nil

	defer func() {
		go Remove("test", outChan, filepath.Join(home, paths[0]))
		res = (<-outChan).(*ResultSet)
	}()

	for _, p := range paths {
		go MakeDirectory("test", outChan, filepath.Join(home, p))
		res = (<-outChan).(*ResultSet)
		if res.Err != "" {
			t.Log(res.Err)
			t.Fail()
		}
		if !Exists(filepath.Join(home, p)) {
			t.Log("Directory not created")
			t.Fail()
		}
	}
}

func TestMakeDirectoryBadPath(t *testing.T) {
	var res *ResultSet = nil
	go func() { res = (<-outChan).(*ResultSet) }()
	MakeDirectory("test", outChan, "does/not/exist")
	if res.Err == "" {
		t.Log("Error should occur.")
		t.FailNow()
	}
	t.Log("Error: ", res.Err)
}

func TestCopy(t *testing.T) {

	home, err := ioutil.TempDir("", "copy")
	FailNotNil(err, t)

	dir := []string{"copy_folder", "copy_folder/sub1",
		"copy_folder/sub2", "copy_folder/sub1/sub3",
		"copy_to"}

	// Create a directory for copying
	for _, p := range dir {
		go MakeDirectory("test", outChan, filepath.Join(home, p))
		res := (<-outChan).(*ResultSet)
		if err := res.Err; err != "" {
			t.Log(err)
			t.FailNow()
		}
	}

	d1 := filepath.Join(home, "copy_folder")
	d2 := filepath.Join(home, "copy_to/")

	defer func() {
		_ = os.Remove(d1)
		_ = os.Remove(d2)
	}()

	go Copy("test", outChan, d1, d2)
	<-outChan
	WaitForOperationsToComplete()

	if CompareDirectory(d1, filepath.Join(d2, "copy_folder")) == false {
		t.Logf("Directories not similar.")
		t.Fail()
	}
}

func TestCopySrcNotExist(t *testing.T) {
	go Copy("test", outChan, "does/not/exist", os.TempDir())
	res := (<-outChan).(*ResultSet)
	if res.Err == "" {
		t.Log("Error should occur.")
		t.FailNow()
	}
	t.Log("Error: ", res.Err)
}

func TestCopyDestNotExist(t *testing.T) {
	f, err := ioutil.TempFile("", "existent_source")
	FailNotNil(err, t)
	f.Close()
	go Copy("test", outChan, f.Name(), "does/not/exist")
	res := (<-outChan).(*ResultSet)
	if res.Err == "" {
		t.Log("Error should occur.")
		t.FailNow()
	}
	t.Log("Error: ", res.Err)
}

func TestGetFile(t *testing.T) {

	var temp, fp string
	var tf *os.File
	outChan := make(chan interface{})

	defer func() {
		_ = os.Remove(temp)
		if tf != nil {
			_ = tf.Close()
		}
		_ = os.Remove(fp)
	}()

	temp, err := ioutil.TempDir("", "GetFile")
	FailNotNil(err, t)
	tf, err = ioutil.TempFile(temp, "getfile")
	FailNotNil(err, t)

	data := []byte{}
	compBuff := []byte{}

	gen := rand.New(rand.NewSource(1))
	size := 3000
	for i := 0; i < int(size); i++ {
		num := gen.Int31()

		b := []byte{0, 0, 0, 0}
		var k int64 = 3
		for k >= 0 {
			b[k] = byte(num & 0xff)
			k--
			num = num >> 8
		}

		for _, j := range b {
			data = append(data, j)
		}
	}

	// File paths
	fp = tf.Name()

	// Write data to temp file and close
	_, err = tf.Write(data)
	FailNotNil(err, t)
	_ = tf.Close()

	// Get the file and write output to outputFile
	go GetFile("test", outChan, fp)
	for {
		res := (<-outChan).(*ResultSet)
		if res.Err != "" {
			t.FailNow()
		}
		for _, b := range res.Data.Data {
			compBuff = append(compBuff, b)
		}
		if res.Data.TotalPieces == res.Data.CurrentPiece {
			break
		}
	}

	WaitForOperationsToComplete()
	if len(data) != len(compBuff) {
		t.FailNow()
	}
}

func TestGetFileEmpty(t *testing.T) {
	var tf *os.File
	var err error
	var res *ResultSet

	tf, err = ioutil.TempFile("", "getfile")
	FailNotNil(err, t)

	_, _ = tf.Write([]byte{})
	_ = tf.Close()

	tp := tf.Name()

	go GetFile("test", outChan, tp)
	res = (<-outChan).(*ResultSet)
	WaitForOperationsToComplete()
	if len(res.Data.Data) != 0 {
		t.FailNow()
	}

}

func TestGetFileNotExist(t *testing.T) {
	t.Skip()
	path := "/this/is/not/a/valid/path"
	var res *ResultSet
	go GetFile("test", outChan, path)
	res = (<-outChan).(*ResultSet)
	WaitForOperationsToComplete()
	if res.Err == "" {
		t.Log("Should fail with an error")
		t.FailNow()
	}
	t.Log("Error: ", res.Err)
}

func TestPutFile(t *testing.T) {
	data := []byte{}

	gen := rand.New(rand.NewSource(1))
	t.Log("Generating bytes: ")
	size := 3000
	for i := 0; i < int(size); i++ {
		num := gen.Int31()

		b := []byte{0, 0, 0, 0}
		var k int64 = 3
		for k >= 0 {
			b[k] = byte(num & 0xff)
			k--
			num = num >> 8
		}

		for _, j := range b {
			data = append(data, j)
		}
	}
	t.Log("Bytes generated: ")

	f, err := ioutil.TempFile("", "put_file_test")
	FailNotNil(err, t)
	newpath := f.Name()
	os.Remove(newpath)
	t.Log(newpath)

	// Write data to newpath using PutFile
	var res *ResultSet
	var count int = 0
	for count < len(data) {
		w := []byte{}
		i := 0
		for i < CHUNKSIZE && count < len(data) {
			w = append(w, data[count])
			i++
			count++
		}
		t.Log("Chunk Length: ", len(w))
		go PutFile("test", outChan, newpath, w)
		res = (<-outChan).(*ResultSet)
		if res.Err != "" {
			t.FailNow()
		}
	}
	file, err := os.OpenFile(newpath, os.O_RDONLY, 0777)
	defer os.Remove(newpath)
	FailNotNil(err, t)
	defer file.Close()
	dataCopy, err := ioutil.ReadAll(file)
	FailNotNil(err, t)
	// t.Log(dataCopy);
	t.Log(len(dataCopy), len(data))

	for i := range data {
		if data[i] != dataCopy[i] {
			t.FailNow()
		}
	}
}

func TestPutFileEmpty(t *testing.T) {
	newpath := filepath.Join(os.TempDir(), "put_file_empty_test")
	go PutFile("test", outChan, newpath, []byte{})
	_ = (<-outChan).(*ResultSet)
	go PutFile("test", outChan, newpath, []byte{})
	_ = (<-outChan).(*ResultSet)
	file, err := os.Open(newpath)
	FailNotNil(err, t)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	FailNotNil(err, t)
	if len(data) != 0 {
		t.FailNow()
	}
}

func TestPutFileNotExists(t *testing.T) {
	path := "this/path/does/not/exist"
	go PutFile("test", outChan, path, []byte{})
	res := (<-outChan).(*ResultSet)
	if res.Err == "" {
		t.Log("Error should occur")
		t.FailNow()
	}
	t.Log("Error: ", res.Err)
}
