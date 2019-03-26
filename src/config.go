package configuration

func Load(jsonFile io.Reader) Config {
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err := json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatalf("Failed to load json file %v", err)
	}
	return config
}
