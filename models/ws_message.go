package models

type WsMessage struct {
    BaseMessage
    Target          string                 `json:"target"`
    Code            string                 `json:"code"`
    Token           string                 `json:"token"`
}
