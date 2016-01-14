package models

type DbQueueMessage struct {
    BaseMessage
    Entity          string                 `json:"entity"`
}
