##File Browser Go

Go binary and JS api of the file browser. It reads msgpack input from stdin and writes msgpack output to stdout.

Currently supports:
* ls -> List contents of a directory
* cp -> copy contents (within remote fs) 
* mv -> move contents (within remote fs)
* rm -> remove file or directory (within remote fs)
* mkdir -> Create a directory on the remote fs
* getfile -> stream file from remote fs to local fs
* putfile -> stream file from local fs to remote fs

###Creating a browser object:
* FileBrowser(inStream, outStream);

eg. 
```javascript
let shell = child_process.spawn('file-browser-go', ...);
let browser = new FileBrowser(shell.stdin, shell.stdout);
```

###Basic usage

1. ls -> list
```js
var result = await browser.ls('path');
console.log(result.files);
```
2. mv -> move
```js
var result = await browser.mv('src', 'dest');
```
3. cp -> copy
```javascript
var result = await browser.cp('src', 'dest');
```
4. rm -> remove
```javascript
var result = await browser.rm('path');
```
5. mkdir -> make directory
```javascript
var result = await browser.mkdir('path');
```
`getfile` and `putfile` are generator functions. <br>
For `getfile` and `putfile` check src/README.md .
