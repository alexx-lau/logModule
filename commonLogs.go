package logModule

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenLogFileDaily(path, filename string, maxFileLength int64, fixedLength int64) (file *os.File, err error) {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	today := time.Now().Format("02")
	_, err = os.OpenFile(path + "/" + year + "/" + month + "/" + today, os.O_RDONLY, 0666)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path + "/" + year + "/" + month + "/" + today, 0666)
		if err != nil {
			return
		}
	}
	path += "/" + year + "/" + month + "/" + today + "/"
	file, err = os.OpenFile(path + filename, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0666)
	}
	if err != nil {
		return
	}
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return
	}
	if fileInfo.Size() + int64(fixedLength) > int64(maxFileLength) {
		temp := strings.Split(filename, ".")
		idx := 1
		filename = temp[0] + "-" + strconv.Itoa(idx)
		var oFile *os.File
	BREAKPOINT:
		for {
			oFile, err = os.OpenFile(path + filename + "." + temp[1], os.O_RDONLY, 0666)
			if os.IsNotExist(err) {
				err = oFile.Close()
				if err != nil {
					return
				}
				oFile, err = os.OpenFile(path + filename + "." + temp[1], os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
				if err != nil {
					return
				}
				break BREAKPOINT
			} else {
				var oFileInfo os.FileInfo
				oFileInfo, err = oFile.Stat()
				if err != nil {
					return
				}
				if oFileInfo.Size() + int64(fixedLength) > int64(maxFileLength) {
					err = oFile.Close()
					if err != nil {
						return
					}
					idx += 1
					filename = temp[0] + "-" + strconv.Itoa(idx)
				}
			}
		}
		var buffer [1024]byte
		var length int
	WRITE_BREAKPOINT:
		for {
			length, err = file.Read(buffer[:])
			if err != nil {
				if io.EOF == err {
					length, err = oFile.Write(buffer[:length])
					if err != nil {
						return
					}
					err = oFile.Sync()
					if err != nil {
						return
					}
					break WRITE_BREAKPOINT
				} else {
					return
				}
			}
			length, err = oFile.Write(buffer[:length])
			if err != nil {
				return
			}
			err = oFile.Sync()
			if err != nil {
				return
			}
		}
		err = file.Truncate(0)
		if err != nil {
			return
		}
		err = oFile.Close()
	}
	return
}
