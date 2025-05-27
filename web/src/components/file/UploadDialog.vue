<template>
  <div class="upload-dialog">
    <div class="upload-area" @drop="handleDrop" @dragover.prevent @dragenter.prevent>
      <n-icon size="48" class="upload-icon">
        <CloudUploadIcon />
      </n-icon>
      <div class="upload-text">拖拽文件到此处上传</div>
      <div class="upload-hint">或者点击选择文件</div>
      <n-button type="primary" @click="selectFiles"> 选择文件 </n-button>
    </div>

    <input ref="fileInput" type="file" multiple style="display: none" @change="handleFileSelect" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NIcon, useMessage } from 'naive-ui'
import { CloudUpload as CloudUploadIcon } from '@vicons/ionicons5'

const emit = defineEmits<{
  close: []
}>()

const message = useMessage()
const fileInput = ref<HTMLInputElement>()

const selectFiles = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    handleFiles(target.files)
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer?.files) {
    handleFiles(event.dataTransfer.files)
  }
}

const handleFiles = (files: FileList) => {
  message.info(`选择了 ${files.length} 个文件`)
  // TODO: 实现文件上传逻辑
  console.log('上传文件:', files)
}
</script>

<style scoped>
.upload-dialog {
  padding: 20px;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  background: #fafafa;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #1890ff;
  background: #f0f8ff;
}

.upload-icon {
  color: #999;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 16px;
  margin-bottom: 8px;
  color: #333;
}

.upload-hint {
  font-size: 12px;
  color: #999;
  margin-bottom: 20px;
}
</style>
