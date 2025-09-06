package models

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WriteRequest struct {
	Items []Item `json:"items"`
}

type ReadResponse struct {
	Items []ReadResponseItem `json:"items"`
}

type ReadResponseItem struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
	Exists bool   `json:"exists"`
}

type KafkaMessage struct {
	Type  string `json:"type"` // "update" or "delete"
	Key   string `json:"key"`
	Value string `json:"value"` // empty for delete operations
}
