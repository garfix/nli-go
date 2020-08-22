package global

import (
	"gopkg.in/yaml.v2"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"path/filepath"
)

type systemBuilder struct {
	log     *common.SystemLog
	baseDir string
	parser  *importer.InternalGrammarParser
}

func NewSystem(systemPath string, log *common.SystemLog) *system {

	system := &system{ log: log }

	absolutePath, err := filepath.Abs(systemPath)
	if err != nil {
		log.AddError(err.Error())
		return system
	}

	builder := NewSystemBuilder(absolutePath, log)
	builder.build(system)

	return system
}

//func readConfig(configPath string, log *common.SystemLog) (systemConfig) {
//
//	config := systemConfig{}
//
//	configJson, err := common.ReadFile(configPath)
//	if err != nil {
//		log.AddError("Error reading JSON file " + configPath + " (" + err.Error() + ")")
//	}
//
//	if log.IsOk() {
//		err := json.Unmarshal([]byte(configJson), &config)
//		if err != nil {
//			log.AddError("Error parsing JSON file " + configPath + " (" + err.Error() + ")")
//		}
//	}
//
//	if config.ParentConfig != "" {
//		parentConfigPath := config.ParentConfig
//		if len(parentConfigPath) > 0 && parentConfigPath[0] != os.PathSeparator {
//			parentConfigPath = filepath.Dir(configPath) + string(os.PathSeparator) + parentConfigPath
//		}
//		parentConfig := system.ReadConfig(parentConfigPath, log)
//
//		config = parentConfig.Merge(config)
//		config.ParentConfig = ""
//	}
//
//	return config
//}

func NewSystemBuilder(baseDir string, log *common.SystemLog) systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return systemBuilder {
		baseDir: baseDir,
		parser:  parser,
		log:     log,
	}
}

func (builder systemBuilder) build(system *system) {
	configPath := builder.baseDir + "/config.yml"

	configYml, err := common.ReadFile(configPath)
	if err != nil {
		builder.log.AddError(err.Error())
	}

	config := config{}

	err = yaml.Unmarshal([]byte(configYml), &config)
	if err != nil {
		builder.log.AddError("Error parsing YAML file " + configPath + " (" + err.Error() + ")")
	}
}
