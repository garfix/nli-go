package common

// Not in use

import (
	"encoding/json"
	"os"
)

type FileStorage struct {
	outputDir   string
	sessionId   string
	storageType string
	storageName string
	log         *SystemLog
}

type StorableObject interface {
}

const StorageNone = "none"
const StorageSession = "session"
const StorageGlobal = "global"

func NewFileStorage(outputDir string, sessionId string, storageType string, storageName string, log *SystemLog) *FileStorage {
	return &FileStorage{
		outputDir:   outputDir,
		sessionId:   sessionId,
		storageType: storageType,
		storageName: storageName,
		log:         log,
	}
}

func (storage FileStorage) Read(object StorableObject) {

	path := storage.getPath()
	if path == "" {
		return
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// session file does not exist yet; it will be created when the session ends
		return
	}

	jsonString, err := ReadFile(path)
	if err != nil {
		storage.log.AddError(err.Error())
		return
	}

	err = json.Unmarshal([]byte(jsonString), object)
	if err != nil {
		return
	}
}

func (storage FileStorage) Write(object StorableObject) {

	path := storage.getPath()
	if path == "" {
		return
	}

	jsonBytes, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		storage.log.AddError("Error serializing object (" + err.Error() + ")")
		return
	}

	jsonString := string(jsonBytes)

	_, err = os.Stat(storage.outputDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(storage.outputDir, 0777)
		if err != nil {
			storage.log.AddError("Error creating dir " + storage.outputDir + " (" + err.Error() + ")")
			return
		}
	}

	err = WriteFile(path, jsonString)
	if err != nil {
		storage.log.AddError("Error writing file " + path + " (" + err.Error() + ")")
		return
	}
}

func (storage FileStorage) Remove(sessionId string) {

	path := storage.getPath()
	if path == "" {
		return
	}

	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			storage.log.AddError("Error removing file " + path + " (" + err.Error() + ")")
			return
		}
	}
}

func (storage FileStorage) getPath() string {
	name := ""

	if storage.storageType == StorageGlobal {
		name = storage.storageName
	} else if storage.storageType == StorageSession {
		if storage.sessionId == "" {
			return ""
		}
		name = storage.storageName + "_" + storage.sessionId
	} else {
		return ""
	}

	return storage.outputDir + "/" + name + ".json"
}
