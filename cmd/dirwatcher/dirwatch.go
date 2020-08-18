package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	parser "github.com/scoir/csv-exercise/pkg/scoirparser"
	log "github.com/sirupsen/logrus"
)

// DirWatch main struct config for watching the input dir and writting other files out
type DirWatch struct {
	inDir    string
	outDir   string
	errorDir string
	unique   bool
	Db       parser.DBClient
}

// ExtractFileName will get the input filename from the path, as well as confirm it is a csv
func (d DirWatch) ExtractFileName(filePath string) *string {
	fi, err := os.Stat(filePath)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("there was an error in the input dir %s", d.inDir))
		return nil
	}
	if fi.Mode().IsDir() {
		log.WithFields(log.Fields{
			"dir": filePath,
		}).Warn("new dir created, will not recusivly find files and assume this is user clean up.")
		return nil
	}
	pathSlice := strings.Split(filePath, "/")
	fileName := pathSlice[1]
	ex := path.Ext(fileName)
	if ex != ".csv" {
		log.WithFields(log.Fields{
			"file":      fileName,
			"input-dir": d.inDir,
		}).Warn("found new file in input dir, but not csv, ignoring")
		return nil
	}
	return &fileName
}

// WatchInputDir watches the input dir for create events
func (d DirWatch) WatchInputDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		FatalLog(err, "error setting up watcher")
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fileName := d.ExtractFileName(event.Name)
					if fileName != nil {
						d.ProcessFile(*fileName)
					}
				}

			case err := <-watcher.Errors:
				FatalLog(err, fmt.Sprintf("error while watching dir %s", d.inDir))
			}
		}
	}()

	if err := watcher.Add(d.inDir); err != nil {
		FatalLog(err, fmt.Sprintf("cannot watch the requested dir %s", d.inDir))
	}

	<-done
}

//ProcessFile will process a new file added to the input dir
func (d DirWatch) ProcessFile(fileName string) {
	if d.unique {
		hasProcessed, err := d.Db.CheckProcessedFile(fileName)
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error quering for file %s", fileName))
			return
		}
		if hasProcessed {
			ErrWarnLog(err, fmt.Sprintf("unique file ON, previosly processed file: %s", fileName))
			return
		}
	}
	data, err := d.GetFile(fileName)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error getting file data %s", fileName))
		return
	}
	jsonData, errList, err := parser.ParseCsvToJsonBytes(data)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error parsing csv into json file:%s", fileName))
		return
	}
	go d.WriteJsonFile(jsonData, fileName)
	if len(errList) > 0 {
		go d.WriteErrorFile(errList, fileName)
	}
}

// GetFile gets the file data from the input dir
func (d DirWatch) GetFile(fileName string) (io.Reader, error) {
	path := fmt.Sprintf("%s/%s", d.inDir, fileName)
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	os.Remove(path)
	return data, nil
}

//WriteJsonFile will write the json file with a same name as the csv file
func (d DirWatch) WriteJsonFile(data []byte, filename string) {
	noExName := strings.Split(filename, ".")[0]
	jsonName := fmt.Sprintf("%s/%s.json", d.outDir, noExName)
	_, err := os.Stat(jsonName)
	if os.IsNotExist(err) {
		f, err := os.Create(jsonName)
		defer f.Close()
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error while making new file: %s", filename))
			return
		}
		_, err = f.Write(data)
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error while making new file: %s", filename))
			return
		}
		SuccessLog(jsonName, "successfully wrote json file")
		err = d.Db.InsertProcessedFile(filename)
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error recording processed file: %s", filename))
		}
		return
	}
	f, err := os.OpenFile(jsonName, os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error while opening file: %s", filename))
		return
	}
	err = f.Truncate(0)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error cleaing old file out file out: %s", filename))
		return
	}
	_, err = f.WriteAt(data, 0)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error writing file: %s", filename))
		return
	}
	//mark file here just incase there is a file in the dir already...
	err = d.Db.InsertProcessedFile(filename)
	SuccessLog(jsonName, "successfully wrote json file")

}

func buildErrorCsvData(errList []parser.FileError) [][]string {
	var csvData = [][]string{
		{"LINE_NUM", "ERROR_MSG"},
	}
	for _, e := range errList {
		l := strconv.Itoa(e.Line)

		row := []string{
			l,
			e.ErrorMsg,
		}
		csvData = append(csvData, row)
	}
	return csvData
}

// WriteErrorFile writes the error csv file
func (d DirWatch) WriteErrorFile(errList []parser.FileError, filename string) {
	errName := fmt.Sprintf("%s/error-%s", d.errorDir, filename)
	_, err := os.Stat(errName)
	data := buildErrorCsvData(errList)
	if os.IsNotExist(err) {
		f, err := os.Create(errName)
		writer := csv.NewWriter(f)
		defer f.Close()
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error while making new file: %s", errName))
			return
		}
		err = writer.WriteAll(data)
		if err != nil {
			ErrWarnLog(err, fmt.Sprintf("error while making new file: %s", errName))
			return
		}
		SuccessLog(errName, "successfully wrote csv error file")
		return
	}
	f, err := os.OpenFile(errName, os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error while making new file: %s", errName))
		return
	}
	writer := csv.NewWriter(f)
	err = writer.WriteAll(data)
	if err != nil {
		ErrWarnLog(err, fmt.Sprintf("error while writing csv error file: %s", errName))
		return
	}
	SuccessLog(errName, "successfully wrote csv error file")
}
