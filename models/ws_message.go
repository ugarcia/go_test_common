package models

type WsMessage struct {
    BaseMessage
    Token       string      `json:"token"`
}
