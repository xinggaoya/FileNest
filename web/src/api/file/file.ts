import { get, post, del } from '@/utils/request'
import type { AxiosProgressEvent } from 'axios'
import { isAxiosError } from 'axios'

export interface FileInfo {
  fileName: string
  filePath: string
  fileSize: number
  fileType?: string
  isDir: boolean
  modTime: string
}

export interface GetFileListParams {
  path: string
}

export interface CreateFolderParams {
  path: string
}

export interface DeleteFileParams {
  path: string
}

export interface UploadParams {
  indexChunk: number
  totalChunks: number
  fileName: string
  path: string
  override: boolean
}

export interface FileStats {
  totalFiles: number
  totalFolders: number
  totalSize: number
}

export interface GetFileStatsParams {
  path: string
}

export interface SearchFilesParams {
  keyword: string
}

export interface Favorite {
  id: number
  filePath: string
  fileName: string
  isDir: boolean
  createdAt: string
}

export interface UploadFileParams {
  file: File
  path: string
  onProgress?: (progress: number) => void
  onSuccess?: () => void
  onError?: (error: string) => void
}

export const getFileList = (params: GetFileListParams) => {
  return get<FileInfo[]>('/file/list', params)
}

export const createFolder = (path: string) => {
  return post('/file/create-folder', null, { params: { path } })
}

export const deleteFile = (path: string) => {
  return del('/file/delete', { path })
}

export const downloadFile = (path: string) => {
  return get('/file/download', { path }, { responseType: 'blob' })
}

/**
 * 上传文件
 */
export const uploadFile = async ({
  file,
  path,
  onProgress,
  onSuccess,
  onError
}: UploadFileParams) => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('fileName', file.name)
  formData.append('path', path)

  try {
    await post('/file/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress?.(progress)
        }
      }
    })
    onSuccess?.()
  } catch (error) {
    if (isAxiosError(error)) {
      onError?.(error.response?.data?.message || '上传失败')
    } else {
      onError?.('上传失败')
    }
  }
}

/**
 * 获取文件下载链接
 * @param path 文件路径
 */
export const getDownloadUrl = (path: string) => {
  return `/api/file/download?path=${encodeURIComponent(path)}`
}

// 获取文件统计信息
export const getFileStats = (params: GetFileStatsParams) => {
  return get<FileStats>('/file/stats', params)
}

export const searchFiles = (params: SearchFilesParams) => {
  return get<FileInfo[]>('/file/search', params)
}

export const addFavorite = (path: string) => {
  return post('/file/favorite', null, { params: { path } })
}

export const removeFavorite = (path: string) => {
  return del('/file/favorite', { path })
}

export const getFavorites = () => {
  return get<Favorite[]>('/file/favorites')
}
