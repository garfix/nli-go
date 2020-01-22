package global

import (
	"encoding/json"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"os"
)

type DialogContextFileStorage struct {
	log *common.SystemLog
}

func NewDialogContextFileStorage(log *common.SystemLog) *DialogContextFileStorage {
	return &DialogContextFileStorage{ log: log }
}

func (storage DialogContextFileStorage) Read(dialogContextPath string, dialogContext *central.DialogContext) {

	if dialogContextPath == "" {
		return
	}

	_, err := os.Stat(dialogContextPath)
	if os.IsNotExist(err) {
		// session file does not exist yet; it will be created when the session ends
		return
	}

	dialogContextJson, err := common.ReadFile(dialogContextPath)
	if err != nil {
		storage.log.AddError(err.Error())
		return
	}

	err = json.Unmarshal([]byte(dialogContextJson), &dialogContext)
	if err != nil {
		storage.log.AddError("Error parsing JSON file " + dialogContextJson + " (" + err.Error() + ")")
		return
	}
}

func (storage DialogContextFileStorage) Write(dialogContextPath string, dialogContext *central.DialogContext) {

	if dialogContextPath == "" {
		return
	}

	jsonBytes, err := json.Marshal(dialogContext)
	if err != nil {
		storage.log.AddError("Error serializing dialog context (" + err.Error() + ")")
		return
	}

	jsonString := string(jsonBytes)

	err = common.WriteFile(dialogContextPath, jsonString)
	if err != nil {
		storage.log.AddError("Error writing dialog context file " + dialogContextPath + " (" + err.Error() + ")")
		return
	}
}
