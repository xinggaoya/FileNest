<template>
  <n-modal
    v-model:show="visible"
    title="上传配置"
    preset="dialog"
    :mask-closable="false"
    style="width: 500px"
    @close="handleClose"
  >
    <n-form :model="form" label-placement="left" label-width="140">
      <n-form-item label="分块大小">
        <n-input-number
          v-model:value="form.chunkSize"
          :min="1"
          :max="100"
          :step="1"
          style="width: 200px"
        />
        <n-text class="ml-2">MB</n-text>
      </n-form-item>

      <n-form-item label="最大并发上传数">
        <n-input-number
          v-model:value="form.maxConcurrent"
          :min="1"
          :max="10"
          :step="1"
          style="width: 200px"
        />
      </n-form-item>

      <n-form-item label="启用分块上传">
        <n-switch v-model:value="form.enableChunked" />
      </n-form-item>

      <n-form-item label="分块上传阈值">
        <n-input-number
          v-model:value="form.chunkThreshold"
          :min="1"
          :max="100"
          :step="1"
          style="width: 200px"
        />
        <n-text class="ml-2">MB</n-text>
      </n-form-item>
    </n-form>

    <template #action>
      <n-button @click="visible = false">取消</n-button>
      <n-button type="primary" @click="handleSave">保存</n-button>
      <n-button @click="handleReset">恢复默认</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { uploadConfig, updateUploadConfig } from '@/config/upload'
import { useMessage } from 'naive-ui'

const visible = ref(false)
const message = useMessage()

// 将字节转换为MB
const bytesToMB = (bytes: number) => Math.round(bytes / (1024 * 1024))
const MBToBytes = (mb: number) => mb * 1024 * 1024

// 表单数据
const form = reactive({
  chunkSize: bytesToMB(uploadConfig.chunkSize),
  maxConcurrent: uploadConfig.maxConcurrent,
  enableChunked: uploadConfig.enableChunked,
  chunkThreshold: bytesToMB(uploadConfig.chunkThreshold)
})

// 保存配置
const handleSave = () => {
  updateUploadConfig({
    chunkSize: MBToBytes(form.chunkSize),
    maxConcurrent: form.maxConcurrent,
    enableChunked: form.enableChunked,
    chunkThreshold: MBToBytes(form.chunkThreshold)
  })
  message.success('配置已保存')
  visible.value = false
}

// 重置为默认配置
const handleReset = () => {
  form.chunkSize = 2
  form.maxConcurrent = 3
  form.enableChunked = true
  form.chunkThreshold = 2
  handleSave()
}

// 关闭对话框时重置表单
const handleClose = () => {
  form.chunkSize = bytesToMB(uploadConfig.chunkSize)
  form.maxConcurrent = uploadConfig.maxConcurrent
  form.enableChunked = uploadConfig.enableChunked
  form.chunkThreshold = bytesToMB(uploadConfig.chunkThreshold)
}

watch(visible, (newValue) => {
  if (!newValue) {
    handleClose()
  }
})

// 暴露方法给父组件
defineExpose({
  show: () => {
    visible.value = true
  }
})
</script>

<style scoped>
.ml-2 {
  margin-left: 8px;
}
</style>
