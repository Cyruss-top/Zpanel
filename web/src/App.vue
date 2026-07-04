<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import { RouterView } from 'vue-router'
import { lightThemeOverrides } from './styles/theme'
</script>

<template>
  <n-config-provider :theme-overrides="lightThemeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <RouterView />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>
