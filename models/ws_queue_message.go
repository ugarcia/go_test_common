package models

type WsQueueMessage struct {
    BaseMessage
    Code            string                 `json:"code"`
}
