<template>
  <div class="file-manager" style="height: calc(100% - 45px)">
    <n-card class="home-max-card" content-style="height: calc(100% - 120px)">
      <template v-slot:header>
        <n-icon>
          <BellTwotone />
        </n-icon>
        File Manager
      </template>
      <template v-slot:default>
        <HomeHeader />
        <div class="file-list" style="height: calc(100% - 55px)">
          <div class="file-item-header">
            <span class="font-title">名称</span>
            <span class="font-title">大小</span>
            <span class="font-title">修改时间</span>
          </div>
          <n-scrollbar>
            <div v-for="(file, index) in files" :key="index" class="file-item" :class="{ selected: selectedFileIndex === index }" @click="toggleSelection(index)" @mouseover="handleMouseOver(index)" @mouseout="handleMouseOut(index)">
              <span>
                <n-icon>
                  <FileTextOutlined />
                </n-icon>
                {{ file.name }}
              </span>
              <span>{{ file.size }} KB</span>
              <span>{{ file.date }}</span>
            </div>
          </n-scrollbar>
        </div>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { FileTextOutlined, BellTwotone } from '@vicons/antd';
import HomeHeader from '@/components/home/HomeHeader.vue';

// 示例数据
const files = [
  { name: 'document.txt', size: 1024, type: 'file', date: '2023-05-18' },
  { name: 'image.jpg', size: 512, type: 'file', date: '2023-05-18' },
  { name: 'video.mp4', size: 20480, type: 'file', date: '2023-05-18' },
  { name: 'report.pdf', size: 1536, type: 'file', date: '2023-05-18' },
  { name: 'notes.txt', size: 768, type: 'file', date: '2023-05-18' },
  { name: 'Documents', size: 0, type: 'folder', date: '2023-05-18' },
  { name: 'document.txt', size: 1024, type: 'file', date: '2023-05-18' }
];

// 选中的文件索引
const selectedFileIndex = ref(-1);

// 鼠标悬停的文件索引
const hoveredFileIndex = ref(-1);

// 切换文件选中状态
function toggleSelection(index: number) {
  selectedFileIndex.value = index === selectedFileIndex.value ? -1 : index;
}

// 处理鼠标悬停事件
function handleMouseOver(index: number) {
  hoveredFileIndex.value = index;
}

// 处理鼠标移出事件
function handleMouseOut() {
  hoveredFileIndex.value = -1;
}
</script>

<style scoped>
.file-manager {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.font-title {
  font-weight: bold;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.3s ease;
}

.file-item.selected {
  background-color: #40a9ff;
  color: white;
}

.file-item:hover {
  background-color: #e6f7ff;
  transform: scale(1.02);
}

.file-item:hover:not(.selected) {
  background-color: #e6f7ff;
  transform: scale(1.02);
}

.n-scrollbar .n-scrollbar__wrap {
  height: calc(100% - 55px);
}
</style>
