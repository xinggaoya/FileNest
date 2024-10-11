<script setup lang="ts">
import { RouterView } from 'vue-router'
import { zhCN } from 'naive-ui'
import themeOverrides from '@/assets/theme/naive-ui-theme-overrides.json'
</script>

<template>
  <n-config-provider class="height-100" :theme-overrides="themeOverrides" :locale="zhCN">
    <RouterView />
  </n-config-provider>
</template>

<style scoped>

.height-100 {
  height: 100%;
}
</style>
