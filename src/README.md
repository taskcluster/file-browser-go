1. _getfile

`_getfile` is a generator function. <br>
The first call to `_getfile` returns a stream object from which data can be read. <br>
The second call to `_getfile` returns a promise which resolves when the whole file has been recieved.

###Basic usage:
```javascript
var gen = browser._getfile('path');
var stream = gen.next().value;
stream.on('data', console.log);
var reader = gen.next().value;
await reader;

```

Two wrapper methods have been provided. <br>
* readFileAsString (remotePath, enc = 'utf8');
* readFileAsBuffer (remotePath);

###Usage:
```javascript
let str = await readFileAsString('path', 'utf8');
let buf = await readFileAsBuffer('path');
```

2. _putfile
The first call to `_putfile` returns a stream to which data can be written. <br>
`_putfile` does not provide an object stream so the data must be of type string or Buffer. <br>
The second call to `_putfile` returns a promise which resolves when `stream.end()` is called on the stream provided by the first call. <br>

###Basic usage:
```javascript
var gen = browser._putfile('path');
var stream = gen.next().value;
var finished = gen.next().value;
stream.write('string or buffer');
stream.end();
await finished;
```
