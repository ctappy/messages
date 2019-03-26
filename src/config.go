package configuration

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func load(jsonFile io.Reader) Config {
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err := json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatalf("Failed to load json file %v", err)
	}
	return config
}

func LoadConfig() Config {
	// if we crash the go code, output file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// setup flags
	configPtr := flag.String("config", "./config.json", "JSON config file location")
	debug := flag.Bool("debug", false, "debug option")
	flag.Parse()

	// load json
	if _, err := os.Stat(*configPtr); err == nil {
		if *debug {
			log.Printf("Loading configuration from %q\n", *configPtr)
		}
	} else if os.IsNotExist(err) {
		log.Fatalf("File not found %q %v\n", *configPtr, err)
	} else {
		log.Fatalf("Issue finding file %q %v\n", *configPtr, err)
	}
	jsonFile, err := os.Open(*configPtr)
	if err != nil {
		log.Fatalf("Failed to open %q %v", *configPtr, err)
	}
	defer jsonFile.Close()
	if *debug {
		log.Printf("Successfully Opened %q\n", *configPtr)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	return load(jsonFile)
}
