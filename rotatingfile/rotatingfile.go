package rotatingfile

import (
	"fmt"
	"os"
)

/* we'll need the total number of files to keep as the backlog as well as the
* max size of each of them */

type LogFile struct {
	basepath    string
	written     int
	filehandle  *os.File
	maxSize     int
	backlogSize int
}

func Create(filepath string, maxsize int, backlog int) *LogFile {
	var (
		f   *os.File
		err error
	)
	if f, err = os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		panic(err)
	}
	return &LogFile{basepath: filepath, filehandle: f, maxSize: maxsize, backlogSize: backlog}
}

//type Writer interface {
//    Write(p []byte) (n int, err error)
//}

func (lf *LogFile) Write(p []byte) (n int, err error) {
	if len(p)+lf.written > lf.maxSize {
		// rotate the existing logs if any
		// TODO: do this more intelligently?
		var (
			f   *os.File
			err error
		)
		for i := lf.backlogSize - 1; i > 0; i -= 1 {
			oldname := fmt.Sprintf("%s.%d", lf.basepath, i)
			newname := fmt.Sprintf("%s.%d", lf.basepath, i+1)
			_ = os.Rename(oldname, newname)
		}
		lf.filehandle.Close()
		os.Rename(lf.basepath, fmt.Sprintf("%s.1", lf.basepath))
		// reopen the log file handle
		if f, err = os.OpenFile(lf.basepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			panic(err)
		}
		lf.filehandle = f
		lf.written = 0
	}
	lf.written += len(p)
	return lf.filehandle.Write(p)
}
