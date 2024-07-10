package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"tool/app/global/variable"
	"tool/pkg/oss"
	"tool/pkg/session"
	"tool/pkg/udp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// User 定义用户模型
type User struct {
	ID            uint   `gorm:"primaryKey"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Avatar        string `json:"avatar"`
	Type          int    `json:"type"`
	LastLoginTime int    `json:"last_login_time"`
	LoginStatus   int    `json:"login_status"`
	CreateTime    int    `json:"create_time"`                           // 自定义 create_time 字段
	UpdateTime    int    `json:"update_time"`                           // 自定义 update_time 字段
	DeleteTime    int    `gorm:"column:delete_time" json:"delete_time"` // 自定义 delete_time 字段
}

// TableName 设置表名前缀
func (User) TableName() string {
	return "h_user"
}

// @Summary 测试接口
// @Description 测试接口
// @Tags 测试接口
// @Accept json
// @Produce json
// @Success 200 {string} string "{"message": "Hello, World!"}"
// @Router /api/v1/test [get]
func Test(c *gin.Context) {
	session.Set(c, "test", "test222")

	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

func Test2(c *gin.Context) {
	sessionTest := session.Get(c, "test")

	fmt.Println(sessionTest)
	c.JSON(200, gin.H{
		"message": sessionTest,
	})
}

// GetUsersHandler 获取用户列表
func GetUsersHandler(c *gin.Context) {
	var users []User
	if err := variable.Mysql.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"users": users})
}

func GetUserMongo(c *gin.Context) {

	params, _ := c.GetQuery("name")

	//Example query usage
	//filter := bson.D{{"updated_at", "2023-03-21 16:54:53"}}
	result := variable.MongoDB.Collection("t_chatgpt_log").FindOne(c, bson.D{})

	// Decode the result
	var doc map[string]interface{}
	err := result.Decode(&doc)
	if err != nil {
		variable.Logs.Error("Failed to decode document", zap.Error(err))
	} else {
		fmt.Printf("Found a document: %v\n", doc)
		variable.Logs.Info("Found a document", zap.Any("document", doc))

		c.JSON(200, gin.H{"users": doc, "params": params})
	}

}

func TestUdp(c *gin.Context) {
	udp.Send("test")
}

func TestAnt(c *gin.Context) {

	// 获取请求参数
	//param := c.Query("param")

	// 示例调用
	// params := map[string]any{
	// 	"age":  30,
	// 	"name": "123123123",
	// }
	// result, err := variable.Pool.SubmitTask(myTask, params)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("Result:", result)
	// }

	//c.JSON(http.StatusOK, gin.H{"result": result})

	task := func(params map[string]any) map[string]any {
		// 模拟一个任务
		time.Sleep(5 * time.Second)
		variable.Logs.Info("Task completed")
		return nil
	}
	variable.Pool.SubmitTask(task, map[string]any{})

	c.JSON(http.StatusOK, gin.H{"message": "task submitted"})
}

func StatusAnt(c *gin.Context) {
	status, s := variable.Pool.GetStatus()

	c.JSON(http.StatusOK, gin.H{"status": status, "s": s})
}

// 示例任务函数
func myTask(params map[string]any) string {

	// Your task implementation
	name, ok := params["name"]
	if !ok {
		return "Invalid parameter"
	}
	age, ok := params["age"]
	if !ok {
		return "Invalid parameter"
	}
	return fmt.Sprintf("Processed: name=%s, age=%d", name, age)
}

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// 保存文件到临时目录
	tempFilePath := filepath.Join(os.TempDir(), header.Filename)
	out, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}
	defer out.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	_, err = out.Write(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	ossConfig := oss.OssConfig{
		Endpoint:        "oss-cn-beijing.aliyuncs.com",
		AccessKeyID:     "DdZh6CfOjy0mKFe4",
		AccessKeySecret: "uQNOq2AU3uTm0DdMnmL6NXU8H7PUoK",
		BucketName:      "wechat-app-yuyue-assets",
	}

	ossClient, _ := oss.NewOSSClient(ossConfig)

	// 上传文件到 OSS
	fileName, err := ossClient.UploadFileFromPath(header.Filename, tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除临时文件
	os.Remove(tempFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_name": fileName})
}
