package models

type QueueMessage struct {
    BaseMessage
    Sender      string      `json:"sender"`
    Receiver    string      `json:"receiver"`
}
