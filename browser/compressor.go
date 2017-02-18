package browser

import (
	"compress/gzip"
	"log"
	"os"

	"gopkg.in/vmihailenco/msgpack.v2"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func EncodeCompressWrite(outChan <-chan *ResultSet, outFile *os.File) {
	zipWriter := gzip.NewWriter(outFile)
	defer zipWriter.Close()
	for {
		res := <-outChan
		msg, err := msgpack.Marshal(res)
		handleErr(err)
		zipWriter.Write(msg)
		zipWriter.Flush()
	}
}

func DecompressDecode(inChan chan<- Command, inFile *os.File) {
	var cmd Command
	zipReader, err := gzip.NewReader(inFile)
	defer zipReader.Close()
	handleErr(err)
	for {
		b := make([]byte, 5120)
		n, err := zipReader.Read(b)
		handleErr(err)
		b = b[:n]
		log.New(os.Stderr, "raw: ", 0).Println(b)
		err = msgpack.Unmarshal(b, &cmd)
		handleErr(err)
		log.New(os.Stderr, "command: ", 0).Println(cmd)
		inChan <- cmd
	}
}
