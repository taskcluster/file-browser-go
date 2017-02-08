## File Browser Go

Go binary and JS api of the file browser. Itreads msgpack input from stdin and writes msgpack output to stdout.

Currently supports:
* ls -> List contents of a directory
* cp -> copy contents (within remote fs) 
* mv -> move contents (within remote fs)
* rm -> remove file or directory (within remote fs)
* getfile -> stream file from remote fs to local fs
* putfile -> stream file from local fs to remote fs
