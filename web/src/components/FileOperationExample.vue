<template>
  <div class="file-operation-example">
    <h3>文件操作示例</h3>

    <!-- 创建文件夹 -->
    <div class="operation-section">
      <n-input
        v-model:value="newFolderName"
        placeholder="输入文件夹名称"
        @keyup.enter="handleCreateFolder"
      />
      <n-button @click="handleCreateFolder" :loading="loading"> 创建文件夹 </n-button>
    </div>

    <!-- 删除选中文件 -->
    <div class="operation-section">
      <n-button
        type="error"
        @click="handleDeleteSelected"
        :disabled="fileStore.selectedFilesArray.length === 0"
        :loading="loading"
      >
        删除选中文件 ({{ fileStore.selectedFilesArray.length }})
      </n-button>
    </div>

    <!-- 搜索文件 -->
    <div class="operation-section">
      <n-input
        v-model:value="searchKeyword"
        placeholder="搜索文件..."
        @keyup.enter="handleSearch"
      />
      <n-button @click="handleSearch" :loading="loading"> 搜索 </n-button>
      <n-button @click="handleClearSearch" v-if="fileStore.isSearchMode"> 清除搜索 </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFileStore } from '@/stores/file'
import { useMessageHandler } from '@/composables/useMessageHandler'

const { handleResult } = useMessageHandler()
const fileStore = useFileStore()

const loading = ref(false)
const newFolderName = ref('')
const searchKeyword = ref('')

// 创建文件夹
const handleCreateFolder = async () => {
  if (!newFolderName.value.trim()) {
    return
  }

  loading.value = true
  const result = await fileStore.createNewFolder(newFolderName.value.trim())
  handleResult(result)

  if (result.success) {
    newFolderName.value = ''
  }
  loading.value = false
}

// 删除选中文件
const handleDeleteSelected = async () => {
  const selectedPaths = fileStore.selectedFilesArray.map((file) => file.filePath)
  if (selectedPaths.length === 0) {
    return
  }

  loading.value = true
  const result = await fileStore.deleteFiles(selectedPaths)
  handleResult(result)
  loading.value = false
}

// 搜索文件
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    return
  }

  loading.value = true
  const result = await fileStore.searchFile(searchKeyword.value.trim())
  if (result && !result.success) {
    handleResult(result)
  }
  loading.value = false
}

// 清除搜索
const handleClearSearch = () => {
  fileStore.exitSearchMode()
  searchKeyword.value = ''
}
</script>

<style scoped>
.file-operation-example {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.operation-section {
  display: flex;
  gap: 8px;
  align-items: center;
}

.operation-section .n-input {
  flex: 1;
  max-width: 300px;
}
</style>
