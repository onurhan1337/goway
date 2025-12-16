package ws

type Message struct {
    Action  string `json:"action"`
    Room    string `json:"room"`
    Content string `json:"content"`
}
