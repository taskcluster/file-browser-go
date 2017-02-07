## File Browser Go

This is the server side component of the file browser. It reads JSON input from stdin and writes msgpack output to stdout.

Input format:
```
{
	"Cmd": "<command to run>",
	"Args": [ <arguments> ],
	"Data": [] //Byte array
}

eg.

{ "cmd": "ls", "args": ["/"] }
{ "cmd": "putfile", "args": ["/home/Hello.txt"], "data":[72, 101, 108, 108, 111] }
{ "cmd": "getfile", "args": ["/home/Hello.txt"] }
{ "cmd": "mkdir", "args": ["/home/Folder/"] }
{ "cmd": "rm", "args": ["/path/to/file/or/folder"] }
{ "cmd": "mv", "args": ["path/to/file/or/folder", "target/path"] }
{ "cmd": "cp", "args": ["src/path", "dest/path"]}
```

Currently supports:
*	List
*	GetFile
*	PutFile
*	MakeDir
*	Remove
*	Move
*	Copy


### GetFile
GetFile initially writes a ResultSet object
with total number of pieces and piece number 0.
```
{
	Cmd: "GetFile",
	Path: < path given in command >,
	Data: {
		TotalPieces: <Total number of pieces>,
		CurrentPiece: 0,
		Data: nil,
	}
}
```
Use TotalPieces and CurrentPiece to reassemble the file.
```
{
	Cmd: "GetFile",
	Path: < path given in command >,
	Data: {
		TotalPieces: <Total number of pieces>,
		CurrentPiece: <Current piece>,
		Data: [ ... ] ,
	}
}
```

### PutFile

Send the file data to the binary in chunks of 2048 bytes.
```
{
	Cmd: "PutFile",
	Args: [ <path of file> ],
	Data: [ ... ], 
}
```
PutFile will append the bytes to a temp file.
After all file data has been sent, send a command of the form

```
{
	Cmd: "PutFile",
	Args: [ <path of file> ],
	Data: [], // Empty array
}
```
This is a signal indicating that all bytes have been sent and 
the temp file is moved to the intended path.
