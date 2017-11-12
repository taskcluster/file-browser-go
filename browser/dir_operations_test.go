// +build stateful

package browser

import (
	"context"
	"os"
	"testing"
)

const baseTestDir = "test_dir"

func Test_Sl_readdir_success(t *testing.T) {
	baseCtx := context.Background()
	req := &Sl_readdirRequest{
		RequestID: 0,
		Path:      baseTestDir,
	}

	res := sl_readdir(baseCtx, req, nil)
	r := res.(*ReaddirResponse)
	if r.Error != nil {
		t.Fatal(r.Error)
	}
	if len(r.Entries) != 2 {
		t.Fatal("wrong count")
	}
}

func Test_Sl_readdir_file(t *testing.T) {
	baseCtx := context.Background()
	req := &Sl_readdirRequest{
		RequestID: 0,
		Path:      baseTestDir + "/read_large_file",
	}

	res := sl_readdir(baseCtx, req, nil)
	r := res.(*ReaddirResponse)
	if r.Error == nil {
		t.Fatal("should fail")
	}
}

func Test_Sl_readdir_BadPath(t *testing.T) {
	baseCtx := context.Background()
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
	baseCtx := context.Background()
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
	if len(readdirRes.Entries) != 2 {
		t.Log(readdirRes)
		t.Fatal("wrong count")
	}
}

func Test_Sf_open_and_readdir_file(t *testing.T) {
	baseCtx := context.Background()
	openReq := &Sf_openRequest{
		RequestID: 0,
		Path:      baseTestDir + "/read_small_file",
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
