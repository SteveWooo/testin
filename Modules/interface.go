package modules

type ConsensusStorageData interface {
	LoadFromJSONString(jsonString string) error
	LoadFromInterface(obj interface{}) error
}
