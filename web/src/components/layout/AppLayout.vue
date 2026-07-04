<script setup lang="ts">
import { ref, computed } from 'vue'
import { useBreakpoints } from '@vueuse/core'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import MobileTabBar from './MobileTabBar.vue'

const bp = useBreakpoints({ mobile: 768, tablet: 1200 })
const isMobile = bp.smaller('mobile')
const isTablet = bp.between('mobile', 'tablet')
const siderCollapsed = computed(() => isTablet.value)
const drawerOpen = ref(false)
</script>

<template>
  <n-layout has-sider style="min-height: 100vh">
    <n-layout-sider
      v-if="!isMobile"
      bordered
      collapse-mode="width"
      :collapsed="siderCollapsed"
      :collapsed-width="64"
      :width="220"
      :native-scrollbar="false"
    >
      <Sidebar />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered>
        <Header :mobile="isMobile" @open-menu="drawerOpen = true" />
      </n-layout-header>

      <n-layout-content :content-style="{ padding: isMobile ? '12px' : '24px' }">
        <router-view />
      </n-layout-content>

      <n-layout-footer v-if="isMobile">
        <MobileTabBar @open-menu="drawerOpen = true" />
      </n-layout-footer>
    </n-layout>
  </n-layout>

  <n-drawer v-if="isMobile" v-model:show="drawerOpen" :width="260" placement="left">
    <n-drawer-content body-content-style="padding: 0; background: #18181b">
      <Sidebar />
    </n-drawer-content>
  </n-drawer>
</template>
