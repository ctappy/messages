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
	fg, err := os.OpenFile(path.Join(logdir, "general.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fe, err := os.OpenFile(path.Join(logdir, "email.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	Log.General.Println("Starting...")
	Log.Email.Println("Starting...")
}

// SetupLogs log level defaults
func SetupLogLevel(logLevel string) {
	if logLevel == "fatal" {
		LogSettings.fatalHandle = os.Stdout
	} else if logLevel == "error" {
		LogSettings.fatalHandle = os.Stdout
		LogSettings.errorHandle = os.Stdout
	} else if logLevel == "warn" {
		LogSettings.fatalHandle = os.Stdout
		LogSettings.errorHandle = os.Stdout
		LogSettings.warnHandle = os.Stdout
	} else if logLevel == "info" {
		LogSettings.fatalHandle = os.Stdout
		LogSettings.errorHandle = os.Stdout
		LogSettings.infoHandle = os.Stdout
		LogSettings.debugHandle = os.Stdout
	} else if logLevel == "debug" {
		LogSettings.fatalHandle = os.Stdout
		LogSettings.errorHandle = os.Stdout
		LogSettings.infoHandle = os.Stdout
		LogSettings.warnHandle = os.Stdout
		LogSettings.debugHandle = os.Stdout
	} else if logLevel == "trace" {
		LogSettings.fatalHandle = os.Stdout
		LogSettings.errorHandle = os.Stdout
		LogSettings.traceHandle = os.Stdout
		LogSettings.infoHandle = os.Stdout
		LogSettings.warnHandle = os.Stdout
		LogSettings.debugHandle = os.Stdout
	} else {
		log.Println("Please use one of the following options for log level")
		log.Println("fatal, error, warn, info, debug, or trace")
		log.Fatalln("Exiting")
	}
}
