package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"
	"tool/global/variable"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ResponseWriter 包装 gin.ResponseWriter 以捕获响应数据
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 使用 zap 和 lumberjack 记录 HTTP 请求日志的中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 捕获响应数据
		w := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// 读取请求数据
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// 请求方法
		reqMethod := c.Request.Method
		// 请求路由
		reqURI := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		// 设置日志中记录字符串的最大长度
		maxLogLength := variable.ConfigYml.GetInt("Logs.ResponseLengthMax") // 最大日志长度

		// 请求体日志处理
		var logReqBody string
		if len(reqBody) > maxLogLength {
			logReqBody = fmt.Sprintf("Request body too large to log, length: %d", len(reqBody))
		} else {
			logReqBody = string(reqBody)
		}

		// 响应体日志处理
		var logRespBody string
		if w.body.Len() > maxLogLength {
			logRespBody = fmt.Sprintf("Response body too large to log, length: %d", w.body.Len())
		} else {
			logRespBody = w.body.String()
		}

		// 记录日志
		logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.Duration("latency", latencyTime),
			zap.String("clientIP", clientIP),
			zap.String("method", reqMethod),
			zap.String("path", reqURI),
			zap.String("reqBody", logReqBody),
			zap.String("respBody", logRespBody),
		)
	}
}

// InitLogger 初始化 zap.Logger 和 lumberjack
func InitLogger() (*zap.Logger, error) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"), // 日志文件路径
		MaxSize:    variable.ConfigYml.GetInt("Logs.MaxSize"),                           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: variable.ConfigYml.GetInt("Logs.MaxBackups"),                        // 日志文件最多保存多少个备份
		MaxAge:     variable.ConfigYml.GetInt("Logs.MaxAge"),                            // 文件最多保存多少天
		Compress:   variable.ConfigYml.GetBool("Logs.Compress"),                         // 是否压缩
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "created_at"

	timePrecision := variable.ConfigYml.GetString("Logs.TimePrecision")
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"
	}

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}
