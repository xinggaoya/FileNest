import type { FileInfo, PreviewConfig } from '@/types/file'

/**
 * 格式化文件大小
 * @param bytes 字节数
 * @returns 格式化后的文件大小字符串
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 获取文件扩展名
 * @param fileName 文件名
 * @returns 文件扩展名（小写）
 */
export function getFileExtension(fileName: string): string {
  const lastDotIndex = fileName.lastIndexOf('.')
  return lastDotIndex === -1 ? '' : fileName.slice(lastDotIndex + 1).toLowerCase()
}

/**
 * 获取文件类型图标名称
 * @param file 文件信息
 * @returns 图标名称
 */
export function getFileIcon(file: FileInfo): string {
  if (file.isDir) {
    return 'folder'
  }

  const extension = getFileExtension(file.fileName)

  // 图片文件
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp', 'ico'].includes(extension)) {
    return 'image'
  }

  // 视频文件
  if (['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm', 'm4v'].includes(extension)) {
    return 'video'
  }

  // 音频文件
  if (['mp3', 'wav', 'flac', 'aac', 'ogg', 'wma', 'm4a'].includes(extension)) {
    return 'audio'
  }

  // 文档文件
  if (['pdf'].includes(extension)) {
    return 'pdf'
  }

  if (['doc', 'docx'].includes(extension)) {
    return 'word'
  }

  if (['xls', 'xlsx'].includes(extension)) {
    return 'excel'
  }

  if (['ppt', 'pptx'].includes(extension)) {
    return 'powerpoint'
  }

  // 压缩文件
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(extension)) {
    return 'archive'
  }

  // 代码文件
  if (
    [
      'js',
      'ts',
      'vue',
      'jsx',
      'tsx',
      'html',
      'css',
      'scss',
      'sass',
      'less',
      'json',
      'xml',
      'yaml',
      'yml'
    ].includes(extension)
  ) {
    return 'code'
  }

  if (['txt', 'md', 'log'].includes(extension)) {
    return 'text'
  }

  return 'file'
}

/**
 * 检查文件是否可预览
 * @param file 文件信息
 * @returns 预览配置
 */
export function getPreviewConfig(file: FileInfo): PreviewConfig {
  if (file.isDir) {
    return { isPreviewable: false, previewType: 'none' }
  }

  const extension = getFileExtension(file.fileName)

  // 图片文件
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'].includes(extension)) {
    return {
      isPreviewable: true,
      previewType: 'image',
      previewUrl: `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    }
  }

  // 视频文件
  if (['mp4', 'webm', 'ogg'].includes(extension)) {
    return {
      isPreviewable: true,
      previewType: 'video',
      previewUrl: `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    }
  }

  // 音频文件
  if (['mp3', 'wav', 'ogg'].includes(extension)) {
    return {
      isPreviewable: true,
      previewType: 'audio',
      previewUrl: `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    }
  }

  // 文本文件
  if (['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'vue', 'log'].includes(extension)) {
    return {
      isPreviewable: true,
      previewType: 'text',
      previewUrl: `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    }
  }

  // PDF文件
  if (['pdf'].includes(extension)) {
    return {
      isPreviewable: true,
      previewType: 'pdf',
      previewUrl: `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    }
  }

  return { isPreviewable: false, previewType: 'none' }
}

/**
 * 生成文件下载URL
 * @param filePath 文件路径
 * @returns 下载URL
 */
export function getDownloadUrl(filePath: string): string {
  return `/api/file/download?path=${encodeURIComponent(filePath)}`
}

/**
 * 检查文件名是否有效
 * @param fileName 文件名
 * @returns 是否有效
 */
export function isValidFileName(fileName: string): boolean {
  if (!fileName || fileName.trim() === '') {
    return false
  }

  // 检查非法字符
  const invalidChars = /[<>:"/\\|?*]/
  if (invalidChars.test(fileName)) {
    return false
  }

  // 检查保留名称（Windows）
  const reservedNames = [
    'CON',
    'PRN',
    'AUX',
    'NUL',
    'COM1',
    'COM2',
    'COM3',
    'COM4',
    'COM5',
    'COM6',
    'COM7',
    'COM8',
    'COM9',
    'LPT1',
    'LPT2',
    'LPT3',
    'LPT4',
    'LPT5',
    'LPT6',
    'LPT7',
    'LPT8',
    'LPT9'
  ]
  const nameWithoutExt = fileName.split('.')[0].toUpperCase()
  if (reservedNames.includes(nameWithoutExt)) {
    return false
  }

  return true
}

/**
 * 解析路径为面包屑导航
 * @param path 文件路径
 * @returns 面包屑导航数组
 */
export function parseBreadcrumbs(path: string): Array<{ label: string; path: string }> {
  if (!path || path === '/') {
    return [{ label: '根目录', path: '/' }]
  }

  const parts = path.split('/').filter((part) => part !== '')
  const breadcrumbs = [{ label: '根目录', path: '/' }]

  let currentPath = ''
  for (const part of parts) {
    currentPath += '/' + part
    breadcrumbs.push({
      label: part,
      path: currentPath
    })
  }

  return breadcrumbs
}

/**
 * 获取父级目录路径
 * @param path 当前路径
 * @returns 父级目录路径
 */
export function getParentPath(path: string): string {
  if (!path || path === '/') {
    return '/'
  }

  const parts = path.split('/').filter((part) => part !== '')
  if (parts.length <= 1) {
    return '/'
  }

  parts.pop()
  return '/' + parts.join('/')
}

/**
 * 合并路径
 * @param basePath 基础路径
 * @param relativePath 相对路径
 * @returns 合并后的路径
 */
export function joinPath(basePath: string, relativePath: string): string {
  if (!basePath || basePath === '/') {
    return '/' + relativePath
  }

  return basePath.endsWith('/') ? basePath + relativePath : basePath + '/' + relativePath
}

/**
 * 排序文件列表
 * @param files 文件列表
 * @param field 排序字段
 * @param order 排序方向
 * @returns 排序后的文件列表
 */
export function sortFiles(files: FileInfo[], field: string, order: 'asc' | 'desc'): FileInfo[] {
  return [...files].sort((a, b) => {
    // 目录优先
    if (a.isDir && !b.isDir) return -1
    if (!a.isDir && b.isDir) return 1

    let comparison = 0

    switch (field) {
      case 'name':
        comparison = a.fileName.localeCompare(b.fileName)
        break
      case 'size':
        comparison = a.fileSize - b.fileSize
        break
      case 'type':
        comparison = getFileExtension(a.fileName).localeCompare(getFileExtension(b.fileName))
        break
      case 'modTime':
        comparison = new Date(a.modTime).getTime() - new Date(b.modTime).getTime()
        break
      default:
        comparison = a.fileName.localeCompare(b.fileName)
    }

    return order === 'asc' ? comparison : -comparison
  })
}
