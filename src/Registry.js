'use strict';

const
  EventEmitter    = require('events').EventEmitter,
  // StringDecoder   = require('string_decoder').StringDecoder,
  // _               = require('lodash'),
	debug						=	require('debug')('registry'),
	msgpack					=	require('msgpack'),
  // FileOperations  = require('./FileOperations');

// const decoder = new StringDecoder();

class Registry extends EventEmitter {

  constructor(outputStream){
    super();

    this.outputStream = outputStream;

		/*
    this.curBuff = "";
    this.pairCount = 0;
    this.curBuffIndex = 0;

    this.outputStream.on('data', data => {
      this.curBuff += decoder.write(data);

      while(this.curBuffIndex < this.curBuff.length){
        if(this.curBuff[this.curBuffIndex] === '{') 
          this.pairCount++;
        else if(this.curBuff[this.curBuffIndex] === '}') 
          this.pairCount--;

        this.curBuffIndex++;

        if(this.pairCount === 0 && this.curBuff[this.curBuffIndex - 1] === '}'){
          var i = 0;
          var residue = "";
          var tempBuff = "";

          while(i < this.curBuffIndex){
            tempBuff += this.curBuff[i]; 
            i++;
          }
          while(i < this.curBuff.length){
            residue += this.curBuff[i];
            i++;
          }

          this.curBuff = residue;
          this.curBuffIndex = 0;

          this.processString(tempBuff);
        }
      }

    });

    this.outputStream.on('end', () => {
      if(this.curBuff === "") return;
      this.processString(this.curBuff);
    });
		*/

		let msgStream = new msgpack.Stream(this.outputStream);
		msgStream.addListener('msg', msg => {
			debug(msg);	
			this.process(msg);
		});

  }

  process(msg){
    if(!msg.id || !msg.cmd) {
      this.emit('error', msg);
      return;
    }

    // Not the last chunk returned by getfile
    this.emit(msg.cmd, msg);
    return;
  }

}

module.exports = Registry;