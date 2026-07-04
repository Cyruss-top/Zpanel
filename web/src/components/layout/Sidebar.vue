<script setup lang="ts">
import { h, computed, type Component } from 'vue'
import { NIcon, type MenuOption } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import {
  SpeedometerOutline,
  GlobeOutline,
  ServerOutline,
  LibraryOutline,
  FolderOutline,
  TimeOutline,
  DocumentTextOutline,
  SettingsOutline,
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

const activeKey = computed(() => {
  const name = route.name as string
  if (name === 'dashboard') return 'dashboard'
  return name || 'dashboard'
})

function icon(comp: Component) {
  return () => h(NIcon, null, { default: () => h(comp) })
}

const menuOptions: MenuOption[] = [
  { label: '概览', key: 'dashboard', icon: icon(SpeedometerOutline) },
  { label: '网站', key: 'sites', icon: icon(GlobeOutline) },
  { label: '环境', key: 'environment', icon: icon(ServerOutline) },
  { label: '数据库', key: 'database', icon: icon(LibraryOutline) },
  { label: '文件', key: 'files', icon: icon(FolderOutline) },
  { label: '计划任务', key: 'cron', icon: icon(TimeOutline) },
  { label: '日志', key: 'logs', icon: icon(DocumentTextOutline) },
  { label: '设置', key: 'settings', icon: icon(SettingsOutline) },
]

function onSelect(key: string) {
  if (key === 'dashboard') {
    router.push('/')
  } else {
    router.push(`/${key}`)
  }
}
</script>

<template>
  <div class="sidebar-brand">Zpanel</div>
  <n-menu
    :value="activeKey"
    :options="menuOptions"
    :collapsed-width="64"
    :collapsed-icon-size="18"
    @update:value="onSelect"
  />
</template>

<style scoped>
.sidebar-brand {
  padding: 20px 20px 12px;
  font-size: 16px;
  font-weight: 600;
  color: #fafafa;
  letter-spacing: 0.02em;
}
</style>
