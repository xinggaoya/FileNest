export interface FileInfo {
  fileName: string
  filePath: string
  fileSize: number
  modTime: string
  isDir: boolean
}

export interface FileStats {
  totalFiles: number
  totalFolders: number
  totalSize: number
}

export interface Favorite {
  id: number
  name: string
  path: string
  isDir: boolean
  createTime: string
} 