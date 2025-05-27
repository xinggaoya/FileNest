<script setup lang="ts">
import { RouterView } from 'vue-router'
import { zhCN, darkTheme, useOsTheme } from 'naive-ui'
import themeOverrides from '@/assets/theme/naive-ui-theme-overrides.json'
import { computed } from 'vue'

import { NMessageProvider, NGlobalStyle, NConfigProvider } from 'naive-ui'

const osTheme = useOsTheme()
const theme = computed(() => (osTheme.value === 'dark' ? darkTheme : null))
const locale = zhCN
</script>

<template>
  <div id="app">
    <n-config-provider :theme="theme" :locale="locale" :theme-overrides="themeOverrides">
      <n-global-style />
      <n-message-provider>
        <RouterView />
      </n-message-provider>
    </n-config-provider>
  </div>
</template>

<style>
@import './styles/global.css';

#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>
