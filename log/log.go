package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Logs struct {
	// standard log levels
	Fatal *log.Logger // fatal to application forcing shutdown
	Err   *log.Logger // error to operation, but continues application
	Warn  *log.Logger // incorrect operation that was recoverable
	Info  *log.Logger // general information
	Debug *log.Logger // information related to operations
	Trace *log.Logger // tracing code
	// logs to output file and stdout
	General *log.Logger // general logs to output file general.log
	Email   *log.Logger // email logs to output file email.log

}

type logType struct {
	fatalHandle io.Writer
	errorHandle io.Writer
	debugHandle io.Writer
	warnHandle  io.Writer
	infoHandle  io.Writer
	traceHandle io.Writer
}

var (
	Log         = Logs{}
	LogSettings = &logType{
		fatalHandle: ioutil.Discard,
		errorHandle: ioutil.Discard,
		debugHandle: ioutil.Discard,
		warnHandle:  ioutil.Discard,
		infoHandle:  ioutil.Discard,
		traceHandle: ioutil.Discard,
	}
)

func LogInit(logdir string, logType logType) {
	Log.Trace = log.New(logType.traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Debug = log.New(logType.debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Info = log.New(logType.infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Warn = log.New(logType.warnHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Err = log.New(logType.errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Fatal = log.New(logType.fatalHandle,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Setup logs to stdout and files
	fg, err := os.OpenFile(path.Join(logdir, "general.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fe, err := os.OpenFile(path.Join(logdir, "email.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer fe.Close() // cannot close inside function
	// defer fg.Close() // cannot close inside function
	Log.General = log.New(io.MultiWriter(os.Stdout, fg),
		"GENERAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Email = log.New(io.MultiWriter(os.Stdout, fe),
		"EMAIL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// SetupLogs log level defaults
func SetupLogLevel(logdir, logLevel string) {
	f, err := os.OpenFile(path.Join(logdir, "stdout.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fe, err := os.OpenFile(path.Join(logdir, "stderr.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close() // cannot close inside function
	// defer fe.Close() // cannot close inside function
	if logLevel == "fatal" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
	} else if logLevel == "error" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.errorHandle = io.MultiWriter(os.Stderr, fe)
	} else if logLevel == "warn" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.errorHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.warnHandle = io.MultiWriter(os.Stdout, f)
	} else if logLevel == "info" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.errorHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.infoHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.debugHandle = io.MultiWriter(os.Stdout, f)
	} else if logLevel == "debug" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.errorHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.infoHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.warnHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.debugHandle = io.MultiWriter(os.Stdout, f)
	} else if logLevel == "trace" {
		LogSettings.fatalHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.errorHandle = io.MultiWriter(os.Stderr, fe)
		LogSettings.traceHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.infoHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.warnHandle = io.MultiWriter(os.Stdout, f)
		LogSettings.debugHandle = io.MultiWriter(os.Stdout, f)
	} else {
		log.Println("Please use one of the following options for log level")
		log.Println("fatal, error, warn, info, debug, or trace")
		log.Fatalln("Exiting")
	}
}
