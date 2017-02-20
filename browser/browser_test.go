package browser

import (
	"bytes"
	"github.com/taskcluster/slugid-go/slugid"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestList(t *testing.T) {
	outChan := make(chan *ResultSet)
	temp, err := ioutil.TempDir("", "list")
	FailNotNil(err, t)
	defer func() {
		os.RemoveAll(temp)
		close(outChan)
	}()
	// Make two directories in temp
	err = os.Mkdir(filepath.Join(temp, "a"), 0777)
	FailNotNil(err, t)
	err = os.Mkdir(filepath.Join(temp, "b"), 0777)
	FailNotNil(err, t)
	go List("test", outChan, temp)
	res := (<-outChan)
	if res.Err != "" {
		t.Fail()
	}
	if len(res.Files) != 2 {
		t.Fail()
	}
}

func TestListNotExist(t *testing.T) {
	outChan := make(chan *ResultSet)
	defer close(outChan)
	go List("test", outChan, "does/not/exist")
	res := (<-outChan)
	if res.Err == "" {
		t.Log("Error should occur.")
		t.Fail()
	}
}

func TestMakeDirectoryAndRemove(t *testing.T) {
	outChan := make(chan *ResultSet)
	dummyChan := make(chan *ResultSet, 10)
	home, err := ioutil.TempDir("", "MakeAndRemove")
	FailNotNil(err, t)
	paths := []string{"/test_folder", "/test_folder/sub_folder"}

	defer func() {
		Remove("test", dummyChan, filepath.Join(home, paths[0]))
		close(outChan)
		close(dummyChan)
	}()

	for _, p := range paths {
		go MakeDirectory("test", outChan, filepath.Join(home, p))
		res := (<-outChan)
		if res.Err != "" {
			t.Fatal(res.Err)
		}
		if !Exists(filepath.Join(home, p)) {
			t.Fatal("Directory not created")
			t.Fail()
		}
	}
}

func TestMakeDirectoryBadPath(t *testing.T) {
	outChan := make(chan *ResultSet)
	defer close(outChan)
	go MakeDirectory("test", outChan, "does/not/exist")
	res := <-outChan
	if res.Err == "" {
		t.Fatal("Error should occur.")
	}
}

func TestCopy(t *testing.T) {
	outChan := make(chan *ResultSet)
	d1, d2 := "", ""
	home, err := ioutil.TempDir("", "copy"+slugid.V4())
	FailNotNil(err, t)

	defer func() {
		_ = os.RemoveAll(home)
	}()

	dir := []string{"copy_folder", "copy_folder/sub1",
		"copy_folder/sub2", "copy_folder/sub1/sub3",
		"copy_to"}

	// Create a directory for copying
	for _, p := range dir {
		err := os.Mkdir(filepath.Join(home, p), 0777)
		FailNotNil(err, t)
	}

	d1 = filepath.Join(home, "copy_folder")
	d2 = filepath.Join(home, "copy_to/")

	go Copy("test", outChan, d1, d2)
	res := (<-outChan)
	if res.Err != "" {
		t.Fatal(res.Err)
	}

	if CompareDirectory(d1, filepath.Join(d2, "copy_folder")) == false {
		t.Fatal("Directories not similar.")
	}
}

func TestCopySrcNotExist(t *testing.T) {
	outChan := make(chan *ResultSet)
	defer close(outChan)
	go Copy("test", outChan, "does/not/exist", os.TempDir())
	res := (<-outChan)
	if res.Err == "" {
		t.Fatal(res.Err)
	}
}

func TestCopyDestNotExist(t *testing.T) {
	outChan := make(chan *ResultSet)
	defer close(outChan)
	f, err := ioutil.TempFile("", "existent_source")
	FailNotNil(err, t)
	f.Close()
	go Copy("test", outChan, f.Name(), "does/not/exist")
	res := (<-outChan)
	if res.Err == "" {
		t.Fatal(res.Err)
	}
}

func TestGetFile(t *testing.T) {
	outChan := make(chan *ResultSet)
	temp, err := ioutil.TempDir("", "GetFile")
	FailNotNil(err, t)
	tf, err := ioutil.TempFile(temp, "getfile")
	FailNotNil(err, t)

	// File paths
	fp := tf.Name()

	defer func() {
		_ = os.RemoveAll(temp)
		if tf != nil {
			_ = tf.Close()
		}
		_ = os.RemoveAll(fp)
		close(outChan)
	}()

	data := []byte{}
	compBuff := []byte{}

	gen := rand.New(rand.NewSource(1))
	size := 3000
	for i := 0; i < int(size); i++ {
		num := gen.Int31()
		for j := uint(0); j < uint(4); j++ {
			b := byte(num >> (j * 8))
			data = append(data, b)
		}
	}

	// Write data to temp file and close
	_, err = tf.Write(data)
	FailNotNil(err, t)
	_ = tf.Close()

	// Get the file and write output to outputFile
	go GetFile("test", outChan, fp)
	res := (<-outChan)
	compBuff = append(compBuff, res.Data.Data...)
	for res.Data.TotalPieces != res.Data.CurrentPiece {
		res = (<-outChan)
		if res.Err != "" {
			t.FailNow()
		}
		compBuff = append(compBuff, res.Data.Data...)
	}

	if len(data) != len(compBuff) {
		t.FailNow()
	}
}

func TestGetFileEmpty(t *testing.T) {
	outChan := make(chan *ResultSet)
	tf, err := ioutil.TempFile("", "getfile")
	FailNotNil(err, t)

	_, _ = tf.Write([]byte{})
	_ = tf.Close()

	tp := tf.Name()

	go GetFile("test", outChan, tp)
	res := (<-outChan)
	if len(res.Data.Data) != 0 {
		t.FailNow()
	}

}

func TestGetFileNotExist(t *testing.T) {
	outChan := make(chan *ResultSet)
	path := "/this/is/not/a/valid/path"
	go GetFile("test", outChan, path)
	res := <-outChan
	if res.Err == "" {
		t.Fatal("Error should occur")
	}
}

func TestPutFile(t *testing.T) {
	outChan := make(chan *ResultSet)
	data := []byte{}

	gen := rand.New(rand.NewSource(1))
	size := 3000
	for i := 0; i < int(size); i++ {
		num := gen.Int31()

		for j := uint(0); j < 4; j++ {
			b := byte(num >> (j * 8))
			data = append(data, b)
		}
	}

	f, err := ioutil.TempFile("", "put_file_test")
	FailNotNil(err, t)
	newpath := f.Name()
	os.Remove(newpath)

	// Write data to newpath using PutFile
	var count int
	for count < len(data) {
		i := Min(CHUNKSIZE, len(data)-count)
		w := data[count : count+i]
		count += i
		go PutFile("test", outChan, newpath, w)
		res := (<-outChan)
		if res.Err != "" {
			t.Fatal(res.Err)
		}
	}
	file, err := os.OpenFile(newpath, os.O_RDONLY, 0777)
	defer os.Remove(newpath)
	FailNotNil(err, t)
	defer file.Close()
	dataCopy, err := ioutil.ReadAll(file)
	FailNotNil(err, t)

	if !bytes.Equal(data, dataCopy) {
		t.FailNow()
	}
}

func TestPutFileEmpty(t *testing.T) {
	outChan := make(chan *ResultSet)
	newpath := filepath.Join(os.TempDir(), "put_file_empty_test")
	go PutFile("test", outChan, newpath, []byte{})
	_ = (<-outChan)
	go PutFile("test", outChan, newpath, []byte{})
	_ = (<-outChan)
	file, err := os.Open(newpath)
	FailNotNil(err, t)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	FailNotNil(err, t)
	if len(data) != 0 {
		t.FailNow()
	}
}

func TestPutFileBadPath(t *testing.T) {
	outChan := make(chan *ResultSet)
	path := "this/path/does/not/exist"
	go PutFile("test", outChan, path, []byte{})
	res := (<-outChan)
	if res.Err == "" {
		t.Fatal("Should fail")
	}
}
