<template>
  <div>
    <n-flex class="search-bar" justify="space-between">
      <div>
        <n-breadcrumb>
          <n-breadcrumb-item>
            <n-icon :component="FolderOpenTwotone" />
            Home
          </n-breadcrumb-item>
          <n-breadcrumb-item>
            <n-icon :component="FolderOpenTwotone" />
            pathA
          </n-breadcrumb-item>
          <n-breadcrumb-item>
            <n-icon :component="FolderOpenTwotone" />
            pathB
          </n-breadcrumb-item>
        </n-breadcrumb>
      </div>
      <div style="width: 160px">
        <n-select v-model="sortType" :options="sortOptions" />
      </div>
    </n-flex>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { FolderOpenTwotone } from '@vicons/antd'

// 排序选项
const sortOptions = [
  { label: 'Name', value: 'name' },
  { label: 'Size', value: 'size' }
]
const sortType = ref('name')
</script>

<style scoped>
.search-bar {
  margin-bottom: 10px;
}
</style>
