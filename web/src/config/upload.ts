import { storage } from '@/utils/storage'

// 默认配置
const defaultConfig = {
  // 分块大小（单位：字节）
  chunkSize: 2 * 1024 * 1024, // 默认2MB
  
  // 最大并发上传数
  maxConcurrent: 3,
  
  // 是否启用分块上传
  enableChunked: true,
  
  // 触发分块上传的文件大小阈值（单位：字节）
  chunkThreshold: 2 * 1024 * 1024 // 默认2MB以上使用分块上传
}

// 从本地存储获取配置，如果没有则使用默认配置
const storedConfig = storage.get('upload_config')
export const uploadConfig = storedConfig || defaultConfig

// 更新配置
export const updateUploadConfig = (newConfig: Partial<typeof defaultConfig>) => {
  Object.assign(uploadConfig, newConfig)
  storage.set('upload_config', uploadConfig)
}
