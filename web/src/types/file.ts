/**
 * 文件信息接口
 */
export interface FileInfo {
  fileName: string // 文件名
  filePath: string // 文件路径
  fileSize: number // 文件大小（字节）
  fileType: string // 文件MIME类型
  isDir: boolean // 是否为目录
  modTime: string // 修改时间
  extension?: string // 文件扩展名
  icon?: string // 文件图标
}

/**
 * 文件统计信息接口
 */
export interface FileStats {
  totalFiles: number // 文件总数
  totalFolders: number // 文件夹总数
  totalSize: number // 总大小（字节）
}

/**
 * 收藏信息接口
 */
export interface Favorite {
  id: number // 收藏ID
  name: string // 文件名
  path: string // 文件路径
  isDir: boolean // 是否是目录
  createTime: string // 创建时间
}

/**
 * 文件操作类型
 */
export type FileOperation = 'copy' | 'move' | 'delete' | 'rename' | 'download' | 'favorite'

/**
 * 视图模式类型
 */
export type ViewMode = 'grid' | 'list'

/**
 * 排序字段类型
 */
export type SortField = 'name' | 'size' | 'type' | 'modTime'

/**
 * 排序方向类型
 */
export type SortOrder = 'asc' | 'desc'

/**
 * 文件排序配置
 */
export interface SortConfig {
  field: SortField
  order: SortOrder
}

/**
 * 文件上传配置
 */
export interface UploadConfig {
  fileName: string
  filePath: string
  fileSize: number
  override: boolean
  relativePath?: string // 文件相对路径，用于文件夹上传
  chunks?: {
    chunkIndex: number
    totalChunks: number
  }
}

/**
 * 上传进度信息
 */
export interface UploadProgress {
  fileName: string
  loaded: number
  total: number
  percentage: number
  status: 'pending' | 'uploading' | 'success' | 'error'
  error?: string
  relativePath?: string // 文件相对路径，用于文件夹上传
}

/**
 * API响应接口
 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

/**
 * 面包屑导航项
 */
export interface BreadcrumbItem {
  label: string
  path: string
}

/**
 * 文件预览配置
 */
export interface PreviewConfig {
  isPreviewable: boolean
  previewType: 'image' | 'video' | 'audio' | 'text' | 'pdf' | 'none'
  previewUrl?: string
}
