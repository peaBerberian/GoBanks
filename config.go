package main

import "encoding/json"
import "os"

const config_file_path = "./config/config.json"

// Exact structure of the config.json file
type configFile struct {
	// Databases map[string]interface{} `json:"databases"`
	Database        interface{} `json:"database"`
	TokenExpiration int         `json:"jwtExpirationOffset"`
	ServerPort      int         `json:"port"`
}

// getConfig parse the config file. See config_file_path.
// Can return an error if:
// 	 - the config file could not be read
// 	 - the config file could not be decoded
func getConfig() (cf configFile, err error) {
	var ret configFile
	var f *os.File
	f, err = os.Open(config_file_path)
	if err != nil {
		return
	}
	dec := json.NewDecoder(f)
	if err = dec.Decode(&ret); err != nil {
		return
	}
	return ret, nil
}
