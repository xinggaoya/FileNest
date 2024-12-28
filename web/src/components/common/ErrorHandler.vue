<template>
  <div v-if="error" class="error-container">
    <n-result :status="error.status" :title="error.title" :description="error.message">
      <template #footer>
        <n-space>
          <n-button @click="retry" v-if="error.retryable">
            重试
          </n-button>
          <n-button @click="goHome">
            返回首页
          </n-button>
        </n-space>
      </template>
    </n-result>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

interface ErrorState {
  status: 403 | 404 | 418 | 500 | '403' | '404' | '418' | '500'
  title: string
  message: string
  retryable?: boolean
}

const error = ref<ErrorState | null>(null)
const router = useRouter()

const setError = (err: ErrorState) => {
  error.value = err
}

const clearError = () => {
  error.value = null
}

const retry = () => {
  if (error.value?.retryable) {
    clearError()
    // 触发重试事件
    emit('retry')
  }
}

const goHome = () => {
  clearError()
  router.push('/')
}

// 暴露方法给父组件
defineExpose({
  setError,
  clearError
})

const emit = defineEmits<{
  (e: 'retry'): void
}>()
</script>

<style scoped>
.error-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
</style> 