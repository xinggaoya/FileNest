import { get, post, del, download } from '@/utils/request'
import type { AxiosProgressEvent } from 'axios'
import type { FileInfo, Favorite, FileStats } from '@/types/file'
import { uploadConfig } from '@/config/upload'

export interface UploadParams {
  indexChunk: number
  totalChunks: number
  fileName: string
  path: string
  override: boolean
}

export interface GetFileStatsParams {
  path: string
}

export interface SearchFilesParams {
  keyword: string
}

export interface UploadFileParams {
  file: File
  path: string
  override?: boolean
  onProgress?: (progress: number) => void
  onSuccess?: () => void
  onError?: (error: string) => void
}

export const getFileList = (path: string) => {
  return get<FileInfo[]>('/file/list', { path })
}

export const createFolder = (path: string) => {
  return post('/file/create-folder', null, { params: { path } })
}

export const deleteFile = (path: string, force: boolean = false) => {
  return del('/file/delete', { path, force })
}

export const downloadFile = (path: string) => {
  return download('/file/download', { path })
}

/**
 * 上传文件分块
 */
export const uploadFileChunk = async ({
  file,
  chunk,
  path,
  fileName,
  chunkIndex,
  totalChunks,
  override = false,
  onProgress,
  onSuccess,
  onError
}: {
  file: File
  chunk: Blob
  path: string
  fileName: string
  chunkIndex: number
  totalChunks: number
  override?: boolean
  onProgress?: (progress: number) => void
  onSuccess?: () => void
  onError?: (error: string) => void
}) => {
  const formData = new FormData()
  formData.append('file', chunk)
  formData.append('fileName', fileName)
  formData.append('path', path)
  formData.append('override', override.toString())
  formData.append('chunkIndex', chunkIndex.toString())
  formData.append('totalChunks', totalChunks.toString())

  try {
    await post('/file/upload-chunk', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent: AxiosProgressEvent) => {
        if (progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress?.(progress)
        }
      }
    })
    onSuccess?.()
  } catch (error: any) {
    onError?.(error.response?.data?.message || '上传失败')
  }
}

/**
 * 合并文件分块
 */
export const mergeFileChunks = async ({
  path,
  fileName,
  totalChunks,
  override = false,
  onSuccess,
  onError
}: {
  path: string
  fileName: string
  totalChunks: number
  override?: boolean
  onSuccess?: () => void
  onError?: (error: string) => void
}) => {
  try {
    await post('/file/merge-chunks', {
      fileName,
      path,
      totalChunks,
      override
    })
    onSuccess?.()
  } catch (error: any) {
    onError?.(error.response?.data?.message || '合并文件失败')
  }
}

/**
 * 上传文件
 */
export const uploadFile = async ({
  file,
  path,
  override = false,
  onProgress,
  onSuccess,
  onError
}: UploadFileParams) => {
  // 如果文件小于阈值或未启用分块上传，直接上传
  if (!uploadConfig.enableChunked || file.size <= uploadConfig.chunkThreshold) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('fileName', file.name)
    formData.append('path', path)
    formData.append('override', override.toString())

    try {
      await post('/file/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        },
        onUploadProgress: (progressEvent: AxiosProgressEvent) => {
          if (progressEvent.total) {
            const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
            onProgress?.(progress)
          }
        }
      })
      onSuccess?.()
    } catch (error: any) {
      onError?.(error.response?.data?.message || '上传失败')
    }
    return
  }

  // 文件分块
  const totalChunks = Math.ceil(file.size / uploadConfig.chunkSize)
  
  try {
    // 创建进度追踪器
    const chunkProgress = new Array(totalChunks).fill(0)
    const updateTotalProgress = () => {
      const totalProgress = Math.round(
        chunkProgress.reduce((acc, curr) => acc + curr, 0) / totalChunks
      )
      onProgress?.(totalProgress)
    }

    // 并发控制
    let completedChunks = 0
    
    // 创建所有分块的上传任务
    const uploadTasks = Array.from({ length: totalChunks }, (_, index) => {
      const start = index * uploadConfig.chunkSize
      const end = Math.min(start + uploadConfig.chunkSize, file.size)
      const chunk = file.slice(start, end)
      
      return async () => {
        await uploadFileChunk({
          file,
          chunk,
          path,
          fileName: file.name,
          chunkIndex: index,
          totalChunks,
          override,
          onProgress: (progress) => {
            chunkProgress[index] = progress
            updateTotalProgress()
          },
          onSuccess: () => {
            completedChunks++
            chunkProgress[index] = 100
            updateTotalProgress()
          },
          onError
        })
      }
    })

    // 执行分块上传
    const executeUploadTasks = async () => {
      const pendingTasks = [...uploadTasks]
      const runningTasks = new Set<Promise<void>>()

      while (pendingTasks.length > 0 || runningTasks.size > 0) {
        // 填充运行中的任务，直到达到最大并发数
        while (runningTasks.size < uploadConfig.maxConcurrent && pendingTasks.length > 0) {
          const task = pendingTasks.shift()!
          const promise = task().then(() => {
            runningTasks.delete(promise)
          })
          runningTasks.add(promise)
        }

        if (runningTasks.size > 0) {
          // 等待任意一个任务完成
          await Promise.race(Array.from(runningTasks))
        }
      }
    }

    await executeUploadTasks()

    // 所有分块上传完成后，请求合并文件
    if (completedChunks === totalChunks) {
      await mergeFileChunks({
        path,
        fileName: file.name,
        totalChunks,
        override,
        onSuccess,
        onError
      })
    }
  } catch (error: any) {
    onError?.(error.response?.data?.message || '上传失败')
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

export const searchFiles = (keyword: string) => {
  return get<FileInfo[]>('/file/search', { keyword })
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

export const renameFile = (path: string, newName: string) => {
  return post('/file/rename', null, { params: { path, newName } })
}

export const copyFile = (srcPath: string, destPath: string) => {
  return post('/file/copy', null, { params: { srcPath, destPath } })
}

export const moveFile = (srcPath: string, destPath: string) => {
  return post('/file/move', null, { params: { srcPath, destPath } })
}
