<template>
  <div class="file-manager" style="height: calc(100% - 45px)">
    <n-card class="home-max-card" content-style="height: calc(100% - 120px)">
      <template v-slot:header>
        <n-flex size="small">
          <n-avatar :src="Logo" :size="40" />
          <span style="line-height: 40px">File Manager</span>
        </n-flex>
      </template>
      <template v-slot:default>
        <HomeHeader v-model:breadcrumb="path" :get-list="getList" />

        <div class="file-list" style="height: calc(100% - 55px)">
          <n-scrollbar>
            <table style="font-size: 16px; width: 100%">
              <thead>
                <tr style="text-align: left">
                  <th width="240">名称</th>
                  <th width="100">大小</th>
                  <th width="150">修改时间</th>
                </tr>
              </thead>
              <transition-group name="list" tag="tbody">
                <tr
                  v-for="(file, index) in files"
                  :key="index"
                  @click="handleFileClick(file)"
                  @contextmenu="(e) => handleContextmenu(e, file)"
                  class="file-item"
                >
                  <td>
                    <n-flex>
                      <n-icon :size="20" :component="getFileIcon(file)" />
                      <n-ellipsis style="max-width: 240px">{{ file.fileName }}</n-ellipsis>
                    </n-flex>
                  </td>
                  <td>{{ (file.fileSize / 1024 / 1024).toFixed(2) }} MB</td>
                  <td>
                    {{ file.modTime }}
                  </td>
                </tr>
              </transition-group>
            </table>
          </n-scrollbar>
        </div>
        <n-dropdown
          placement="bottom-start"
          trigger="manual"
          :options="options"
          :x="x"
          :y="y"
          :show="showDropdown"
          :on-clickoutside="onClickoutside"
          @select="handleSelect"
        />
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h, nextTick } from 'vue'
import {
  FileExcelOutlined,
  FileTextOutlined,
  FileWordOutlined,
  FolderOpenTwotone
} from '@vicons/antd'
import { ImageOutline, VideocamOutline } from '@vicons/ionicons5'
import Logo from '@/assets/logo.png'
import HomeHeader from '@/components/home/HomeHeader.vue'
import { deleteFile, getFileList } from '@/api/file/file'
import { type DropdownOption, useMessage, useDialog } from 'naive-ui'

const path = ref(['首页'])
const files = ref<any>([])
const showDropdown = ref(false)
// 选中的文件
const selectFile = ref()
const x = ref(0)
const y = ref(0)
const message = useMessage()
const dialog = useDialog()
const options: DropdownOption[] = [
  {
    label: '下载',
    key: 'download'
  },
  {
    label: () => h('span', { style: { color: 'red' } }, '删除'),
    key: 'delete'
  }
]

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

// 通过文件类型显示图标
const getFileIcon = (file: any) => {
  if (file.isDir) {
    return FolderOpenTwotone
  }
  switch (file.fileType) {
    case '.png':
      return ImageOutline
    case '.jpg':
      return ImageOutline
    case '.mp4':
      return VideocamOutline
    case '.docx':
      return FileWordOutlined
    case '.xlsx':
      return FileExcelOutlined
    case '.xls':
      return FileExcelOutlined
    default:
      return FileTextOutlined
  }
}

const handleContextmenu = (e: any, file: any) => {
  e.preventDefault()
  showDropdown.value = file
  selectFile.value = file
  nextTick().then(() => {
    showDropdown.value = true
    x.value = e.clientX
    y.value = e.clientY
  })
}

const handleSelect = (key: string) => {
  showDropdown.value = false
  if (key === 'download') {
    // a 标签
    const a = document.createElement('a')
    a.href = '/api/file/download?path=' + selectFile.value.filePath
    a.download = selectFile.value.fileName
    a.click()
    a.remove()
    return
  }
  if (key === 'delete') {
    dialog.success({
      title: '提示',
      content: '请再次确认删除该文件？',
      positiveText: '确定',
      negativeText: '取消',
      maskClosable: false,
      onPositiveClick: () => {
        deleteFile(selectFile.value.filePath).then((res: any) => {
          message.success('删除成功')
          getList()
        })
      }
    })
    return
  }
}
const onClickoutside = () => {
  showDropdown.value = false
}

const handleFileClick = (file: any) => {
  if (file.isDir) {
    path.value.push(file.fileName)
    getList()
  }
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

.file-item {
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  transition:
    background-color 0.3s ease,
    transform 0.3s ease;
}

.file-item:hover {
  background-color: #b9d5fd;
  transform: scale(1.02);
}

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
