<script setup lang="ts">
import { RouterView } from 'vue-router'
import { zhCN } from 'naive-ui'
import themeOverrides from '@/assets/theme/naive-ui-theme-overrides.json'
</script>

<template>
  <n-config-provider class="height-100" :theme-overrides="themeOverrides" :locale="zhCN">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <n-modal-provider>
            <RouterView />
          </n-modal-provider>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
.height-100 {
  height: 100%;
}
</style>
