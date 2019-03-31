package core

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type logs struct {
	fatal *log.Logger // fatal to application forcing shutdown
	err   *log.Logger // error to operation, but continues application
	warn  *log.Logger // incorrect operation that was recoverable
	info  *log.Logger // general information
	debug *log.Logger // information related to operations
	trace *log.Logger // tracing code
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
	Log         = logs{}
	logSettings = &logType{
		fatalHandle: ioutil.Discard,
		errorHandle: ioutil.Discard,
		debugHandle: ioutil.Discard,
		warnHandle:  ioutil.Discard,
		infoHandle:  ioutil.Discard,
		traceHandle: ioutil.Discard,
	}
)

func logInit(logType logType) {
	Log.trace = log.New(logType.traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.debug = log.New(logType.debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.info = log.New(logType.infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.warn = log.New(logType.warnHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.err = log.New(logType.errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Log.fatal = log.New(logType.fatalHandle,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// setLogLevel log level defaults
func setLogLevel(logLevel string) {
	if logLevel == "fatal" {
		logSettings.fatalHandle = os.Stdout
	} else if logLevel == "error" {
		logSettings.fatalHandle = os.Stdout
		logSettings.errorHandle = os.Stdout
	} else if logLevel == "warn" {
		logSettings.fatalHandle = os.Stdout
		logSettings.errorHandle = os.Stdout
		logSettings.warnHandle = os.Stdout
	} else if logLevel == "info" {
		logSettings.fatalHandle = os.Stdout
		logSettings.errorHandle = os.Stdout
		logSettings.infoHandle = os.Stdout
		logSettings.debugHandle = os.Stdout
	} else if logLevel == "debug" {
		logSettings.fatalHandle = os.Stdout
		logSettings.errorHandle = os.Stdout
		logSettings.infoHandle = os.Stdout
		logSettings.warnHandle = os.Stdout
		logSettings.debugHandle = os.Stdout
	} else if logLevel == "trace" {
		logSettings.fatalHandle = os.Stdout
		logSettings.errorHandle = os.Stdout
		logSettings.traceHandle = os.Stdout
		logSettings.infoHandle = os.Stdout
		logSettings.warnHandle = os.Stdout
		logSettings.debugHandle = os.Stdout
	} else {
		log.Println("Please use one of the following options for log level")
		log.Println("fatal, error, warn, info, debug, or trace")
		log.Fatalln("Exiting")
	}
}
