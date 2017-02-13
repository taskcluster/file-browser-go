mkdir test_dir

export TEST_HOME=`pwd`/test_dir
echo "Running tests in $TEST_HOME"

make_ls () {
	mkdir -p $TEST_HOME/ls/dir
	mkdir -p $TEST_HOME/ls/.hidden_dir
	touch $TEST_HOME/ls/file
	touch $TEST_HOME/ls/.hidden_file
}

make_rm () {
	mkdir -p $TEST_HOME/rmdir
	touch $TEST_HOME/rmfile
}

make_mv () {
	mkdir -p $TEST_HOME/mvdir
	touch $TEST_HOME/mvfile
}

make_cp () {
	mkdir -p $TEST_HOME/cpsrc/dir
	mkdir -p $TEST_HOME/cpsrc/dir2
	touch $TEST_HOME/cpsrc/dir2/cpfile
	
	mkdir -p $TEST_HOME/cpdest
}

make_ls
make_rm
make_mv
make_cp
DEBUG=browser,test,registry npm test
rm -rf test_dir
