package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// HandleFileUpload 处理大文件、多文件和单文件上传，支持断点续传
func HandleFileUpload(ctx *gin.Context, uploadDir string) error {
	files := ctx.Request.MultipartForm.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "没有上传的文件"})
		return nil
	}

	var uploadedFiles []map[string]string
	for _, file := range files {
		filePath := filepath.Join(uploadDir, file.Filename)

		out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("文件打开失败: %v", err)
		}
		defer out.Close()

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("文件读取失败: %v", err)
		}
		defer src.Close()

		buffer := make([]byte, 10*1024*1024) // 10MB 缓冲区
		for {
			n, readErr := src.Read(buffer)
			if n > 0 {
				if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
					return fmt.Errorf("文件写入失败: %v", writeErr)
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				return fmt.Errorf("文件读取错误: %v", readErr)
			}
		}

		uploadedFiles = append(uploadedFiles, map[string]string{
			"filename": file.Filename,
			"filepath": filePath,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"files":   uploadedFiles,
	})
	return nil
}

// HandleFileDownload 处理大文件、多文件和单文件下载，支持断点续传
func HandleFileDownload(ctx *gin.Context, downloadDir string) {
	filename := ctx.Query("filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	filePath := filepath.Join(downloadDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("文件未找到: %s", filename)})
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader == "" {
		ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		ctx.File(filePath)
		return
	}

	var start int64
	n, _ := fmt.Sscanf(rangeHeader, "bytes=%d-", &start)
	if n != 1 || start >= fileSize {
		ctx.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": "无效的范围请求"})
		return
	}

	ctx.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, fileSize-1, fileSize))
	ctx.Writer.Header().Set("Accept-Ranges", "bytes")
	ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileSize-start, 10))
	ctx.Writer.WriteHeader(http.StatusPartialContent)

	file.Seek(start, io.SeekStart)
	buffer := make([]byte, 10*1024*1024) // 10MB 缓冲区
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			ctx.Writer.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
			return
		}
	}
}
