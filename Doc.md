# file-browser-go Documentaion

## Sync Commands

1. ### List
	Lists the contents of a directory.

	Input:
	```json
	{
		"cmd": "list",
		"args": [ "/absolute/path/of/directory/"]
	}
	```

	Result:
	```json
	{
		"cmd": "list",
		"path": "/absolute/path/of/directory/",
		"error": "",
		"files": [ <FileInfo> ]
		
	}
	```
	FileInfo consists of 3 fields:
		- name : string : Path of file or directory
		- size : integer: Size of file
		- dir  : boolean: true if directory

2. ### MkDir
	  Creates a directory at the given path. Similar to mkdir.
  
	  Input:
	  ```json
	  {
		  "cmd": "mkdir",
		  "args": ["/absolute/path/of/new/directory"]
	  }
	  ```
	  
	  Result:
	  ```json
	  {
		  "cmd": "mkdir",
		  "path": "/absolute/path/of/new/directory"
	  }
	  ```
3. ### Remove
	  Removes file or directory at given path (rm -rf).
  
	  Input:
	  ```json
	  {
		  "cmd": "rm",
		  "args": ["/absolute/path/of/directory/or/file"]
	  }
	  ```
	  
	  Output:
	  ```json
	  {
		  "cmd": "rm",
		  "path": "/absolute/path/of/directory/or/file"
	  }
	  ```

4. ### Move
	  Rename/Move a file or directory (mv).
	  
	  Input:
	  ```json
	  {
		  "cmd": "mv",
		  "args": [ "oldpath", "newpath" ]
	  }
	  ```
	  
	  Output:
	  ```json
	  {
		  "cmd": "mv",
		  "path": "newpath"
	  }
	  ```
	  
## Async Commands	
	
1. ### GetFile
