import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  FileInfo,
  FileStats,
  Favorite,
  ViewMode,
  SortConfig,
  UploadProgress,
  BreadcrumbItem
} from '@/types/file'
import {
  getFileList,
  searchFiles,
  createFolder,
  deleteFile,
  renameFile,
  copyFile,
  moveFile,
  getFileStats,
  addFavorite,
  removeFavorite,
  getFavorites
} from '@/api/file'
import { parseBreadcrumbs, sortFiles, getParentPath } from '@/utils/file'

export const useFileStore = defineStore('file', () => {
  // 状态
  const currentPath = ref('/') // 当前目录路径
  const files = ref<FileInfo[]>([]) // 文件列表
  const isLoading = ref(false) // 加载状态
  const viewMode = ref<ViewMode>('grid') // 视图模式
  const sortConfig = ref<SortConfig>({ field: 'name', order: 'asc' }) // 排序配置
  const selectedFiles = ref<Set<string>>(new Set()) // 选中的文件
  const searchKeyword = ref('') // 搜索关键词
  const searchResults = ref<FileInfo[]>([]) // 搜索结果
  const isSearchMode = ref(false) // 是否为搜索模式
  const favorites = ref<Favorite[]>([]) // 收藏列表
  const stats = ref<FileStats>({ totalFiles: 0, totalFolders: 0, totalSize: 0 }) // 统计信息
  const uploadQueue = ref<Map<string, UploadProgress>>(new Map()) // 上传队列
  const copiedFiles = ref<string[]>([]) // 复制的文件列表
  const isCutMode = ref(false) // 是否为剪切模式

  // 计算属性
  const displayFiles = computed(() => {
    const targetFiles = isSearchMode.value ? searchResults.value : files.value
    return sortFiles(targetFiles, sortConfig.value.field, sortConfig.value.order)
  })

  const breadcrumbs = computed<BreadcrumbItem[]>(() => {
    if (isSearchMode.value) {
      return [{ label: '搜索结果', path: '' }]
    }
    return parseBreadcrumbs(currentPath.value)
  })

  const hasParent = computed(() => {
    return !isSearchMode.value && currentPath.value !== '/'
  })

  const selectedFilesArray = computed(() => {
    return displayFiles.value.filter((file) => selectedFiles.value.has(file.filePath))
  })

  const canGoBack = computed(() => {
    return hasParent.value || isSearchMode.value
  })

  const hasFilesToPaste = computed(() => copiedFiles.value.length > 0)

  // 操作方法

  /**
   * 获取文件列表
   */
  const fetchFiles = async (path?: string) => {
    try {
      isLoading.value = true
      const targetPath = path ?? currentPath.value

      const response = await getFileList(targetPath)
      files.value = response.data.data

      if (path !== undefined) {
        currentPath.value = targetPath
      }

      // 同时获取统计信息
      await fetchStats(targetPath)
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '获取文件列表失败'
      }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 进入目录
   */
  const enterDirectory = async (dirPath: string) => {
    if (isSearchMode.value) {
      // 退出搜索模式
      exitSearchMode()
    }
    // 清空当前文件列表
    files.value = []
    // 更新当前路径
    currentPath.value = dirPath
    // 获取新目录的文件列表
    await fetchFiles(dirPath)
  }

  /**
   * 返回上级目录
   */
  const goBack = async () => {
    if (isSearchMode.value) {
      exitSearchMode()
    } else if (hasParent.value) {
      const parentPath = getParentPath(currentPath.value)
      await fetchFiles(parentPath)
    }
  }

  /**
   * 导航到指定路径
   */
  const navigateTo = async (path: string) => {
    if (isSearchMode.value) {
      exitSearchMode()
    }
    await fetchFiles(path)
  }

  /**
   * 搜索文件
   */
  const searchFile = async (keyword: string) => {
    if (!keyword.trim()) {
      exitSearchMode()
      return { success: true }
    }

    try {
      isLoading.value = true
      searchKeyword.value = keyword

      const response = await searchFiles(keyword)
      searchResults.value = response.data.data
      isSearchMode.value = true
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '搜索失败'
      }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 退出搜索模式
   */
  const exitSearchMode = () => {
    isSearchMode.value = false
    searchKeyword.value = ''
    searchResults.value = []
  }

  /**
   * 创建新文件夹
   */
  const createNewFolder = async (name: string) => {
    try {
      const folderPath = currentPath.value === '/' ? `/${name}` : `${currentPath.value}/${name}`
      await createFolder(folderPath)
      await fetchFiles()
      return { success: true, message: '文件夹创建成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '创建文件夹失败'
      }
    }
  }

  /**
   * 删除文件或文件夹
   */
  const deleteFiles = async (paths: string[], force = false) => {
    try {
      await Promise.all(paths.map((path) => deleteFile(path, force)))
      clearSelection()
      await fetchFiles()
      return { success: true, message: '删除成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '删除失败'
      }
    }
  }

  /**
   * 重命名文件或文件夹
   */
  const renameItem = async (path: string, newName: string) => {
    try {
      await renameFile(path, newName)
      await fetchFiles()
      return { success: true, message: '重命名成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '重命名失败'
      }
    }
  }

  /**
   * 复制文件或文件夹
   */
  const copyFiles = async (srcPaths: string[], destPath: string) => {
    try {
      await Promise.all(srcPaths.map((srcPath) => copyFile(srcPath, destPath)))
      await fetchFiles()
      return { success: true, message: '复制成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '复制失败'
      }
    }
  }

  /**
   * 移动文件或文件夹
   */
  const moveFiles = async (srcPaths: string[], destPath: string) => {
    try {
      await Promise.all(srcPaths.map((srcPath) => moveFile(srcPath, destPath)))
      clearSelection()
      await fetchFiles()
      return { success: true, message: '移动成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '移动失败'
      }
    }
  }

  /**
   * 获取统计信息
   */
  const fetchStats = async (path?: string) => {
    try {
      const targetPath = path ?? currentPath.value
      const response = await getFileStats(targetPath)
      stats.value = response.data.data
    } catch (error) {
      console.error('获取统计信息失败:', error)
    }
  }

  /**
   * 获取收藏列表
   */
  const fetchFavorites = async () => {
    try {
      const response = await getFavorites()
      favorites.value = response.data.data
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '获取收藏列表失败'
      }
    }
  }

  /**
   * 添加到收藏
   */
  const addToFavorites = async (path: string) => {
    try {
      await addFavorite(path)
      await fetchFavorites()
      return { success: true, message: '已添加到收藏' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '添加收藏失败'
      }
    }
  }

  /**
   * 移除收藏
   */
  const removeFromFavorites = async (path: string) => {
    try {
      await removeFavorite(path)
      await fetchFavorites()
      return { success: true, message: '已移除收藏' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '移除收藏失败'
      }
    }
  }

  /**
   * 切换视图模式
   */
  const toggleViewMode = () => {
    viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid'
  }

  /**
   * 设置排序配置
   */
  const setSortConfig = (field: string, order?: 'asc' | 'desc') => {
    if (sortConfig.value.field === field) {
      // 如果是同一字段，切换排序方向
      sortConfig.value.order = sortConfig.value.order === 'asc' ? 'desc' : 'asc'
    } else {
      // 不同字段，使用指定方向或默认升序
      sortConfig.value.field = field as any
      sortConfig.value.order = order || 'asc'
    }
  }

  /**
   * 选择文件
   */
  const selectFile = (filePath: string) => {
    selectedFiles.value.add(filePath)
  }

  /**
   * 取消选择文件
   */
  const unselectFile = (filePath: string) => {
    selectedFiles.value.delete(filePath)
  }

  /**
   * 切换文件选择状态
   */
  const toggleFileSelection = (filePath: string) => {
    if (selectedFiles.value.has(filePath)) {
      unselectFile(filePath)
    } else {
      selectFile(filePath)
    }
  }

  /**
   * 选择所有文件
   */
  const selectAllFiles = () => {
    displayFiles.value.forEach((file) => {
      selectedFiles.value.add(file.filePath)
    })
  }

  /**
   * 清除选择
   */
  const clearSelection = () => {
    selectedFiles.value.clear()
  }

  /**
   * 检查文件是否已选择
   */
  const isFileSelected = (filePath: string) => {
    return selectedFiles.value.has(filePath)
  }

  /**
   * 刷新当前目录
   */
  const refresh = async () => {
    if (isSearchMode.value) {
      await searchFile(searchKeyword.value)
    } else {
      await fetchFiles()
    }
  }

  /**
   * 复制文件
   */
  const copySelectedFiles = () => {
    copiedFiles.value = Array.from(selectedFiles.value)
    isCutMode.value = false
  }

  /**
   * 剪切文件
   */
  const cutSelectedFiles = () => {
    copiedFiles.value = Array.from(selectedFiles.value)
    isCutMode.value = true
  }

  /**
   * 粘贴文件
   */
  const pasteFiles = async () => {
    if (copiedFiles.value.length === 0) return

    try {
      if (isCutMode.value) {
        // 移动文件
        await moveFiles(copiedFiles.value, currentPath.value)
      } else {
        // 复制文件
        await copyFiles(copiedFiles.value, currentPath.value)
      }
      // 清空复制列表
      copiedFiles.value = []
      isCutMode.value = false
      return { success: true, message: isCutMode.value ? '移动成功' : '复制成功' }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : isCutMode.value ? '移动失败' : '复制失败'
      }
    }
  }

  /**
   * 获取文件类型
   */
  const getFileType = (fileName: string): string => {
    const ext = fileName.split('.').pop()?.toLowerCase() || ''
    return ext
  }

  /**
   * 判断是否为可预览的文本文件
   */
  const isPreviewableText = (fileName: string): boolean => {
    const textExtensions = ['txt', 'md', 'json', 'js', 'ts', 'html', 'css', 'xml', 'yaml', 'yml']
    return textExtensions.includes(getFileType(fileName))
  }

  /**
   * 判断是否为图片文件
   */
  const isImage = (fileName: string): boolean => {
    const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg']
    return imageExtensions.includes(getFileType(fileName))
  }

  /**
   * 下载文件（无类型判断，直接下载）
   */
  const downloadFile = async (filePath: string, fileName: string) => {
    try {
      const response = await fetch(`/api/file/download?path=${encodeURIComponent(filePath)}`)
      if (!response.ok) throw new Error('下载失败')
      const contentDisposition = response.headers.get('content-disposition')
      let downloadFileName = fileName
      if (contentDisposition) {
        const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition)
        if (matches != null && matches[1]) {
          downloadFileName = matches[1].replace(/['"]/g, '')
        }
      }
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = downloadFileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('下载文件失败:', error)
      throw error
    }
  }

  /**
   * 处理文件下载或预览（只做预览）
   */
  const handleFileAction = async (filePath: string, fileName: string) => {
    if (isPreviewableText(fileName) || isImage(fileName)) {
      window.open(`/api/file/preview?path=${encodeURIComponent(filePath)}`, '_blank')
      return
    }
    // 其他类型可提示不支持预览
  }

  return {
    // 状态
    currentPath,
    files,
    isLoading,
    viewMode,
    sortConfig,
    selectedFiles,
    searchKeyword,
    searchResults,
    isSearchMode,
    favorites,
    stats,
    uploadQueue,
    copiedFiles,
    isCutMode,

    // 计算属性
    displayFiles,
    breadcrumbs,
    hasParent,
    selectedFilesArray,
    canGoBack,
    hasFilesToPaste,

    // 方法
    fetchFiles,
    enterDirectory,
    goBack,
    navigateTo,
    searchFile,
    exitSearchMode,
    createNewFolder,
    deleteFiles,
    renameItem,
    copyFiles,
    moveFiles,
    fetchStats,
    fetchFavorites,
    addToFavorites,
    removeFromFavorites,
    toggleViewMode,
    setSortConfig,
    selectFile,
    unselectFile,
    toggleFileSelection,
    selectAllFiles,
    clearSelection,
    isFileSelected,
    refresh,
    copySelectedFiles,
    cutSelectedFiles,
    pasteFiles,
    getFileType,
    isPreviewableText,
    isImage,
    downloadFile,
    handleFileAction
  }
})
