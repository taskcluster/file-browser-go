package browser;

import (
	"testing";
	"os";
	"path/filepath";
	"io";
	"io/ioutil";
	"math/rand";
	"encoding/json";
)

func TestList(t *testing.T) {
	temp, err := ioutil.TempDir("","list");
	FailNotNil(err, t);
	defer os.Remove(temp);
	// Make two directories in temp
	err = os.Mkdir(filepath.Join(temp, "a"), 0777);
	FailNotNil(err, t);
	err = os.Mkdir(filepath.Join(temp, "b"), 0777);
	FailNotNil(err, t);
	res := List(temp).(*ResultSet);
	if res.Err != "" {
		t.Fail();
	}
	if len(res.Files) != 2 {
		t.Fail();
	}
}

func TestListNotExist(t *testing.T) {
	res := List("does/not/exist").(*ResultSet);
	if res.Err == "" {
		t.Log("Error should occur.");
		t.Fail();
	}
	t.Log("Error: ", res.Err);
}

func TestMakeDirectoryAndRemove(t *testing.T) {
	home , err:= ioutil.TempDir("","MakeAndRemove");
	FailNotNil(err, t);
	paths := []string{"/test_folder","/test_folder/sub_folder"};

	defer func() {
		r := Remove(filepath.Join(home, paths[0])).(*ResultSet);
		if r.Err != "" {
			t.Log(r.Err);
			t.Fail();
		}
	}();

	for _, p := range paths {
		res := MakeDirectory(filepath.Join(home, p)).(*ResultSet);
		if res.Err != "" {
			t.Log(res.Err);
			t.Fail();
		}
		if !Exists(filepath.Join(home, p)) {
			t.Log("Directory not created");
			t.Fail();
		}
	}
}

func TestMakeDirectoryBadPath(t *testing.T) {
	res := MakeDirectory("does/not/exist").(*ResultSet);
	if res.Err == "" {
		t.Log("Error should occur.");
		t.FailNow();
	}
	t.Log("Error: ", res.Err);
}

func TestCopy (t *testing.T) {

	home, err := ioutil.TempDir("","copy");
	FailNotNil(err, t);

	dir := []string{"copy_folder", "copy_folder/sub1",
	"copy_folder/sub2", "copy_folder/sub1/sub3",
	"copy_to"};

	// Create a directory for copying
	for _, p := range dir {
		res := MakeDirectory(filepath.Join(home,p));
		if err := res.(*ResultSet).Err; err != "" {
			t.Log(err);
			t.FailNow();
		}
	}

	d1 := filepath.Join(home,"copy_folder");
	d2 := filepath.Join(home,"copy_to/");

	defer func() {
		_ = Remove(d1);
		_ = Remove(d2);
	}();

	res := Copy(d1, d2, ioutil.Discard).(*ResultSet);
	WaitForOperationsToComplete();

	if CompareDirectory(d1,res.Path) == false {
		t.Logf("Directories not similar.");
		t.Fail();
	}
}

func TestCopySrcNotExist (t *testing.T){
	res := Copy("does/not/exist", os.TempDir(), ioutil.Discard).(*ResultSet);
	if res.Err == "" {
		t.Log("Error should occur.");
		t.FailNow();
	}
	t.Log("Error: ", res.Err);
}

func TestCopyDestNotExist(t *testing.T) {
	f, err := ioutil.TempFile("", "existent_source");
	FailNotNil(err, t);
	f.Close();
	res := Copy(f.Name(), "does/not/exist", ioutil.Discard).(*ResultSet);
	if res.Err == "" {
		t.Log("Error should occur.");
		t.FailNow();
	}
	t.Log("Error: ", res.Err);
}

func TestGetFile (t *testing.T) {

	var temp string;
	var tf, outputFile *os.File;

	defer func() {
		_ = os.Remove(temp);
		if tf != nil {
			_ = tf.Close();
		}
	}();

	temp, err := ioutil.TempDir("", "GetFile");
	FailNotNil(err, t);
	tf, err = ioutil.TempFile(temp, "getfile");
	FailNotNil(err, t);
	outputFile, err = os.OpenFile(filepath.Join(temp, "catch"), os.O_CREATE | os.O_WRONLY, 0777);
	FailNotNil(err, t);

	data := []byte{};

	gen := rand.New(rand.NewSource(1));
	size := 3000;
	for i := 0; i < int(size); i++ {
		num := gen.Int31();

		b := []byte{0,0,0,0};
		var k int64 = 3;
		for k >= 0 {
			b[k] = byte(num & 0xff);
			k--;
			num = num >> 8;
		}

		for _, j := range b {
			data = append(data, j);
		}
	}

	// File paths
	fp := tf.Name();
	op := outputFile.Name();

	defer func() {
		_ = os.Remove(fp);
		_ = os.Remove(op);
	}();

	// Write data to temp file and close
	_, err = tf.Write(data);
	FailNotNil(err, t);
	_ = tf.Close();

	// Get the file and write output to outputFile
	GetFile(fp, outputFile);
	WaitForOperationsToComplete();
	// Close the file after writing
	outputFile.Close();

	// Reopen file in READ ONLY mode
	outputFile, err = os.OpenFile(op, os.O_RDONLY, 0777);
	FailNotNil(err, t);

	// Create decoder from GetFile output
	jsonDec := json.NewDecoder(outputFile);
	err = nil;
	res := &ResultSet{};

	// This helps to ensure that the pieces are in order
	var max int64 = -1;

	// Buffer to hold the data that's been read
	compBuff := []byte{};

	for err != io.EOF {
		err = jsonDec.Decode(res);
		if err != nil {
			t.Log(err);
			break;
		}
		t.Log(res);
		if res.Data.CurrentPiece == 0 {
			t.Logf("Got total pieces %d", res.Data.TotalPieces);
			max = res.Data.TotalPieces;
			continue;
		}
		for _, b := range res.Data.Data {
			compBuff = append(compBuff, b);
		}
		max--;
	}

	outputFile.Close();

	if max != 0 {
		t.Logf("All pieces were not recieved");
		t.FailNow();
	}
	for i, _ := range data {
		if data[i] != compBuff[i]{
			t.Logf("Data not same");
			t.FailNow();
		}
	}
}

