<template>
  <div class="file-stats">
    <n-space>
      <n-statistic label="文件数" :value="stats.totalFiles" />
      <n-statistic label="文件夹数" :value="stats.totalFolders" />
      <n-statistic label="总大小" :value="formattedTotalSize" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { formatFileSize } from '@/utils/format'
import { useFileStore } from '@/stores/file'
import { getFileStats } from '@/api/file/file'
import type { FileStats } from '@/types/file'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])
const fileStore = useFileStore()

// 统计数据
const stats = ref<FileStats>({
  totalFiles: 0,
  totalFolders: 0,
  totalSize: 0
})

// 格式化后的总大小
const formattedTotalSize = computed(() => {
  return formatFileSize(stats.value.totalSize)
})

// 获取统计数据
const fetchStats = async () => {
  try {
    const { data } = await getFileStats({ path: fileStore.currentPathString })
    if (data) {
      stats.value = data
    }
  } catch (error) {
    message.error('获取统计信息失败')
  }
}

// 监听路径变化，重新获取统计数据
watch(() => fileStore.currentPathString, () => {
  fetchStats()
})

// 初始化时获取统计数据
onMounted(() => {
  fetchStats()
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