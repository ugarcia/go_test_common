package models

type BaseMessage struct {
    Sender          string                 `json:"sender"`
    ConnectionId    uint                   `json:"connection_id"`
    Broadcast       bool                   `json:"broadcast"`
    Action          string                 `json:"action"`
    Data            map[string]interface{} `json:"data"`
}
