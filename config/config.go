package config

import "encoding/json"
import "os"

const CONFIG_FILE_PATH = "./config/config.json"

// Exact structure of the config.json file
type configFile struct {
	// Databases map[string]interface{} `json:"databases"`
	Databases       DatabasesConfig
	TokenExpiration int `json:"jwtExpirationOffset"`
}

type DatabasesConfig struct {
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

// Returns map of current config file
// Can return an error if:
// 	 - the config file could not be read
// 	 - the config file could not be decoded
func GetConfig() (cf configFile, err error) {
	var ret configFile
	var f *os.File
	f, err = os.Open(CONFIG_FILE_PATH)
	if err != nil {
		return
	}
	dec := json.NewDecoder(f)
	if err = dec.Decode(&ret); err != nil {
		return
	}
	return ret, nil
}
