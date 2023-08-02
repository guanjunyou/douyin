package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type ResponseBody struct {
	Url      string `json:"url"`
	CoverUrl string `json:"cover_url"`
}

func UploadToServer(data *multipart.FileHeader) (string, error) {
	fileName := data.Filename
	fileType := getFileCategory(getFileType(fileName))
	// 构建HTTP请求的Body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件数据
	file, _ := data.Open()
	defer file.Close()

	part, err := writer.CreateFormFile("file", data.Filename)
	if err != nil {
		return "", fmt.Errorf("error writing file to request: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file to request: %w", err)
	}

	// 添加fileType参数
	err = writer.WriteField("fileType", fmt.Sprintf("%d", fileType))
	if err != nil {
		return "", fmt.Errorf("error writing filetype to request: %w", err)
	}

	// 结束写入
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest(
		config.Config.VideoServer.Api.Upload.Method,
		"http://"+config.Config.VideoServer.Addr+config.Config.VideoServer.Api.Upload.Path,
		body,
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", respBody)
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("resource upload failed, HTTP status code: %d", resp.StatusCode)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(respBody, &responseBody)
	fmt.Println("资源保存成功！")
	return responseBody.CoverUrl, nil
}

func getFileType(filename string) string {
	// 使用path/filepath包的Ext函数获取文件扩展名
	extension := filepath.Ext(filename)

	// 去除扩展名中的点号，并转换为小写字母
	fileType := strings.ToLower(strings.TrimPrefix(extension, "."))

	return fileType
}

func getFileCategory(fileType string) int {
	// 视频类型判断
	videoTypes := []string{"mp4", "avi", "mov"}
	for _, t := range videoTypes {
		if fileType == t {
			return 1 // 视频类返回1
		}
	}

	// 图片类型判断
	imageTypes := []string{"jpg", "jpeg", "png", "gif"}
	for _, t := range imageTypes {
		if fileType == t {
			return 2 // 图片类返回2
		}
	}

	// 其他类型
	return 0
}
