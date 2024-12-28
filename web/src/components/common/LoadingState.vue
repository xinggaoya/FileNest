<template>
  <div v-if="loading" class="loading-container">
    <n-spin :size="size" :description="description">
      <template #icon>
        <n-icon class="loading-icon" :size="iconSize">
          <component :is="icon" />
        </n-icon>
      </template>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CloudDownloadOutline } from '@vicons/ionicons5'

const props = withDefaults(
  defineProps<{
    loading?: boolean
    description?: string
    size?: 'small' | 'medium' | 'large'
    icon?: any
  }>(),
  {
    loading: false,
    description: '加载中...',
    size: 'large',
    icon: CloudDownloadOutline
  }
)

const iconSize = computed(() => {
  switch (props.size) {
    case 'small':
      return 24
    case 'medium':
      return 36
    case 'large':
      return 48
    default:
      return 36
  }
})
</script>

<style scoped>
.loading-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.loading-icon {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

:deep(.n-spin-description) {
  margin-top: 12px;
  font-size: 14px;
  color: var(--n-text-color);
}
</style> 