package browser;

import (
	"testing";
	"os";
	"path/filepath";
	"io";
	"io/ioutil";
	"math/rand";
	"time";
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

	gen := rand.New(rand.NewSource(time.Now().Unix()));
	size := gen.Int31() % 1000 + 1; // Max size would be 4000 bytes
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
