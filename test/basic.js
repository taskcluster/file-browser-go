/* Check if it works */
const
  child_process = require('child_process'),
	debug					=	require('debug')('test'),
  FileBrowser   = require('../lib/FileBrowser'),
  assert        = require('assert'),
  fs            = require('fs'),
  _             = require('lodash'),
  path          = require('path'),
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
    let result = await browser.ls(path.join(TEST_HOME, 'ls'));
    assert(result.files.length === 4);
  });

  it('can remove a directory', async function() {

    assert(fs.existsSync(path.join(TEST_HOME,'/rmdir')));
    assert(fs.existsSync(path.join(TEST_HOME,'/rmfile')));

    let result = await browser.rm(path.join(TEST_HOME,'/rmdir'));
    
    result = await browser.rm(path.join(TEST_HOME,'/rmfile'));

    assert(!fs.existsSync(path.join(TEST_HOME,'/rmdir')));
    assert(!fs.existsSync(path.join(TEST_HOME,'/rmfile')));
  });

  it('can create a directory', async function() {
    assert(!fs.existsSync(path.join(TEST_HOME,'/mkdir')));
    let result = await browser.mkdir(path.join(TEST_HOME,'/mkdir'));
    assert(fs.existsSync(path.join(TEST_HOME,'/mkdir')));
  });

  it('can copy a directory', async function() {
    let result = await browser.cp(path.join(TEST_HOME,'/cpsrc'), 
      path.join(TEST_HOME,'/cpdest'));
    assert(fs.existsSync(path.join(TEST_HOME,'/cpdest/cpsrc')));
  });

  it('can move a directory', async function() {
    assert(fs.existsSync(path.join(TEST_HOME,'/mvdir')));
    assert(fs.existsSync(path.join(TEST_HOME,'/mvfile')));

    let result = await browser.mv(path.join(TEST_HOME,'/mvfile'), path.join(TEST_HOME,'/mvdir/mvfile'));
    debug(result);

    assert(fs.existsSync(path.join(TEST_HOME,'/mvdir/mvfile')));

    result = await browser.mv(path.join(TEST_HOME,'/mvdir'), path.join(TEST_HOME,'/mvdest'));

    assert(fs.existsSync(path.join(TEST_HOME,'/mvdest')));
  });

  it('_getfile test', async function () {
    const fileName = path.join(TEST_HOME,'/ls/_getfileTest.txt');
    const destFile = path.join(TEST_HOME,'/_getfileTest.txt');
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

  });

  it('_putfile test', async function () {
    const fileName = path.join(TEST_HOME,'/putFileTest.txt');
    const dest = path.join(TEST_HOME,'/ls/putFileTest.txt');
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
  });

  it('test readFileAsString', async function() {
    const fileName = path.join(TEST_HOME,'/readFileAsString.txt');
    fs.appendFileSync(fileName, 'Hello');
    let str = await browser.readFileAsString(fileName);
    debug(str);
    assert(str == 'Hello');
  });

  it('test writeToFile', async function() {
    const fileName = path.join(TEST_HOME,'/writeToFile.txt');
    await browser.writeToFile(fileName, 'Hello');
    let str = fs.readFileSync(fileName).toString();
    debug(str);
    assert(str === 'Hello');

  });

});
