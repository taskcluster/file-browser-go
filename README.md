## File Browser Go

This is the server side component of the file browser. It reads JSON input from stdin and writes msgpack output to stdout.

Input format:
```
{
	"Cmd": "<command to run>",
	"Args": [ <arguments> ]
}

eg.

{ "Cmd": "ls", "Args": ["/"] }

```

Currently supports:
*	ls
*	cat
