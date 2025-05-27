<template>
  <n-icon :size="size" :class="iconClass">
    <component :is="iconComponent" />
  </n-icon>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NIcon } from 'naive-ui'
import {
  Folder as FolderIcon,
  FolderOpen as FolderOpenIcon,
  Image as ImageIcon,
  VideocamSharp as VideoIcon,
  MusicalNotes as AudioIcon,
  DocumentText as DocumentIcon,
  Code as CodeIcon,
  Archive as ArchiveIcon,
  Document as FileIcon
} from '@vicons/ionicons5'
import type { FileInfo } from '@/types/file'
import { getFileIcon } from '@/utils/file'

interface Props {
  file: FileInfo
  size?: string | number
  open?: boolean // 文件夹是否打开状态
}

const props = withDefaults(defineProps<Props>(), {
  size: 24,
  open: false
})

// 根据文件类型获取图标组件
const iconComponent = computed(() => {
  if (props.file.isDir) {
    return props.open ? FolderOpenIcon : FolderIcon
  }

  const iconType = getFileIcon(props.file)

  switch (iconType) {
    case 'image':
      return ImageIcon
    case 'video':
      return VideoIcon
    case 'audio':
      return AudioIcon
    case 'pdf':
    case 'word':
    case 'excel':
    case 'powerpoint':
    case 'text':
      return DocumentIcon
    case 'code':
      return CodeIcon
    case 'archive':
      return ArchiveIcon
    default:
      return FileIcon
  }
})

// 图标样式类
const iconClass = computed(() => {
  if (props.file.isDir) {
    return 'file-icon folder-icon'
  }

  const iconType = getFileIcon(props.file)
  return `file-icon ${iconType}-icon`
})
</script>

<style scoped>
.file-icon {
  color: #666;
  transition: color 0.2s;
}

.folder-icon {
  color: #faad14;
}

.image-icon {
  color: #52c41a;
}

.video-icon {
  color: #722ed1;
}

.audio-icon {
  color: #eb2f96;
}

.pdf-icon,
.word-icon,
.excel-icon,
.powerpoint-icon,
.text-icon {
  color: #1890ff;
}

.code-icon {
  color: #13c2c2;
}

.archive-icon {
  color: #fa8c16;
}

.file-icon:hover {
  opacity: 0.8;
}
</style>
