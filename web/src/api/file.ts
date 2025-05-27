import request from '@/config/request'
import type { AxiosProgressEvent } from 'axios'
import type { FileInfo, FileStats, Favorite, UploadConfig } from '@/types/file'

/**
 * 获取文件列表
 * @param path 目录路径
 * @returns 文件列表
 */
export function getFileList(path = '/') {
  return request.get<{ data: FileInfo[] }>('/file/list', {
    params: { path }
  })
}

/**
 * 搜索文件
 * @param keyword 搜索关键词
 * @returns 搜索结果
 */
export function searchFiles(keyword: string) {
  return request.get<{ data: FileInfo[] }>('/file/search', {
    params: { keyword }
  })
}

/**
 * 创建文件夹
 * @param path 文件夹路径
 * @returns 创建结果
 */
export function createFolder(path: string) {
  return request.post('/file/create-folder', null, {
    params: { path }
  })
}

/**
 * 删除文件或文件夹
 * @param path 文件路径
 * @param force 是否强制删除
 * @returns 删除结果
 */
export function deleteFile(path: string, force = false) {
  return request.delete('/file/delete', {
    params: { path, force }
  })
}

/**
 * 重命名文件或文件夹
 * @param path 原文件路径
 * @param newName 新文件名
 * @returns 重命名结果
 */
export function renameFile(path: string, newName: string) {
  return request.post('/file/rename', null, {
    params: { path, newName }
  })
}

/**
 * 复制文件或文件夹
 * @param srcPath 源文件路径
 * @param destPath 目标文件路径
 * @returns 复制结果
 */
export function copyFile(srcPath: string, destPath: string) {
  return request.post('/file/copy', null, {
    params: { srcPath, destPath }
  })
}

/**
 * 移动文件或文件夹
 * @param srcPath 源文件路径
 * @param destPath 目标文件路径
 * @returns 移动结果
 */
export function moveFile(srcPath: string, destPath: string) {
  return request.post('/file/move', null, {
    params: { srcPath, destPath }
  })
}

/**
 * 上传文件
 * @param file 文件对象
 * @param config 上传配置
 * @param onProgress 进度回调
 * @returns 上传结果
 */
export function uploadFile(
  file: File,
  config: UploadConfig,
  onProgress?: (progressEvent: AxiosProgressEvent) => void
) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('fileName', config.fileName)
  formData.append('path', config.filePath)
  formData.append('override', String(config.override))

  if (config.relativePath) {
    formData.append('relativePath', config.relativePath)
  }

  return request.post('/file/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: onProgress
  })
}

/**
 * 分块上传文件
 * @param chunk 文件分块
 * @param config 上传配置
 * @param onProgress 进度回调
 * @returns 上传结果
 */
export function uploadChunk(
  chunk: Blob,
  config: UploadConfig,
  onProgress?: (progressEvent: AxiosProgressEvent) => void
) {
  const formData = new FormData()
  formData.append('file', chunk)
  formData.append('fileName', config.fileName)
  formData.append('path', config.filePath)
  formData.append('override', String(config.override))

  if (config.relativePath) {
    formData.append('relativePath', config.relativePath)
  }

  if (config.chunks) {
    formData.append('chunkIndex', String(config.chunks.chunkIndex))
    formData.append('totalChunks', String(config.chunks.totalChunks))
  }

  return request.post('/file/upload-chunk', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: onProgress
  })
}

/**
 * 合并文件分块
 * @param fileName 文件名
 * @param path 文件路径
 * @param totalChunks 总分块数
 * @param override 是否覆盖
 * @param relativePath 相对路径
 * @returns 合并结果
 */
export function mergeChunks(
  fileName: string,
  path: string,
  totalChunks: number,
  override = false,
  relativePath?: string
) {
  return request.post('/file/merge-chunks', {
    fileName,
    path,
    totalChunks,
    override,
    relativePath
  })
}

/**
 * 获取文件统计信息
 * @param path 目录路径
 * @returns 统计信息
 */
export function getFileStats(path = '/') {
  return request.get<{ data: FileStats }>('/file/stats', {
    params: { path }
  })
}

/**
 * 添加收藏
 * @param path 文件路径
 * @returns 添加结果
 */
export function addFavorite(path: string) {
  return request.post('/file/favorite', null, {
    params: { path }
  })
}

/**
 * 移除收藏
 * @param path 文件路径
 * @returns 移除结果
 */
export function removeFavorite(path: string) {
  return request.delete('/file/favorite', {
    params: { path }
  })
}

/**
 * 获取收藏列表
 * @returns 收藏列表
 */
export function getFavorites() {
  return request.get<{ data: Favorite[] }>('/file/favorites')
}

/**
 * 获取文件下载URL
 * @param filePath 文件路径
 * @returns 下载URL
 */
export function getDownloadUrl(filePath: string): string {
  return `/api/file/download?path=${encodeURIComponent(filePath)}`
}

/**
 * 获取文件预览URL
 * @param filePath 文件路径
 * @returns 预览URL
 */
export function getPreviewUrl(filePath: string): string {
  return `/api/file/preview?path=${encodeURIComponent(filePath)}`
}
