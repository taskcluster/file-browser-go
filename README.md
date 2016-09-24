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

{ "Cmd": "List", "Args": ["/"] }
{ "Cmd": "PutFile", "Args": ["/home/Hello.txt"], "Data":[72, 101, 108, 108, 111] }
{ "Cmd": "GetFile", "Args": ["/home/Hello.txt"] }

```

Currently supports:
*	List
*	GetFile
*	PutFile
