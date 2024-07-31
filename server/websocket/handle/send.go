package handle

import "github.com/gorilla/websocket"

// type Person struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// 发送信息
type Msg struct {
	Content  []byte `json:"content"`   // 消息内容
	Type     string `json:"type"`      // 消息类型 1:文本 2:图片 3:视频
	UserName string `json:"user_name"` // 用户名
}

//发送信息
func sendMsg(cli *websocket.Conn, msg []byte) {
	// person := Msg{Name: "Alice", Age: 30}

	// // 编码为 JSON
	// jsonBytes, err := json.Marshal(person)
	// if err != nil {
	// 	// 处理错误
	// }

	// fmt.Println(string(jsonBytes))
}
