<template>
  <div>
    <n-flex class="search-bar" justify="space-between">
      <div>
        <n-breadcrumb>
          <n-breadcrumb-item
            v-for="(item, index) in breadcrumb"
            :key="index"
            @click="handelBreadcrumbClick(item)"
          >
            <n-icon :component="FolderOpenTwotone" />
            {{ item }}
          </n-breadcrumb-item>
        </n-breadcrumb>
      </div>
      <n-flex>
        <n-button type="primary" size="small" @click="handelBreadcrumbClick('')">
          <n-icon :component="FolderOpenTwotone" />
          新建文件夹
        </n-button>
        <div>
          <n-upload action="#" :default-upload="false" @change="handleUpload" :show-file-list="false">
            <n-button type="primary" size="small" @click="handleUploadClick">
              <n-icon :component="FolderOpenTwotone" />
              上传文件
            </n-button>
          </n-upload>
        </div>
      </n-flex>
    </n-flex>
  </div>
</template>

<script setup lang="ts">
import { FolderOpenTwotone } from '@vicons/antd'

const breadcrumb = defineModel('breadcrumb')

const emit = defineEmits(['breadcrumb-click', 'upload-click', 'upload-change'])

const handelBreadcrumbClick = (item: string) => {
  emit('breadcrumb-click', item)
}

const handleUploadClick = () => {
  emit('upload-click')
}
const handleUpload = (file) => {
  emit('upload-change', file)
}
</script>

<style scoped>
.search-bar {
  margin-bottom: 10px;
}
</style>
