const
  child_process = require('child_process'),
  FileBrowser   = require('./FileBrowser');

var run = async () => {

  try{
    let shell = child_process.spawn('./file-browser-go', [] , {
      stdio: ['pipe','pipe', 'ignore'],
      detached: true
    });
    let fb = new FileBrowser(shell);
    let result;
		ls(fb, '/');
		ls(fb, '/Users/chinmaykousik/');
		await fb.kill();
  }catch(e){
    console.error(e);
  }
	return;

}

var ls = async (fileBrowser, path) => {
	let result = await fileBrowser.ls(path);
	console.log(result);
}

run();
