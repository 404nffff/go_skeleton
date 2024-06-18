package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"tool/app/global/variable"
)

// 嵌入 templates 目录下的所有文件
//
//go:embed admin/layouts/*.html
//go:embed admin/*.html
//go:embed admin/user/*.html
var Content embed.FS

// 嵌入 components 目录下的所有文件
//
//go:embed components/*
var Components embed.FS

// Load 封装模板加载逻辑
func Load() (*template.Template, error) {
	return template.New("").ParseFS(Content, "admin/*.html", "admin/user/*.html", "admin/layouts/*.html")
}

// 列出嵌入文件系统中的所有文件
func List() {

	//初始化接受的数组
	var pathArray []string

	// 列出嵌入文件系统中的所有文件
	fs.WalkDir(Content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		pathArray = append(pathArray, path)
		return nil
	})

	fs.WalkDir(Components, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		pathArray = append(pathArray, path)
		return nil
	})

	// 打印所有嵌入文件
	variable.Logs.Sugar().Info(fmt.Sprintf("嵌入文件系统中的所有文件: %v", pathArray))
}
