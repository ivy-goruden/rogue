package serializer

import (
	"encoding/json"
	"fmt"
)

type Serializer struct{}

// создает новый сериализатор
func MakeSerializer() *Serializer {
	return &Serializer{}
}

// Serialize сериализует любую структуру, реализующую Serializable
func (s *Serializer) Serialize(obj Serializable) ([]byte, error) {
	data := obj.ToMap()
	return json.Marshal(data)
}

// Deserialize десериализует данные в объект
func (s *Serializer) Deserialize(data []byte, obj Serializable) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return obj.FromMap(raw)
}
