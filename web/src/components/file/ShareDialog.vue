<template>
  <n-modal v-model:show="visible" preset="card" title="创建分享" style="width: 500px">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="left" label-width="100px">
      <n-form-item label="文件路径" path="filePath">
        <n-input v-model:value="formData.filePath" readonly />
      </n-form-item>
      
      <n-form-item label="分享密码" path="password">
        <n-input
          v-model:value="formData.password"
          type="password"
          placeholder="留空则无需密码（可选）"
          show-password-on="click"
          clearable
        />
      </n-form-item>
      
      <n-form-item label="有效期" path="expireHours">
        <n-select
          v-model:value="formData.expireHours"
          :options="expireOptions"
          placeholder="选择有效期"
          clearable
        />
      </n-form-item>
      
      <n-form-item label="下载次数" path="maxDownload">
        <n-input-number
          v-model:value="formData.maxDownload"
          :min="0"
          :max="1000"
          placeholder="0表示不限制"
          style="width: 100%"
        />
      </n-form-item>
      
      <n-form-item label="分享描述" path="description">
        <n-input
          v-model:value="formData.description"
          type="textarea"
          :rows="3"
          placeholder="输入分享描述（可选）"
          :maxlength="200"
          show-count
        />
      </n-form-item>
    </n-form>
    
    <template #footer>
      <div class="dialog-footer">
        <n-button @click="handleCancel">取消</n-button>
        <n-button type="primary" :loading="creating" @click="handleCreate">创建分享</n-button>
      </div>
    </template>
  </n-modal>
  
  <!-- 分享结果对话框 -->
  <n-modal v-model:show="resultVisible" preset="card" title="分享创建成功" style="width: 500px">
    <div class="share-result">
      <n-alert type="success" title="分享链接已生成" show-icon>
        <div class="share-info">
          <div class="info-item">
            <span class="label">分享码：</span>
            <span class="value">{{ shareResult?.shareCode }}</span>
            <n-button text @click="copyShareCode">复制</n-button>
          </div>
          
          <div class="info-item">
            <span class="label">分享链接：</span>
            <span class="value break-all">{{ shareUrl }}</span>
            <n-button text @click="copyShareUrl">复制</n-button>
          </div>
          
          <div v-if="shareResult?.password" class="info-item">
            <span class="label">提取密码：</span>
            <span class="value">{{ shareResult.password }}</span>
            <n-button text @click="copyPassword">复制</n-button>
          </div>
          
          <div v-if="shareResult?.expireTime" class="info-item">
            <span class="label">过期时间：</span>
            <span class="value">{{ formatTime(shareResult.expireTime) }}</span>
          </div>
          
          <div v-if="shareResult?.maxDownload" class="info-item">
            <span class="label">下载限制：</span>
            <span class="value">{{ shareResult.maxDownload }}次</span>
          </div>
        </div>
      </n-alert>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <n-button @click="handleClose">关闭</n-button>
        <n-button type="primary" @click="copyAllInfo">复制分享信息</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NButton,
  NAlert,
  useMessage
} from 'naive-ui'
import { createShare } from '@/api/share/share'
import type { CreateShareRequest, ShareItem } from '@/api/share/share'

// Props
interface Props {
  modelValue: boolean
  filePath: string
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'success': [share: ShareItem]
}>()

// 状态
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formRef = ref()
const creating = ref(false)
const resultVisible = ref(false)
const shareResult = ref<ShareItem | null>(null)
const message = useMessage()

// 表单数据
const formData = ref<CreateShareRequest>({
  filePath: '',
  password: '',
  expireHours: undefined,
  maxDownload: 0,
  description: ''
})

// 表单规则
const formRules = {
  filePath: [
    { required: true, message: '文件路径不能为空', trigger: 'blur' }
  ]
}

// 过期时间选项
const expireOptions = [
  { label: '1小时', value: 1 },
  { label: '6小时', value: 6 },
  { label: '12小时', value: 12 },
  { label: '1天', value: 24 },
  { label: '3天', value: 72 },
  { label: '7天', value: 168 },
  { label: '30天', value: 720 },
  { label: '永久有效', value: 0 }
]

// 分享链接
const shareUrl = computed(() => {
  if (!shareResult.value) return ''
  const baseUrl = window.location.origin
  return `${baseUrl}/#/share/${shareResult.value.shareCode}`
})

// 监听文件路径变化
const initForm = () => {
  formData.value.filePath = props.filePath
}

// 创建分享
const handleCreate = async () => {
  try {
    await formRef.value?.validate()
    creating.value = true
    
    const { data } = await createShare(formData.value)
    shareResult.value = data
    
    emit('success', data)
    visible.value = false
    resultVisible.value = true
    
    message.success('分享创建成功')
  } catch (error: any) {
    if (error.message) {
      message.error(error.message)
    }
  } finally {
    creating.value = false
  }
}

// 取消
const handleCancel = () => {
  visible.value = false
  resetForm()
}

// 关闭结果对话框
const handleClose = () => {
  resultVisible.value = false
  resetForm()
}

// 重置表单
const resetForm = () => {
  formData.value = {
    filePath: '',
    password: '',
    expireHours: undefined,
    maxDownload: 0,
    description: ''
  }
  shareResult.value = null
}

// 复制分享码
const copyShareCode = async () => {
  if (shareResult.value?.shareCode) {
    await navigator.clipboard.writeText(shareResult.value.shareCode)
    message.success('分享码已复制')
  }
}

// 复制分享链接
const copyShareUrl = async () => {
  if (shareUrl.value) {
    await navigator.clipboard.writeText(shareUrl.value)
    message.success('分享链接已复制')
  }
}

// 复制密码
const copyPassword = async () => {
  if (shareResult.value?.password) {
    await navigator.clipboard.writeText(shareResult.value.password)
    message.success('提取密码已复制')
  }
}

// 复制所有分享信息
const copyAllInfo = async () => {
  if (!shareResult.value) return
  
  let info = `文件分享\n`
  info += `文件名：${shareResult.value.fileName}\n`
  info += `分享链接：${shareUrl.value}\n`
  
  if (shareResult.value.password) {
    info += `提取密码：${shareResult.value.password}\n`
  }
  
  if (shareResult.value.expireTime) {
    info += `过期时间：${formatTime(shareResult.value.expireTime)}\n`
  }
  
  if (shareResult.value.description) {
    info += `描述：${shareResult.value.description}\n`
  }
  
  await navigator.clipboard.writeText(info)
  message.success('分享信息已复制')
}

// 格式化时间
const formatTime = (time: string): string => {
  return new Date(time).toLocaleString('zh-CN')
}

// 监听props变化
watch(() => props.filePath, () => {
  if (props.filePath) {
    initForm()
  }
}, { immediate: true })

// 监听visible变化
watch(visible, (newVal) => {
  if (newVal && props.filePath) {
    initForm()
  }
})
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.share-result {
  margin: 20px 0;
}

.share-info {
  margin-top: 15px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  font-size: 14px;
}

.label {
  min-width: 80px;
  color: #666;
  font-weight: 500;
}

.value {
  flex: 1;
  color: #333;
  margin-right: 10px;
}

.break-all {
  word-break: break-all;
}
</style> 