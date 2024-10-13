package model

/**
  @author: XingGao
  @date: 2024/10/11
**/

type FileInfo struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
	IsDir    bool   `json:"isDir"`
	ModTime  string `json:"modTime"`
}
