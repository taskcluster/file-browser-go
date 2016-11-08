'use strict';

let fs      = require('fs');
let Buffer  = require('buffer');
let Command = require('./Command.js');

// Maps path on the container to path in the local fs.
let getfilemap = {};
let putfilemap = {};

let FileOperations = {};

const CHUNKSIZE = 2048;

FileOperations.getfileSetPath = (path, localPath) => {
  getfilemap[path] = localPath;
}

FileOperations.getfileWrite = (path, data, complete) => { 
  // Create local file to hold data
  if (!getfilemap[path]){
    return;
  }
  let localPath = getfilemap[path];
  fs.appendFileSync(localPath, data);
  if (complete) {
    getfilemap[path] = null;
  }
}

FileOperations.putfile = (localPath, path) => {
  if (!putfilemap[path]) {

    if (fs.existsSync(localPath) && fs.statSync(localPath).isFile()){
      let fd = fs.openSync(localPath);
      putfilemap[path] = [ fd, 0 ];
    }else{
      return null;
    }

  }
  
  let fd = putfilemap[path][0];
  let pos = putfilemap[path][1];

  let buf = Buffer.alloc(CHUNKSIZE);

  let bytes = fs.readSync(fd, buf, 0, buf.length, pos);
  let trailingCommand = null;
  let command = null;

  // Blank file
  if(bytes == 0 && pos == 0) {
    trailingCommand = Command.putfile(path, Buffer.alloc(0));
    putfilemap[path] = null;
    return [trailingCommand, trailingCommand];
  }

  // Special case which may occur if size of file
  // is a multiple of CHUNKSIZE
  if(bytes == 0) {
    putfilemap[path] = null;
    return [ Command.putfile(path, Buffer.alloc(0)) ];
  }

  if (bytes < CHUNKSIZE) {
    trailingCommand = Command.putfile(path, Buffer.alloc(0));
    putfilemap[path] = null;
  }

  command = Command.putfile(path, buf);

  // Increment position
  putfilemap[path][1] += bytes;

  if (trailingCommand) {
    return [ command, trailingCommand ];
  }else {
    return [ command ];
  }

}


module.exports = FileOperations;
