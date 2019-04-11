package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"openpitrix.io/openpitrix/pkg/logger"
)

type AnyMap map[interface{}]interface{}

//return: content contentMap error
func ReadYamlFile(f string) ([]byte, AnyMap, error) {
	//read global_config file and convert to map
	content, err := ioutil.ReadFile(f)
	if err != nil {
		logger.Error(nil, "Failed to read file %s!", f)
		return nil, nil, err
	}

	contentMap := make(AnyMap)
	err = yaml.Unmarshal(content, contentMap)
	if err != nil {
		logger.Error(nil, "Failed to Unmarshal yaml to map!")
	}
	return content, contentMap, err
}

func UpdateMap(m AnyMap, updateM AnyMap) {
	for k, _ := range m {
		if updateM[k] != nil {
		}
	}
}
