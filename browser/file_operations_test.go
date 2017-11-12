package browser

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// Tests for stateless reads

// Try reading a file that does not exist
func Test_Sl_read_FileDoesNotExist(t *testing.T) {
	ctx := context.Background()
	req := &Sl_readRequest{
		RequestID:   0,
		BytesToRead: 1024, // Note that this file only contains 100 bytes
		Path:        baseTestDir + "/does_not_exist",
		Offset:      0,
	}
	res := (sl_read(ctx, req, nil)).(*ReadResponse)
	defer res.closeFile()

	if res.Error == nil {
		t.Fatal("error should be returned")
	}
}

// Read from small file in test_dir to check if EOF is returned
func Test_Sl_read_EOF(t *testing.T) {
	ctx := context.Background()
	req := &Sl_readRequest{
		RequestID:   0,
		BytesToRead: 1024, // Note that this file only contains 100 bytes
		Path:        baseTestDir + "/read_small_file",
		Offset:      0,
	}
	r := sl_read(ctx, req, nil)
	res := r.(*ReadResponse)

	if !res.IsStreamResponse() {
		t.Fatal("response should be streamed")
	}

	out := make(chan interface{}, 1)
	done := make(chan struct{})
	go func() {
		_ = res.StreamToChannel(out)
		close(done)
	}()

	var err error
	data := []byte{}
	for err == nil && len(data) < 1024 {
		rf := (<-out).(*ReadResponseFrame)
		err = rf.Error
		if rf.Bytes != 0 {
			data = append(data, rf.Buffer...)
		}
	}
	<-done
	if len(data) != 100 {
		res.closeFile()
		t.Fatal("buffer should only contain 100 bytes")
	}
	if err != io.EOF {
		res.closeFile()
		t.Fatalf("expected: %v, got error: %v", io.EOF, res.Error)
	}

	file, err := os.Open(baseTestDir + "/read_small_file")
	if err != nil {
		res.closeFile()
		t.Fatalf("unable to open file for checking: %v", err)
	}
	testData, err := ioutil.ReadAll(file)
	_ = file.Close()
	if err != nil {
		res.closeFile()
		t.Fatalf("could not read file for checking: %v", err)
	}
	if !bytes.Equal(data, testData) {
		res.closeFile()
		t.Fatal("data should not differ")
	}

	res.closeFile()
}

// Try reading from offset in file
func Test_sl_read_ReadFromOffset(t *testing.T) {
	ctx := context.Background()
	req := &Sl_readRequest{
		RequestID:   0,
		Offset:      1234,
		BytesToRead: 5000,
		Path:        baseTestDir + "/read_large_file",
	}
	res := (sl_read(ctx, req, nil)).(*ReadResponse)
	defer res.closeFile()

	if res.Error != nil {
		res.closeFile()
		t.Fatalf("expected: %v, got error: %v", nil, res.Error)
	}

	if !res.IsStreamResponse() {
		res.closeFile()
		t.Fatalf("response should be streamed")
	}

	data, err := streamReadResponse(res)

	// file is closed by this point
	if err != nil {
		t.Fatalf("streaming error: %v", err)
	}

	file, err := os.Open(baseTestDir + "/read_large_file")
	if err != nil {
		t.Fatalf("unable to open file for checking: %v", err)
	}

	defer file.Close()

	testData := make([]byte, 5000)
	n, err := file.ReadAt(testData, 1234)

	testData = testData[:n]

	if err != nil {
		file.Close()
		t.Fatalf("error while reading file: %v", err)
	}

	if !bytes.Equal(data, testData) {
		file.Close()
		t.Fatal("data should not differ")
	}
}
