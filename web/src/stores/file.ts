import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { FileInfo, Favorite } from '@/types/file'
import { createDiscreteApi } from 'naive-ui'
import { 
  getFileList, 
  createFolder, 
  downloadFile, 
  deleteFile, 
  searchFiles,
  addFavorite as addFavoriteApi,
  removeFavorite as removeFavoriteApi,
  getFavorites as getFavoritesApi,
  renameFile as renameFileApi,
  copyFile as copyFileApi,
  moveFile as moveFileApi
} from '@/api/file/file'

const { message } = createDiscreteApi(['message'])

// 文件大小格式化
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + units[i]
}

export const useFileStore = defineStore('file', () => {
  // 当前路径
  const currentPath = ref<string[]>([])
  const currentPathString = computed(() => currentPath.value.join('/'))

  // 文件列表
  const files = ref<FileInfo[]>([])
  const isLoading = ref(false)

  // 收藏列表
  const favorites = ref<Favorite[]>([])

  // 视图模式
  const viewMode = ref<'grid' | 'list'>('grid')
  const toggleViewMode = () => {
    viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid'
  }

  // 排序后的文件列表
  const sortedFiles = computed(() => {
    const sorted = [...files.value]
    // 文件夹排在前面
    return sorted.sort((a, b) => {
      if (a.isDir && !b.isDir) return -1
      if (!a.isDir && b.isDir) return 1
      return a.fileName.localeCompare(b.fileName)
    })
  })

  // 获取文件列表
  const fetchFiles = async () => {
    try {
      isLoading.value = true
      const { data } = await getFileList(currentPathString.value)
      files.value = data || []
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取文件列表失败')
      files.value = [] // 出错时清空文件列表
    } finally {
      isLoading.value = false
    }
  }

  // 获取收藏列表
  const fetchFavorites = async () => {
    try {
      const { data } = await getFavoritesApi()
      favorites.value = data || []
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取收藏列表失败')
      favorites.value = []
    }
  }

  // 添加收藏
  const addToFavorites = async (path: string) => {
    try {
      await addFavoriteApi(path)
      message.success('添加收藏成功')
      await fetchFavorites()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '添加收藏失败')
    }
  }

  // 取消收藏
  const removeFromFavorites = async (path: string) => {
    try {
      await removeFavoriteApi(path)
      message.success('取消收藏成功')
      await fetchFavorites()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '取消收藏失败')
    }
  }

  // 进入目录
  const enterDirectory = async (path: string) => {
    if (path === '/') {
      currentPath.value = []
    } else {
      // 如果是绝对路径（包含完整路径），则直接解析
      if (path.includes('/')) {
        currentPath.value = path.split('/').filter(Boolean)
      } else {
        // 如果是单个目录名，则追加到当前路径
        currentPath.value.push(path)
      }
    }
    await fetchFiles()
  }

  // 创建文件夹
  const createNewFolder = async (name: string) => {
    try {
      const path = currentPathString.value ? `${currentPathString.value}/${name}` : name
      await createFolder(path)
      message.success('创建文件夹成功')
      await fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '创建文件夹失败')
      throw error
    }
  }

  // 下载文件
  const downloadFileAction = async (path: string) => {
    try {
      message.loading('正在下载...', { duration: 0 })
      const blob = await downloadFile(path)
      
      if (!(blob instanceof Blob)) {
        throw new Error('下载失败，返回数据格式错误')
      }

      // 从路径中获取文件名
      const fileName = path.split('/').pop() || 'download'

      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.style.display = 'none'
      a.href = url
      a.download = fileName
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      
      message.success('下载成功')
    } catch (error) {
      console.error('下载错误:', error)
      message.error(error instanceof Error ? error.message : '下载文件失败')
    } finally {
      message.destroyAll() // 清除所有消息，包括 loading
    }
  }

  // 删除文件
  const deleteFileAction = async (path: string, force: boolean = false) => {
    try {
      await deleteFile(path, force)
      message.success('删除成功')
      await fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '删除文件失败')
      throw error
    }
  }

  // 删除文件（别名，为了保持 API 一致性）
  const removeFile = deleteFileAction

  // 搜索文件
  const searchFile = async (keyword: string) => {
    try {
      isLoading.value = true
      const { data } = await searchFiles(keyword)
      files.value = data || []
    } catch (error) {
      message.error(error instanceof Error ? error.message : '搜索文件失败')
      files.value = [] // 出错时清空文件列表
    } finally {
      isLoading.value = false
    }
  }

  // 返回上目录
  const goBack = async () => {
    if (currentPath.value.length > 0) {
      currentPath.value.pop()
      await fetchFiles()
    }
  }

  // 重命名文件或文件夹
  const renameFile = async (path: string, newName: string) => {
    try {
      await renameFileApi(path, newName)
      message.success('重命名成功')
      await fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '重命名失败')
      throw error
    }
  }

  // 复制文件或文件夹
  const copyFile = async (srcPath: string, destPath: string) => {
    try {
      await copyFileApi(srcPath, destPath)
      message.success('复制成功')
      await fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '复制失败')
      throw error
    }
  }

  // 移动文件或文件夹
  const moveFile = async (srcPath: string, destPath: string) => {
    try {
      await moveFileApi(srcPath, destPath)
      message.success('移动成功')
      await fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '移动失败')
      throw error
    }
  }

  return {
    currentPath,
    currentPathString,
    files,
    isLoading,
    viewMode,
    sortedFiles,
    favorites,
    fetchFiles,
    enterDirectory,
    createNewFolder,
    downloadFile: downloadFileAction,
    deleteFile: deleteFileAction,
    removeFile,
    toggleViewMode,
    searchFile,
    goBack,
    formatFileSize,
    fetchFavorites,
    addToFavorites,
    removeFromFavorites,
    renameFile,
    copyFile,
    moveFile
  }
})