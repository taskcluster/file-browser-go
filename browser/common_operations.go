package browser

import (
	"io"
	"os"
	"path/filepath"
)

func init() {
	RegisterCommand("mv", TwoPathWrapper(Move))
	RegisterCommand("rm", OnePathWrapper(Remove))
	RegisterCommand("cp", TwoPathWrapper(Copy))
}

func Move(id string, outChan chan<- *ResultSet, oldpath, newpath string) {
	OpAdd()
	defer func() {
		UnlockPath(oldpath)
		UnlockPath(newpath)
		OpDone()
	}()
	if IsLocked(oldpath) {
		outChan <- FailedResultSet(id, "Path locked for another operation.")
		return
	}
	LockPath(oldpath)
	LockPath(newpath)

	err := os.Rename(oldpath, newpath)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	outChan <- &ResultSet{
		Id: id,
	}
}

func Remove(id string, outChan chan<- *ResultSet, path string) {
	OpAdd()
	defer func() {
		UnlockPath(path)
		OpDone()
	}()
	if IsLocked(path) {
		outChan <- FailedResultSet(id, "Path locked for another operation.")
		return
	}
	LockPath(path)
	err := os.RemoveAll(path)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	outChan <- &ResultSet{
		Id: id,
	}
}

// Function for copying file/dirs

func Copy(id string, outChan chan<- *ResultSet, src, dest string) {
	OpAdd()
	defer OpDone()
	odir, ndir := true, true
	oinfo, err := os.Stat(src)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	odir = oinfo.IsDir()

	ninfo, err := os.Stat(dest)
	if ninfo == nil {
		outChan <- FailedResultSet(id, "FileInfo nil")
		return
	} else {
		ndir = ninfo.IsDir()
	}
	if os.IsNotExist(err) {
		// check if parent exists
		pinfo, err := os.Stat(filepath.Dir(dest))
		if os.IsNotExist(err) {
			outChan <- FailedResultSet(id, err.Error())
			return
		}
		if !pinfo.IsDir() {
			outChan <- FailedResultSet(id, filepath.Dir(dest)+" is not a directory.")
			return
		}
		ndir = false
	}

	// Case 1: if src is a directory and dest is a file
	if odir && !ndir {
		outChan <- FailedResultSet(id, "Directory cannot be copied to file")
		return
	}

	// Case 2: if src is a file and dest is a directory
	if !odir && ndir {
		// Convert to Case 3
		_, f := filepath.Split(src)
		dest = filepath.Join(dest, f)
		ndir = false
	}

	// Case 3: If src and dest are both files
	if !(ndir || odir) {
		of, err := os.Open(src)
		if err != nil {
			outChan <- FailedResultSet(id, err.Error())
			return
		}
		nf, err := os.Open(dest)
		defer func() {
			of.Close()
			nf.Close()
		}()
		if err != nil {
			outChan <- FailedResultSet(id, err.Error())
			return
		}
		_, err = io.Copy(nf, of)
		if err != nil {
			outChan <- FailedResultSet(id, err.Error())
			return
		}
		outChan <- &ResultSet{Id: id}
		return
	}

	// Case 4: src and dest are directories
	_, f := filepath.Split(src)
	dest = filepath.Join(dest, f)
	walkFn := func(path string, info os.FileInfo, err error) error {
		npath := dest + path[len(src):]
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = os.Mkdir(npath, info.Mode().Perm())
			if err != nil {
				return err
			}
		} else {
			file, err := os.OpenFile(npath, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return err
			}
			oldfile, err := os.Open(path)
			if err != nil {
				file.Close()
				return err
			}
			defer func() {
				_ = file.Close()
				_ = oldfile.Close()
			}()
			_, err = io.Copy(file, oldfile)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = filepath.Walk(src, walkFn)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
	}
	outChan <- &ResultSet{
		Id: id,
	}
}
