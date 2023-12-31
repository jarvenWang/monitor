package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}
type ReadFromFile struct {
	path string //读取文件的路径
}

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取模块
	//打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(err)
	}
	//从文件末尾读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(err)
		}
		//line := "message"
		rc <- line[:len(line)-1]
	}

}

type WriteToInfluxDB struct {
	influxDBDsn string //data source
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块
	for v := range wc {
		fmt.Println(v)
	}

}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

func (l *LogProcess) Process() {
	//解析模块
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}

}

func main() {
	//pi := flag.Int("a", 10, "apple")
	//flag.Parse()
	//fmt.Printf("%v\n", *pi)

	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "user&pwd",
	}
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}
