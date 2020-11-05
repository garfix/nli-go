package global

import (
	"encoding/json"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"os"
)

type DialogContextFileStorage struct {
	cacheDir string
	log      *common.SystemLog
}

func NewDialogContextFileStorage(varDir string, log *common.SystemLog) *DialogContextFileStorage {
	return &DialogContextFileStorage{
		cacheDir: varDir,
		log:      log,
	}
}

func (storage DialogContextFileStorage) Read(sessionId string, dialogContext *central.DialogContext, clearWhenCorrupt bool) {

	if sessionId == "" {
		return
	}

	sessionPath := storage.cacheDir + "/" + sessionId + ".json"

	_, err := os.Stat(sessionPath)
	if os.IsNotExist(err) {
		// session file does not exist yet; it will be created when the session ends
		return
	}

	dialogContextJson, err := common.ReadFile(sessionPath)
	if err != nil {
		storage.log.AddError(err.Error())
		return
	}

	err = json.Unmarshal([]byte(dialogContextJson), &dialogContext)
	if err != nil {
		if !clearWhenCorrupt {
			storage.log.AddError("Error parsing YAML file " + dialogContextJson + " (" + err.Error() + ")")
			return
		}
	}
}

func (storage DialogContextFileStorage) Write(sessionId string, dialogContext *central.DialogContext) {

	if sessionId == "" {
		return
	}

	sessionPath := storage.cacheDir + "/" + sessionId + ".json"

	jsonBytes, err := json.Marshal(dialogContext)
	if err != nil {
		storage.log.AddError("Error serializing dialog context (" + err.Error() + ")")
		return
	}

	jsonString := string(jsonBytes)

	_, err = os.Stat(storage.cacheDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(storage.cacheDir, 0777)
		if err != nil {
			storage.log.AddError("Error creating cache dir " + storage.cacheDir + " (" + err.Error() + ")")
			return
		}
	}


	err = common.WriteFile(sessionPath, jsonString)
	if err != nil {
		storage.log.AddError("Error writing dialog context file " + sessionPath + " (" + err.Error() + ")")
		return
	}
}

func (storage DialogContextFileStorage) Remove(sessionId string) {

	if sessionId == "" {
		return
	}

	sessionPath := storage.cacheDir + "/" + sessionId + ".json"

	_, err := os.Stat(sessionPath)
	if err == nil {
		err = os.Remove(sessionPath)
		if err != nil {
			storage.log.AddError("Error removing session file " + sessionPath + " (" + err.Error() + ")")
			return
		}
	}
}
