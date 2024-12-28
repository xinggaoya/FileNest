<template>
  <div class="file-stats">
    <n-space>
      <n-statistic label="文件数" :value="fileCount" />
      <n-statistic label="文件夹数" :value="folderCount" />
      <n-statistic label="总大小" :value="formattedTotalSize" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatFileSize } from '@/utils/format'
import { useFileStore } from '@/stores/file'
import type { FileInfo } from '@/types/file'

const fileStore = useFileStore()

// 计算文件数量
const fileCount = computed(() => {
  return fileStore.files.filter((file: FileInfo) => !file.isDir).length
})

// 计算文件夹数量
const folderCount = computed(() => {
  return fileStore.files.filter((file: FileInfo) => file.isDir).length
})

// 计算总大小
const totalSize = computed(() => {
  return fileStore.files.reduce((total: number, file: FileInfo) => {
    return total + (file.isDir ? 0 : file.fileSize)
  }, 0)
})

// 格式化后的总大小
const formattedTotalSize = computed(() => {
  return formatFileSize(totalSize.value)
})
</script>

<style scoped>
.file-stats {
  margin-top: 16px;
  padding: 16px;
  background-color: var(--n-card-color);
  border-radius: 8px;
}
</style> 