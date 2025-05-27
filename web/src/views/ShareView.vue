<template>
  <div class="share-page">
    <div class="share-container">
      <!-- 分享信息展示 -->
      <div v-if="shareInfo" class="share-info">
        <div class="file-icon">
          <n-icon size="48" :color="shareInfo.isDir ? '#ffa940' : '#1890ff'">
            <FolderOutlined v-if="shareInfo.isDir" />
            <FileOutlined v-else />
          </n-icon>
        </div>
        
        <div class="file-details">
          <h2 class="file-name">{{ shareInfo.fileName }}</h2>
          <div class="file-meta">
            <span class="file-size">{{ formatFileSize(shareInfo.fileSize) }}</span>
            <span class="file-type">{{ shareInfo.isDir ? '文件夹' : '文件' }}</span>
          </div>
          
          <div v-if="shareInfo.description" class="file-description">
            {{ shareInfo.description }}
          </div>
          
          <div class="share-meta">
            <div class="meta-item">
              <span class="label">下载次数：</span>
              <span class="value">{{ shareInfo.downloaded }}</span>
              <span v-if="shareInfo.maxDownload > 0" class="limit">
                / {{ shareInfo.maxDownload }}
              </span>
            </div>
            
            <div v-if="shareInfo.expireTime" class="meta-item">
              <span class="label">过期时间：</span>
              <span class="value">{{ formatTime(shareInfo.expireTime) }}</span>
            </div>
            
            <div class="meta-item">
              <span class="label">分享时间：</span>
              <span class="value">{{ formatTime(shareInfo.createTime) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 密码输入 -->
      <div v-if="shareInfo?.hasPassword && !passwordValidated" class="password-section">
        <n-form @submit.prevent="handleValidate">
          <n-form-item label="提取密码">
            <n-input
              v-model:value="password"
              type="password"
              placeholder="请输入提取密码"
              :disabled="validating"
              @keyup.enter="handleValidate"
            />
          </n-form-item>
          <n-form-item>
            <n-button
              type="primary"
              :loading="validating"
              @click="handleValidate"
              block
            >
              验证密码
            </n-button>
          </n-form-item>
        </n-form>
      </div>
      
      <!-- 下载按钮 -->
      <div v-if="canDownload" class="download-section">
        <n-button
          type="primary"
          size="large"
          :loading="downloading"
          @click="handleDownload"
          block
        >
          <template #icon>
            <n-icon><DownloadOutlined /></n-icon>
          </template>
          下载{{ shareInfo?.isDir ? '文件夹' : '文件' }}
        </n-button>
      </div>
      
      <!-- 错误信息 -->
      <div v-if="error" class="error-section">
        <n-alert type="error" :title="error" />
      </div>
      
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-section">
        <n-spin size="large" />
        <p>正在加载分享信息...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { NIcon, NButton, NForm, NFormItem, NInput, NAlert, NSpin, useMessage } from 'naive-ui'
import { FolderOutlined, FileOutlined, DownloadOutlined } from '@vicons/antd'
import { getShareInfo, validateShare, downloadSharedFile } from '@/api/share/share'
import type { ShareInfo } from '@/api/share/share'

const route = useRoute()
const message = useMessage()

// 状态
const shareInfo = ref<ShareInfo | null>(null)
const password = ref('')
const passwordValidated = ref(false)
const loading = ref(false)
const validating = ref(false)
const downloading = ref(false)
const error = ref('')

// 计算属性
const canDownload = computed(() => {
  if (!shareInfo.value) return false
  if (shareInfo.value.hasPassword && !passwordValidated.value) return false
  return true
})

// 获取分享码
const shareCode = computed(() => route.params.shareCode as string)

// 初始化
onMounted(() => {
  if (shareCode.value) {
    loadShareInfo()
  } else {
    error.value = '分享码无效'
  }
})

// 加载分享信息
const loadShareInfo = async () => {
  try {
    loading.value = true
    error.value = ''
    
    const { data } = await getShareInfo(shareCode.value)
    shareInfo.value = data
    
    // 如果没有密码，则直接可以下载
    if (!data.hasPassword) {
      passwordValidated.value = true
    }
  } catch (err: any) {
    error.value = err.message || '获取分享信息失败'
  } finally {
    loading.value = false
  }
}

// 验证密码
const handleValidate = async () => {
  if (!password.value.trim()) {
    message.warning('请输入提取密码')
    return
  }
  
  try {
    validating.value = true
    await validateShare(shareCode.value, password.value)
    passwordValidated.value = true
    message.success('密码验证成功')
  } catch (err: any) {
    message.error(err.message || '密码验证失败')
  } finally {
    validating.value = false
  }
}

// 下载文件
const handleDownload = async () => {
  try {
    downloading.value = true
    
    const blob = await downloadSharedFile(
      shareCode.value,
      shareInfo.value?.hasPassword ? password.value : undefined
    )
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.style.display = 'none'
    a.href = url
    a.download = shareInfo.value?.fileName || 'download'
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
    
    message.success('下载成功')
  } catch (err: any) {
    message.error(err.message || '下载失败')
  } finally {
    downloading.value = false
  }
}

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + units[i]
}

// 格式化时间
const formatTime = (time: string): string => {
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
.share-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.share-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  max-width: 500px;
  width: 100%;
}

.share-info {
  text-align: center;
  margin-bottom: 30px;
}

.file-icon {
  margin-bottom: 20px;
}

.file-details {
  text-align: left;
}

.file-name {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0 0 10px 0;
  word-break: break-all;
}

.file-meta {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  font-size: 14px;
  color: #666;
}

.file-description {
  padding: 12px;
  background: #f5f5f5;
  border-radius: 6px;
  margin-bottom: 20px;
  font-size: 14px;
  color: #666;
}

.share-meta {
  margin-top: 20px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.label {
  color: #999;
}

.value {
  color: #333;
  font-weight: 500;
}

.limit {
  color: #999;
}

.password-section,
.download-section,
.error-section,
.loading-section {
  margin-top: 30px;
}

.loading-section {
  text-align: center;
}

.loading-section p {
  margin-top: 15px;
  color: #666;
}

@media (max-width: 768px) {
  .share-container {
    padding: 20px;
    margin: 10px;
  }
  
  .file-name {
    font-size: 20px;
  }
}
</style> 