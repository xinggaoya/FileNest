/**
 * 格式化文件大小
 * @param size 文件大小（字节）
 * @returns 格式化后的文件大小
 */
export function formatFileSize(size: number): string {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024

  if (size === 0) return '0 B'
  
  const i = Math.floor(Math.log(size) / Math.log(k))
  const value = size / Math.pow(k, i)
  
  // 如果是 B，不显示小数点
  if (i === 0) {
    return Math.round(value) + ' ' + units[i]
  }
  
  // 其他单位显示两位小数
  return value.toFixed(2) + ' ' + units[i]
}

/**
 * 格式化时间
 * @param time 时间字符串
 * @returns 格式化后的时间
 */
export function formatTime(time: string): string {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  // 如果是今天
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 如果是今年
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 其他情况
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
} 