package models

import (
	"General_Framework_Gin/config"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// UpdateInfo 更新信息结构体
type UpdateInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// DownloadFile 下载文件到指定路径
func DownloadFile(url, dest string) error {
	log.Printf("开始下载文件: %s", url)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	log.Printf("文件已下载到: %s", dest)
	return nil
}

// ReplaceAndRestart 替换当前程序并重启
func ReplaceAndRestart(newAppPath string) error {
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前程序路径失败: %v", err)
	}
	currentPath = filepath.Clean(currentPath)

	scriptPath := "./update_and_restart.sh"
	if err := CreateUpdateScript(scriptPath, currentPath, newAppPath); err != nil {
		return fmt.Errorf("生成脚本失败: %v", err)
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("脚本文件未成功创建")
	}

	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("无法设置脚本的执行权限: %v", err)
	}

	cmd := exec.Command("sh", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动脚本失败: %v", err)
	}
	log.Println("更新脚本已启动，退出当前程序...")
	os.Exit(0)
	return nil
}

// CreateUpdateScript 创建用于替换和重启的脚本
func CreateUpdateScript(scriptPath, currentPath, newAppPath string) error {
	scriptContent := fmt.Sprintf(`#!/bin/bash
killall "$(basename %s)"
if [ -f "%s" ]; then
	rm -f "%s"
fi
if [ -f "%s" ]; then
	mv "%s" "%s"
else
	echo "新程序文件不存在，更新失败！" >&2
	exit 1
fi
echo "请手动启动新程序： ./$(basename %s)"
rm -- "$0"`, currentPath, currentPath, currentPath, newAppPath, newAppPath, currentPath, currentPath)
	return os.WriteFile(scriptPath, []byte(scriptContent), 0755)
}

// InitUpdate 检查更新并执行更新
func InitUpdate(NoteSegment []byte) {
	serverURL := config.AppConfig.Update.ServerURL
	platform := config.AppConfig.Update.Platform
	appName := config.AppConfig.Update.AppName
	log.Println("检查更新中...")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	checkURL := fmt.Sprintf("%s/check-update?platform=%s&app=%s", serverURL, platform, appName)
	resp, err := client.Get(checkURL)
	if err != nil {
		log.Printf("检查更新失败：%v", err)
		return
	}
	defer resp.Body.Close()

	var updateInfo UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&updateInfo); err != nil {
		log.Println("解析更新信息失败")
		return
	}

	localVersion, err := GetProductVersion(NoteSegment)
	if err != nil {
		log.Printf("获取本地文件版本失败：%v", err)
		return
	}

	if localVersion < updateInfo.Version {
		log.Println("发现新版本，正在下载更新...")
		tempFile := "./temp_app"
		if err := DownloadFile(updateInfo.URL, tempFile); err != nil {
			log.Printf("下载失败：%v", err)
			return
		}
		log.Println("下载完成，正在替换程序...")
		if err := ReplaceAndRestart(tempFile); err != nil {
			log.Printf("更新失败：%v", err)
		} else {
			log.Println("更新完成，程序已重启！")
		}
	} else {
		log.Println("当前是最新版本，无需更新。")
	}
}

// GetProductVersion 获取 ELF 文件中的版本信息
func GetProductVersion(NoteSegment []byte) (string, error) {
	if len(NoteSegment) < 12 {
		return "", fmt.Errorf("NoteSegment 数据过短，无法解析")
	}

	namesz := binary.LittleEndian.Uint32(NoteSegment[0:4])
	descsz := binary.LittleEndian.Uint32(NoteSegment[4:8])
	noteType := binary.LittleEndian.Uint32(NoteSegment[8:12])

	nameStart := 12
	nameEnd := nameStart + int(namesz)
	if nameEnd > len(NoteSegment) {
		return "", fmt.Errorf("NoteSegment 名称超出范围")
	}

	name := bytes.Trim(NoteSegment[nameStart:nameEnd], "\x00")
	if string(name) != "Version" {
		return "", fmt.Errorf("NoteSegment 名称不是 'Version'")
	}

	descStart := nameEnd
	descEnd := descStart + int(descsz)
	if descEnd > len(NoteSegment) {
		return "", fmt.Errorf("NoteSegment 描述符超出范围")
	}
	desc := bytes.Trim(NoteSegment[descStart:descEnd], "\x00")
	if noteType != 1 {
		return "", fmt.Errorf("NoteSegment 类型无效，期望值为 1，实际值为 %d", noteType)
	}
	return string(desc), nil
}
