package models

type BaseMessage struct {
    Source          string                 `json:"source"`
    Target          string                 `json:"target"`
    Code            string                 `json:"code"`
    Action          string                 `json:"action"`
    ConnectionType  string                 `json:"connection_type"`
    ConnectionId    uint                   `json:"connection_id"`
    Broadcast       bool                   `json:"broadcast"`
    Data            map[string]interface{} `json:"data"`
}
