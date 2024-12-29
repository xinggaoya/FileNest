<template>
  <div class="file-stats">
    <div class="stats-container">
      <n-card class="stat-card" :bordered="false" size="small">
        <div class="stat-content">
          <div class="stat-icon-wrapper">
            <n-icon size="20" class="stat-icon">
              <document-text-outline />
            </n-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">文件数</div>
            <div class="stat-value">
              <n-number-animation :from="0" :to="stats.totalFiles" :duration="1000" />
            </div>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false" size="small">
        <div class="stat-content">
          <div class="stat-icon-wrapper">
            <n-icon size="20" class="stat-icon">
              <folder-outline />
            </n-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">文件夹数</div>
            <div class="stat-value">
              <n-number-animation :from="0" :to="stats.totalFolders" :duration="1000" />
            </div>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false" size="small">
        <div class="stat-content">
          <div class="stat-icon-wrapper">
            <n-icon size="20" class="stat-icon">
              <server-outline />
            </n-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总大小</div>
            <div class="stat-value">{{ formattedTotalSize }}</div>
          </div>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { formatFileSize } from '@/utils/format'
import { useFileStore } from '@/stores/file'
import { getFileStats } from '@/api/file/file'
import type { FileStats } from '@/types/file'
import { createDiscreteApi } from 'naive-ui'
import {
  DocumentTextOutline,
  FolderOutline,
  ServerOutline
} from '@vicons/ionicons5'

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
  margin: 12px 0;
}

.stats-container {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.stat-card {
  flex: 1;
  min-width: 180px;
  transition: all 0.3s ease;
  background: linear-gradient(145deg, var(--n-card-color), var(--n-card-color-hover));
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--n-primary-color-suppl);
}

.stat-icon {
  color: var(--n-primary-color);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: var(--n-text-color-2);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--n-text-color);
}

@media (max-width: 768px) {
  .stats-container {
    gap: 8px;
  }

  .stat-card {
    min-width: 140px;
  }

  .stat-value {
    font-size: 16px;
  }
}
</style> 