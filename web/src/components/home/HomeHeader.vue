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
        <n-button type="primary" size="small" @click="handelCreateFolder">
          <template v-slot:icon>
            <n-icon :component="CreateOutline" />
          </template>
          新建文件夹
        </n-button>
        <n-button type="primary" size="small" @click="handleUploadPanClick">
          <template v-slot:icon>
            <n-icon :component="CloudUploadOutline" />
          </template>
          上传文件
        </n-button>
      </n-flex>
    </n-flex>
    <div>
      <n-drawer v-model:show="drawerShow" :mask-closable="false" :width="502">
        <n-drawer-content closable>
          <template v-slot:header>
            <n-icon :component="FolderOpenTwotone" :size="20" />
            <n-text>上传文件</n-text>
          </template>
          <template v-slot:default>
            <n-flex vertical>
              <div>
                <n-flex justify="space-between">
                  <n-form-item label="目录">
                    <n-highlight
                      :text="'当前目录：' + breadcrumb?.join('/')"
                      :patterns="[breadcrumb?.join('/')]"
                    />
                  </n-form-item>
                  <n-form-item label="上传类型">
                    <n-radio-group
                      v-model:value="uploadType"
                      size="small"
                      style="margin-bottom: 10px"
                    >
                      <n-radio
                        v-for="(item, index) in ['文件', '文件夹']"
                        :key="index"
                        :value="item"
                        >{{ item }}
                      </n-radio>
                    </n-radio-group>
                  </n-form-item>
                </n-flex>
                <n-form-item label="分块大小">
                  <n-input-number :min="1" :max="100" v-model:value="chunkSizeIndex" />
                </n-form-item>
              </div>
              <n-upload
                multiple
                :directory="uploadType === '文件夹'"
                action="#"
                :default-upload="false"
                v-model:file-list="fileList"
              >
                <n-upload-dragger>
                  <div style="margin-bottom: 12px">
                    <n-icon size="48" :depth="3">
                      <ArchiveIcon />
                    </n-icon>
                  </div>
                  <n-text style="font-size: 16px"> 点击或者拖动文件到该区域来上传</n-text>
                  <n-p depth="3" style="margin: 8px 0 0 0">
                    请不要上传敏感数据，比如你的银行卡号和密码，信用卡号有效期和安全码
                  </n-p>
                </n-upload-dragger>
                <div>
                  <n-text depth="3">下载速度：{{ currentUploadSpeed.toFixed(2) }} MB/s</n-text>
                </div>
              </n-upload>
            </n-flex>
          </template>
          <template v-slot:footer>
            <n-button type="primary" :loading="uploading" @click="handleUploadClick">上传</n-button>
          </template>
        </n-drawer-content>
      </n-drawer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { FolderOpenTwotone } from '@vicons/antd'
import { ArchiveOutline as ArchiveIcon, CloudUploadOutline, CreateOutline } from '@vicons/ionicons5'
import { NButton, NInput, type UploadFileInfo, useMessage, useModal } from 'naive-ui'
import { createFolder, uploadFile } from '@/api/file/file'

const modal = useModal()
const drawerShow = ref(false)
const uploading = ref(false)
const chunkSizeIndex = ref(90)
const chunkSize = computed(() => {
  return chunkSizeIndex.value * 1024 * 1024
})
const message = useMessage()
const fileList = ref<UploadFileInfo[]>([])
const uploadType = ref('文件')
const dirValue = ref('')
const breadcrumb = defineModel<Array<string>>('breadcrumb', {
  default: []
})

// 上传速度
let currentUploadSpeed = ref(0)

const props = defineProps({
  // 刷新
  getList: {
    type: Function,
    default: () => {}
  }
})

const handelBreadcrumbClick = (item: string) => {
  breadcrumb.value = breadcrumb.value.slice(0, breadcrumb.value.indexOf(item) + 1)
  props.getList()
}

// 新建文件夹
const handelCreateFolder = () => {
  const m = modal.create({
    title: '新建文件夹',
    preset: 'card',
    style: {
      width: '500px'
    },
    content: () =>
      h(
        NInput,
        {
          placeholder: '请输入文件夹名称',
          autofocus: true,
          onUpdateValue: (value: string) => {
            dirValue.value = value
          }
        },
        {}
      ),
    action: () =>
      h(
        NButton,
        {
          type: 'primary',
          onClick: () => {
            const value = dirValue.value
            createNewFolder(value)
            m.destroy()
          }
        },
        () => '提交'
      )
  })
}

// 创建文件夹
const createNewFolder = (name: string) => {
  const pathStr = [...breadcrumb.value]
  // 移除第一个
  pathStr.shift()
  pathStr.push(name)
  createFolder(pathStr.join('/')).then((res: any) => {
    props.getList()
  })
}

// 上传文件
const handelUpload = async (file: any, fullPath: any, indexFile: number) => {
  if (!file) {
    message.warning('请选择文件')
    return
  }

  const totalChunks = Math.ceil(file.size / chunkSize.value)
  const fileName = fullPath

  // 文件总大小
  let fileTotalSize = 0
  // 当前上传文件大小
  let currentUploadSize = 0

  // 开始下载时间
  const startTime = Date.now()
  // 存储总片大小
  fileTotalSize = file.size
  // 上传文件分片
  const uploadChunk = async (chunk, index) => {
    // 保存路径
    const path = [...breadcrumb.value]
    path.shift()

    // 进度条
    const progress = (event: any) => {
      currentUploadSize += event.bytes
      // 计算速度 mb 保留两位小数
      currentUploadSpeed.value =
        Math.round((event.bytes / (Date.now() - startTime)) * 1000) / 1024 / 1024
      const file = fileList.value[indexFile]
      file.status = 'uploading'
      file.percentage = Math.round((currentUploadSize / fileTotalSize) * 100)
      if (currentUploadSize >= fileTotalSize) {
        file.status = 'finished'
      }
    }

    await uploadFile(
      chunk,
      {
        indexChunk: index + 1,
        totalChunks: totalChunks,
        fileName: fileName,
        path: path.join('/')
      },
      progress
    ).then((res: any) => {
      // console.log(res)
    })
  }

  // Start the upload process
  for (let i = 0; i < totalChunks; i++) {
    const start = i * chunkSize.value
    const end = Math.min(start + chunkSize.value, file.size)
    const chunk = file.slice(start, end)

    await uploadChunk(chunk, i)
  }

  props.getList()
  currentUploadSpeed.value = 0
}

const handleUploadPanClick = () => {
  drawerShow.value = true
}
const handleUploadClick = async () => {
  if (fileList.value.length === 0) {
    message.warning('请选择需要上传的文件')
    return
  }
  uploading.value = true
  for (let i = 0; i < fileList.value.length; i++) {
    if (fileList.value[i].status !== 'finished') {
      await handelUpload(fileList.value[i].file, fileList.value[i].fullPath, i)
    }
  }

  message.success('上传成功')
  uploading.value = false
  // 清理已完成任务
  fileList.value = fileList.value.filter((item: any) => item.status !== 'finished')
}
</script>

<style scoped>
.search-bar {
  margin-bottom: 10px;
}
</style>
