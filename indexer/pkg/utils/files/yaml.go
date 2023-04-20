package files

import (
	v1 "github.com/ghodss/yaml"
	"io/ioutil"
)

type Yaml struct {
	FilePath string
}

func NewYaml(filePath string) *Yaml {
	return &Yaml{FilePath: filePath}
}

func (y *Yaml) Read(data interface{}) error {
	fileContent, err := ioutil.ReadFile(y.FilePath)
	if err != nil {
		return err
	}
	return v1.Unmarshal(fileContent, &data)
}

func (y *Yaml) Write(data interface{}) error {
	marshal, err := v1.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(y.FilePath, marshal, 0777)
}