func TestGetFileEmpty(t *testing.T) {
	var tf, of *os.File;
	var err error;
	tf, err = ioutil.TempFile("","getfile");
	FailNotNil(err, t);
	of, err = ioutil.TempFile("", "catch");

	_,_ = tf.Write([]byte{});
	_ = tf.Close();

	tp := tf.Name();
	op := of.Name();
	GetFile(tp, of);
	WaitForOperationsToComplete();

	of.Close();

	of, err = os.OpenFile(op, os.O_RDONLY, 0777);
	jsonDec := json.NewDecoder(of);

	var max int64 = -1;
	err = nil;
	res := &ResultSet{};

	tempBuff := []byte{};

	for err != io.EOF {
		err = jsonDec.Decode(res);
		if err != nil {
			t.Log(err);
			break;
		}
		if res.Data.CurrentPiece == 0 {
			t.Logf("Got total pieces %d", res.Data.TotalPieces);
			max = res.Data.TotalPieces;
			continue;
		}
		max--;
		for _, b := range res.Data.Data {
			tempBuff = append(tempBuff, b);
		}
	}

	if max != 0 {
		t.Log("All pieces not recieved.");
		t.FailNow();
	}
	// t.Log(tempBuff);
	if len(tempBuff) != 0 {
		t.Logf("%d bytes expected.", 0);
		t.FailNow();
	}
}

func TestGetFileNotExist(t *testing.T) {
	path := "/this/is/not/a/valid/path";
	op, err := ioutil.TempFile("", "get_file_not_exist_op");
	FailNotNil(err, t);

	GetFile(path, op);
	WaitForOperationsToComplete();
	op.Close();

	op, err = os.Open(op.Name());
	FailNotNil(err, t);

	dec := json.NewDecoder(op);
	var res *ResultSet = nil;
	dec.Decode(&res);
	if res.Err == "" {
		t.Log("Should return error.");
		t.FailNow();
	}

	t.Log("Error: ", res.Err);
}

func TestPutFile (t *testing.T) {
	data := []byte{};

	gen := rand.New(rand.NewSource(1));
	t.Log("Generating bytes: ");
	size := 3000;
	for i := 0; i < int(size); i++ {
		num := gen.Int31();

		b := []byte{0,0,0,0};
		var k int64 = 3;
		for k >= 0 {
			b[k] = byte(num & 0xff);
			k--;
			num = num >> 8;
		}

		for _, j := range b {
			data = append(data, j);
		}
	}
	t.Log("Bytes generated: ");

	f, err := ioutil.TempFile("", "put_file_test");
	FailNotNil(err, t);
	newpath := f.Name();
	os.Remove(newpath);
	t.Log(newpath);

	// Write data to newpath using PutFile2
	var count int = 0;
	for count < len(data) {
		w := []byte{};
		i := 0;
		for i < CHUNKSIZE && count < len(data) {
			w = append(w, data[count]);
			i++;
			count++;
		}
		t.Log("Chunk Length: ", len(w));
		res := PutFile2(newpath, w).(*ResultSet);
		if res.Err != "" {
			t.FailNow();

			_ = PutFile(newpath, []byte{});
		}
	}
	_ = PutFile2(newpath, []byte{});
	file, err := os.OpenFile(newpath, os.O_RDONLY, 0777);
	defer os.Remove(newpath);
	FailNotNil(err, t);
	defer file.Close();
	dataCopy, err := ioutil.ReadAll(file);
	FailNotNil(err ,t);
	// t.Log(dataCopy);
	t.Log(len(dataCopy) , len(data));

	for i, _ := range data {
		if data[i] != dataCopy[i] {
			t.FailNow();
		}
	}
}

func TestPutFileEmpty(t *testing.T) {
	newpath := filepath.Join(os.TempDir(), "put_file_empty_test");
	PutFile2(newpath, []byte{});
	PutFile2(newpath, []byte{});
	file, err := os.Open(newpath);
	FailNotNil(err, t);
	defer file.Close();
	data, err := ioutil.ReadAll(file);
	FailNotNil(err, t);
	if len(data) != 0 {
		t.FailNow();
	}
}

func TestPutFileNotExists(t *testing.T) {
	path := "this/path/does/not/exist";
	res := PutFile2(path, []byte{}).(*ResultSet);
	if res.Err == "" {
		t.Log("Error should occur");
		t.FailNow();
	}
	t.Log("Error: ", res.Err);
}
