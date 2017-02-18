/* Check if it works */
const
  child_process = require('child_process'),
	debug					=	require('debug')('test'),
  FileBrowser   = require('../lib/FileBrowser'),
  assert        = require('assert'),
  fs            = require('fs'),
  _             = require('lodash'),
  StringDecoder = require('string_decoder').StringDecoder;

const decoder = new StringDecoder();

let browser, shell;

const TEST_HOME = process.env.TEST_HOME;

describe ('Basic', function(){

  before(function(){
    shell = child_process.spawn ('./file-browser-go', [] , {
      stdio: ['pipe', 'pipe', 'pipe']
    });
    browser = new FileBrowser(shell.stdin, shell.stdout);
    shell.stderr.pipe(process.stderr)
  });

  it('can list contents of a directory', async function() {
    try{
      let result = await browser.ls(TEST_HOME + '/ls');
      assert(result.error.length === 0);
      assert(result.files.length === 4);
      return null;
    }catch(err){
      return err;
    }
  });

  it('can remove a directory', async function() {
    try{

      assert(fs.existsSync(TEST_HOME + '/rmdir'));
      assert(fs.existsSync(TEST_HOME + '/rmfile'));

      let result = await browser.rm(TEST_HOME + '/rmdir');
      assert(result.error.length === 0);
      
      result = await browser.rm(TEST_HOME + '/rmfile');
      assert(result.error.length === 0);

      assert(!fs.existsSync(TEST_HOME + '/rmdir'));
      assert(!fs.existsSync(TEST_HOME + '/rmfile'));
      return null;
      
    }catch(err){

      return err;

    }
  });

  it('can create a directory', async function() {
    try{

      assert(!fs.existsSync(TEST_HOME + '/mkdir'));
      let result = await browser.mkdir(TEST_HOME + '/mkdir');
      assert(result.error.length === 0);
      assert(fs.existsSync(TEST_HOME + '/mkdir'));
      return null;

    }catch(err){
      return err;
    }
  });

  it('can copy a directory', async function() {
    try{

      let result = await browser.cp(TEST_HOME + '/cpsrc', TEST_HOME + '/cpdest');
      assert(fs.existsSync(TEST_HOME + '/cpdest/cpsrc'));
      return null;

    }catch(err){

      return err;
    }
  });

  it('can move a directory', async function() {
    try{
      
      assert(fs.existsSync(TEST_HOME + '/mvdir'));
      assert(fs.existsSync(TEST_HOME + '/mvfile'));

      let result = await browser.mv(TEST_HOME + '/mvfile', TEST_HOME + '/mvdir/mvfile');
      debug(result);
      assert(result.error.length === 0);

      assert(fs.existsSync(TEST_HOME + '/mvdir/mvfile'));

      result = await browser.mv(TEST_HOME + '/mvdir', TEST_HOME + '/mvdest');
      assert(result.error.length === 0);

      assert(fs.existsSync(TEST_HOME + '/mvdest'));

      return null;

    }catch(err){
      return err;
    }
  });

  it('_getfile test', async function () {
    const fileName = TEST_HOME + '/ls/_getfileTest.txt';
    const destFile = TEST_HOME + '/_getfileTest.txt';
    try {
      // Create a new file and append 'Hello'
      fs.appendFileSync(fileName, "Hello");
      let gen = browser._getfile(fileName);
      
      // Create a write stream for the destination file
      let outStream = fs.createWriteStream(destFile);
      let readStream = gen.next().value;
      readStream.pipe(outStream);
      // Run _getfile
      let result = await gen.next().value;
      debug('_getfile result: ',result);
      assert(result != false);
      outStream.close();
      
      //Check if contents are the same
      let str = decoder.write(fs.readFileSync(destFile));
      let target = decoder.write(fs.readFileSync(fileName));
      debug(str, target);
      assert(str == target);
      return null;

    }catch(err) {
      debug(err);
      debug(err);
      return err;
      return err;
    } 
  });

  it('_putfile test', async function () {
    const fileName = TEST_HOME + '/putFileTest.txt';
    const dest = TEST_HOME + '/ls/putFileTest.txt';
    try{
      // Create file 
      fs.appendFileSync(fileName, "Hello");
      
      // Create readable stream for file
      let inStream = fs.createReadStream(fileName);
      let gen = browser._putfile(dest);
      let stream = gen.next().value;
      let writer = gen.next().value;
      inStream.pipe(stream);
      let result = await writer;
      // assert(result.error === "");

      let src = fs.readFileSync(fileName);
      let target = fs.readFileSync(dest);
      debug(src,target);
      assert(src.equals(target));
      return null;
    }catch(err){
      debug(err);
      return err;
    }
  });

  it('test readFileAsString', async function() {
    const fileName = TEST_HOME + '/readFileAsString.txt';
    try{
      fs.appendFileSync(fileName, 'Hello');
      let str = await browser.readFileAsString(fileName);
      debug(str);
      assert(str == 'Hello');
    } catch(err) {
      debug(err);
      return err;
    }
  });

  it('test writeToFile', async function() {
    const fileName = TEST_HOME + '/writeToFile.txt';
    try {
      await browser.writeToFile(fileName, 'Hello');
      let str = fs.readFileSync(fileName).toString();
      debug(str);
      assert(str === 'Hello');
    } catch (err) {
      debug(err);
      return err;
    }
  });

});
