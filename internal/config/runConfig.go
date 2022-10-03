package config

import "strings"

type RunConfig struct {
	ConfigFile     string
	ConnectionName string
	QueryName      string
	Wrapper        string
	FirstOnly      bool
	KeyName        string
	ValueName      string
	Output         string
	Format         string
	JsonKeys       string
}

func (config *RunConfig) GetOutputFileName() string {
	fileName := config.Output
	ext := "." + config.Format

	if !strings.HasSuffix(strings.ToLower(fileName), ext) {
		fileName = fileName + ext
	}

	return fileName
}
