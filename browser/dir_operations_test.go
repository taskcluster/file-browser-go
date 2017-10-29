// +build stateful

package browser

import (
	"context"
	"os"
	"testing"
)

const baseTestDir = "test_dir"

var baseCtx = context.Background()

func Test_Sl_readdir_success(t *testing.T) {
	req := &Sl_readdirRequest{
		RequestID: 0,
		Path:      baseTestDir,
	}

	res := sl_readdir(baseCtx, req, nil)
	r := res.(*ReaddirResponse)
	if r.Error != nil {
		t.Fatal(r.Error)
	}
	if len(r.Entries) != 7 {
		t.Fatal("wrong count")
	}
}

func Test_Sl_readdir_file(t *testing.T) {
	req := &Sl_readdirRequest{
		RequestID: 0,
		Path:      baseTestDir + "/file_0",
	}

	res := sl_readdir(baseCtx, req, nil)
	r := res.(*ReaddirResponse)
	if r.Error == nil {
		t.Fatal("should fail")
	}
}

func Test_Sl_readdir_BadPath(t *testing.T) {
	req := &Sl_readdirRequest{
		RequestID: 0,
		Path:      baseTestDir + "/does/not/exist/",
	}

	res := sl_readdir(baseCtx, req, nil)
	r := res.(*ReaddirResponse)
	if r.Error == nil {
		t.Fatal("should fail")
	}
}

func Test_Sf_open_and_readdir(t *testing.T) {
	openReq := &Sf_openRequest{
		RequestID: 0,
		Path:      baseTestDir,
		Flags:     os.O_RDONLY,
		Perm:      os.FileMode(0444),
	}

	res := sf_open(baseCtx, openReq, nil)

	openRes := res.(*Sf_openResponse)
	if openRes.Error != nil {
		t.Fatal(openRes.Error)
	}

	defer func() {
		_ = localFileRegistry.closeFile(openRes.HandleID)
	}()

	readdirReq := &Sf_readdirRequest{
		RequestID: 1,
		HandleID:  openRes.HandleID,
	}

	res = sf_readdir(baseCtx, readdirReq, nil)
	readdirRes := res.(*ReaddirResponse)
	if readdirRes.Error != nil {
		t.Fatal(readdirRes.Error)
	}
	if len(readdirRes.Entries) != 7 {
		t.Log(readdirRes)
		t.Fatal("wrong count")
	}
}

func Test_Sf_open_and_readdir_file(t *testing.T) {
	openReq := &Sf_openRequest{
		RequestID: 0,
		Path:      baseTestDir + "/file_0",
		Flags:     os.O_RDONLY,
		Perm:      os.FileMode(0444),
	}

	res := sf_open(baseCtx, openReq, nil)

	openRes := res.(*Sf_openResponse)
	if openRes.Error != nil {
		t.Fatal(openRes.Error)
	}

	defer func() {
		_ = localFileRegistry.closeFile(openRes.HandleID)
	}()

	readdirReq := &Sf_readdirRequest{
		RequestID: 1,
		HandleID:  openRes.HandleID,
	}

	res = sf_readdir(baseCtx, readdirReq, nil)

	readdirRes := res.(*ReaddirResponse)
	if readdirRes.Error == nil {
		t.Fatal("should fail")
	}
}
