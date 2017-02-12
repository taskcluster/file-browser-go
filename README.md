##File Browser Go

Go binary and JS api of the file browser. Itreads msgpack input from stdin and writes msgpack output to stdout.

Currently supports:
* ls -> List contents of a directory
* cp -> copy contents (within remote fs) 
* mv -> move contents (within remote fs)
* rm -> remove file or directory (within remote fs)
* getfile -> stream file from remote fs to local fs
* putfile -> stream file from local fs to remote fs

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
`getfile` and `putfile` are generator functions. <br>
For `getfile` and `putfile` check src/README.md .
