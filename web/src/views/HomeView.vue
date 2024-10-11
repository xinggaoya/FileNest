<template>
  <div class="file-manager" style="height: calc(100% - 45px)">
    <n-card class="home-max-card" content-style="height: calc(100% - 120px)">
      <template v-slot:header>
        <n-image :src="Logo" width="30px" height="30px" style="vertical-align: middle" />
        File Manager
      </template>
      <template v-slot:default>
        <HomeHeader
          v-model:breadcrumb="path"
          @breadcrumb-click="handelBreadcrumbClick"
          @upload-change="handelUpload"
        />
        <n-button type="primary" size="small" @click="handelUploadClick">上传文件</n-button>

        <div class="file-list" style="height: calc(100% - 55px)">
          <n-scrollbar>
            <table style="font-size: 16px; width: 100%">
              <tr style="text-align: left">
                <th>名称</th>
                <th>大小</th>
                <th>修改时间</th>
              </tr>
              <tr
                v-for="(file, index) in files"
                :key="index"
                @click="handleFileClick(file)"
                class="file-item"
              >
                <td>
                  <n-icon>
                    <FolderOpenTwotone v-if="file.isDir" />
                    <FileTextOutlined v-else />
                  </n-icon>
                  {{ file.fileName }}
                </td>
                <td>{{ (file.fileSize / 1024 / 1024).toFixed(2) }} MB</td>
                <td>
                  <n-time :time="file.modTime"></n-time>
                </td>
              </tr>
            </table>
          </n-scrollbar>
        </div>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { FileTextOutlined, FolderOpenTwotone } from '@vicons/antd'
import Logo from '@/assets/logo.png'
import HomeHeader from '@/components/home/HomeHeader.vue'
import { getFileList, uploadFile } from '@/api/file/file'
import { useMessage } from 'naive-ui'

const path = ref(['home'])
const files = ref([])
const file = ref()
const uploading = ref(false)
const chunkSize = 5 * 1024 * 1024 // 5MB
const message = useMessage()

const getList = () => {
  // 拼接路径
  const list = [...path.value]
  list[0] = '/'
  let pathStr = list.join('/')
  getFileList({
    path: pathStr
  }).then((res: any) => {
    files.value = res.data
  })
}

const handelUpload = (f: any) => {
  file.value = f.file.file
}

const handleFileClick = (file: any) => {
  if (file.isDir) {
    path.value.push(file.fileName)
    getList()
  }
}

const handelBreadcrumbClick = (p: number) => {
  path.value.splice(p)
  getList()
}

// 上传文件
const handelUploadClick = async () => {
  if (!file.value) {
    message.warning('请选择文件')
    return
  }

  uploading.value = true
  const totalChunks = Math.ceil(file.value.size / chunkSize)

  // Start the upload process
  for (let i = 0; i < totalChunks; i++) {
    const start = i * chunkSize
    const end = Math.min(start + chunkSize, file.value.size)
    const chunk = file.value.slice(start, end)

    await uploadChunk(chunk, i, totalChunks)
  }

  uploading.value = false
  getList()
}

const uploadChunk = async (chunk, index, totalChunks) => {
  await uploadFile(file.value, {
    indexChunk: index + 1,
    totalChunks: totalChunks
  })
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.file-manager {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.home-max-card {
  border-radius: 0.7rem;
  height: 100%;
  box-shadow:
    0 10px 15px -3px rgb(0 0 0 / 0.09),
    0 4px 6px -4px rgb(0 0 0 / 0.1);
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.file-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-radius: 5px;
  font-size: 16px;
}

.file-item {
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  transition:
    background-color 0.3s ease,
    transform 0.3s ease;
}

.file-item:hover {
  background-color: #e6f7ff;
  transform: scale(1.02);
}
</style>
