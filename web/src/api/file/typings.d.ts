declare namespace API {
  interface FileInfo {
    fileName: string
    filePath: string
    isDir: boolean
    size: number
    modTime: string
  }

  interface FileStats {
    totalFiles: number
    totalFolders: number
    totalSize: number
  }
} 