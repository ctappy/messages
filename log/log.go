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
	fatalHandle   io.Writer
	errorHandle   io.Writer
	debugHandle   io.Writer
	warnHandle    io.Writer
	infoHandle    io.Writer
	traceHandle   io.Writer
	generalHandle io.Writer
	emailHandle   io.Writer
}

var (
	Log         = Logs{}
	LogSettings = &logType{
		fatalHandle:   ioutil.Discard,
		errorHandle:   ioutil.Discard,
		debugHandle:   ioutil.Discard,
		warnHandle:    ioutil.Discard,
		infoHandle:    ioutil.Discard,
		traceHandle:   ioutil.Discard,
		generalHandle: ioutil.Discard,
		emailHandle:   ioutil.Discard,
	}
)

func LogInit(logType logType) {
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

	Log.General = log.New(logType.generalHandle,
		"GENERAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.Email = log.New(logType.emailHandle,
		"EMAIL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

// SetupLog log level defaults and log files
func SetupLog(logdir, logLevel string) {
	// Setup log files

	/////////////////////
	// STDOUT LOG FILE //
	/////////////////////
	f, err := os.OpenFile(path.Join(logdir, "stdout.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	/////////////////////
	// STDERR LOG FILE //
	/////////////////////
	fe, err := os.OpenFile(path.Join(logdir, "stderr.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close() // cannot close inside function
	// defer fe.Close() // cannot close inside function

	////////////////////////////////
	// GENERAL AND EMAIL LOG FILE //
	////////////////////////////////
	fg, err := os.OpenFile(path.Join(logdir, "general.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	femail, err := os.OpenFile(path.Join(logdir, "email.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer fe.Close() // cannot close inside function
	// defer fg.Close() // cannot close inside function

	//////////////////////
	// Setup log levels //
	//////////////////////
	LogSettings.generalHandle = io.MultiWriter(os.Stderr, fg)
	LogSettings.emailHandle = io.MultiWriter(os.Stderr, femail)

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
