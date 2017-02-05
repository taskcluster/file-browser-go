'use strict';

let fs      = require('fs');
let Buffer  = require('buffer');
let Command = require('./Command.js');

// Maps path on the container to path in the local fs.
let putfilemap = {};

let FileOperations = {};

const CHUNKSIZE = 2048;

// dest -> File on remote fs
// src -> File on local fs

FileOperations.putfile = (src, dest) => {
  if (!putfilemap[dest]) {

    if (fs.existsSync(src) && fs.statSync(src).isFile()){
      let fd = fs.openSync(src);
      putfilemap[dest] = [ fd, 0 ];
    }else{
      return null;
    }

  }
  
  let fd = putfilemap[dest][0];
  let pos = putfilemap[dest][1];

  let buf = Buffer.alloc(CHUNKSIZE);

  let bytes = fs.readSync(fd, buf, 0, buf.length, pos);
  let trailingCommand = null;
  let command = null;

  // Blank file
  if(bytes == 0 && pos == 0) {
    trailingCommand = Command.putfile(dest, Buffer.alloc(0));
    putfilemap[dest] = null;
    return [trailingCommand, trailingCommand];
  }

  // Special case which may occur if size of file
  // is a multiple of CHUNKSIZE
  if(bytes == 0) {
    putfilemap[dest] = null;
    return [ Command.putfile(dest, Buffer.alloc(0)) ];
  }

  if (bytes < CHUNKSIZE) {
    trailingCommand = Command.putfile(dest, Buffer.alloc(0));
    putfilemap[dest] = null;
  }

  command = Command.putfile(dest, buf);

  // Increment position
  putfilemap[dest][1] += bytes;

  if (trailingCommand) {
    return [ command, trailingCommand ];
  }else {
    return [ command ];
  }

}

FileOperations.putfileClean = (dest) => {
  putfilemap[dest] = null;  
}


module.exports = FileOperations;
