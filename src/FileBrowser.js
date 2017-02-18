'use strict';
/* API for all file-browser commands */

const
  assert    = require('assert'),
  debug     = require('debug')('browser'),
  Buffer    = require('buffer').Buffer,
  msp       = require('msgpack-lite'),
  through2  = require('through2'),
  Promise   = require('bluebird'),
  lock      = require('lock')(),
  _         = require('lodash'),
  slugid    = require('slugid'),
  fs        = require('fs');

const
  Command   = require('./Command.js');

const CHUNKSIZE = 4072;

const lck = key => new Promise(resolve => lock(key, resolve));

class FileBrowser {
  
  constructor(inStream, outStream) {

    assert(inStream);
    assert(outStream);

    this.stdin = through2.obj(function (chunk, enc, cb) {
      this.push(msp.encode(chunk));
      cb();
    });

    // Table that maps from id to resolver
    this._cb = {}

    this.stdout = msp.createDecodeStream();
    // this.testOut = msp.createDecodeStream();

    this.stdin.pipe(inStream);
    // this.stdin.pipe(this.testOut);
    outStream.pipe(this.stdout);

    // this.stdout.setMaxListeners(0);

    // this.testOut.on('data', (data) => {
    //   debug("Wrote: ", data);
    // });

    this.stdout.on('data', (data) => {
      if(!data.id) return;
      debug('Received: ', data);
      if (!this._cb[data.id]) {
        debug("message for unknown id: ", data.id);
        return;
      }

      this._cb[data.id](data);
    });

    this.stdout.on('error', debug);
    this.stdin.on('error', debug);

    let self = this;

    ["mv", "cp"].forEach(c => {

      self[c] = async (src, dest) => {
        let cmd = Command[c](src, dest);
        let result = await self._writeAndResolve(cmd);
        delete this._cb[cmd.id];
        if (result.error) throw new Error(result.error);
        return result;
      }

    });

    ["ls", "rm", "mkdir"].forEach(c => {

      self[c] = async (path) => {
        let cmd = Command[c](path);
        let result = await self._writeAndResolve(cmd);
        delete this._cb[cmd.id];
        if (result.error) throw new Error(result.error);
        return result;
      }

    });
  }
  
  _writeAndResolve (cmd) {
    let self = this;
    return new Promise(resolve => {
      self.stdin.write(cmd);
      this._cb[cmd.id] = resolve;
    });
  }

  * _putfile (dest) {  

    if(typeof dest != 'string') {
      throw new Error('\'dest\' must be of type string.');
    }
    let self = this;

    let srcStream = through2(function (chunk, enc, cb) { this.push(chunk); cb(); });
    yield srcStream;

    yield new Promise((resolve, reject) => {

      // Lock to guarantee chunks are written in the correct order
      let lk_id = slugid.v4();
      // debug('Lock id:', lk_id);

      srcStream.on('error', err => {
        debug(err);
        reject(err);
      });

      srcStream.on('data', async data => {
        let unlock = await lck(lk_id);
        let fail = false;
        let result;
        try{

          if (typeof data === 'string') {
            data = Buffer.from(data);
          }

          let chunks = _.chunk(data.toJSON().data, CHUNKSIZE);

          for (let i in chunks) {
            let ch = Buffer.from(chunks[i]);
            let cmd = Command.putfile(dest, ch);
            // debug(cmd);
            result = await self._writeAndResolve(cmd);
            // debug(result);
            if (result.error != '') {
              fail = true;
              break;
            }
          }

        }finally{
          unlock();
          if (fail) {
            reject(new Error(result.error));
          }
        }
      });

      srcStream.on('end', resolve);
    });

  }

  * _getfile (src) {

    let cmd = Command.getfile(src); 
    let outStream = through2(function (chunk, enc, cb) { this.push(chunk); cb(); });
    let self = this;
    
    yield outStream;

    yield new Promise((resolve, reject) => {
      self.stdin.write(cmd);
      self._cb[cmd.id] = (data) => {
        // debug('Getfile :', data);
        if (data.error) {
          outStream.end();
          return reject(new Error(data.error));
        }
        outStream.write(data.fileData.data);
        if (data.fileData.totalPieces == data.fileData.currentPiece) {
          outStream.end();
          delete self._cb[cmd.id];
          return resolve(true);
        }
      }
    });

  }

  kill () {
    // this.stdin.destroy();
    this.stdin.end();
    this.stdout.end();
  }
 
  // Wrapper method for putfile
  async writeToFile (remotePath , data) {
    if (typeof data !== 'string' && typeof data !== 'Buffer') {
      throw new Error('\'data\' must be of type string or Buffer');
    }
    let gen = this._putfile(remotePath);
    let stream = gen.next().value;
    let writer = gen.next().value;
    stream.write(data);
    stream.end();
    return await writer;
  }

  // Wrapper methods for _getfile
  async readFileAsString (remotePath, enc = 'utf8') {
    let gen = this._getfile(remotePath);
    let stream = gen.next().value;
    let str = "";
    stream.on('data', data => {
      if(typeof data === 'Buffer') {
        data = data.toString(enc);
      }
      str += data;
    });
    let reader = gen.next().value;
    await reader;
    return str;
  }

  async readFileAsBuffer (){
    let gen = this._getfile(remotePath);
    let stream = gen.next().value;
    let buff = Buffer.alloc(0); // Initialize a zero length buffer
    stream.on('data', data => {
      if(typeof data === 'string') {
        data = Buffer.from(data);
      }
      buff = Buffer.concat([buff, data]);
    });
    let reader = gen.next().value();
    await reader;
    return str;
  }

}

module.exports = FileBrowser;
