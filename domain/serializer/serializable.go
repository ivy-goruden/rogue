package serializer

type Serializable interface {
	ToMap() map[string]interface{}
	FromMap(map[string]interface{}) error
}
