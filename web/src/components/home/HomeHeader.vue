<template>
  <div>
    <n-flex class="search-bar" justify="space-between">
      <div>
        <n-breadcrumb>
          <n-breadcrumb-item v-for="(item, index) in breadcrumb"
                             :key="index"
                             @click="handelBreadcrumbClick(item)">
            <n-icon :component="FolderOpenTwotone" />
            {{ item}}
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

const breadcrumb = defineModel('breadcrumb')

const emit = defineEmits(['breadcrumb-click'])

const handelBreadcrumbClick = (item: string) => {
  emit('breadcrumb-click', item)
}
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
