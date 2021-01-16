/*
   Copyright (C) 2021 Dennis Jung

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Provisioners       []Provisioner `json:"provisioners"`
	DefaultProvisioner bool          `json:"defaultProvisioner"`
}

type Provisioner struct {
	// The name of the provisioner
	Name string `json:"name"`
	// The directory to create the persistant volume directories in.
	Directory string `json:"directory"`
	// Override the default reclaim-policy of dynamicly provisioned volumes (The default is remove).
	ReclaimPolicy string `json:"reclaimPolicy"`
}

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
