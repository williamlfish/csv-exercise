package main

import (
	"flag"
	"os"
	"strconv"

	parser "github.com/scoir/csv-exercise/pkg/scoirparser"

	log "github.com/sirupsen/logrus"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
func FatalLog(err error, message string) {
	log.WithFields(log.Fields{
		"error": err,
	}).Fatal(message)
}

func ErrWarnLog(err error, message string) {
	log.WithFields(log.Fields{
		"error": err,
	}).Warn(message)
}
func SuccessLog(file string, message string) {
	log.WithFields(log.Fields{
		"file": file,
	}).Info(message)
}

func checkMakeDirs(dirList []string) {
	for _, dir := range dirList {
		logInfo := log.Fields{
			"dir": dir,
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.WithFields(logInfo).Warn("dir not found, making it")
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				FatalLog(err, "Can't make requested dirs, dying now..")
			}
		} else {
			log.WithFields(logInfo).Info("found dir")
		}
	}
}

func main() {
	var dirList []string
	inEnv := getEnv("INPUT_DIR", "./input-dir")
	outEnv := getEnv("OUTPUT_DIR", "./out-dir")
	errorEnv := getEnv("ERROR_DIR", "./error-dir")
	// default to the envs or overwrite with the flag
	inDir := flag.String("indir", inEnv, "the input dir or watch dir, for users to put their csv's in.")
	outDir := flag.String("outdir", outEnv, "the output dir, for for parsed user files")
	errorDir := flag.String("errordir", errorEnv, "the error dir, for csv files of errors caused the input file")
	uflag := flag.Bool("u", false, "pass this flag in to force parsing only unique file names")
	flag.Parse()

	dirList = append(dirList, *inDir, *outDir, *errorDir)

	dirwatch := DirWatch{
		inDir:    *inDir,
		outDir:   *outDir,
		errorDir: *errorDir,
		unique:   *uflag,
	}
	if dirwatch.unique {
		log.Info("unique file is ON building db connection")
		dbHost := getEnv("PGHOST", "localhost")
		dbPort := getEnv("PGPORT", "5432")
		port, err := strconv.Atoi(dbPort)
		if err != nil {
			ErrWarnLog(err, "env PGPORT does not appere to be an in, defaulting to 5432")
			port = 5432
		}
		dbDb := getEnv("PGDATABASE", "processed_files")
		dbUser := getEnv("PGUSER", "root")
		dbPw := getEnv("PGPASSWORD", "password")
		dbCon, err := parser.NewDBClient(dbUser, dbPw, dbHost, dbDb, port)
		if err != nil {
			FatalLog(err, "could not connect to be, need db for only unique files")
		}
		defer dbCon.DB.Close()
		dirwatch.Db = *dbCon
	}
	checkMakeDirs(dirList)
	dirwatch.WatchInputDir()
}
