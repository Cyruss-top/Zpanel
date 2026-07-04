<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const tabs = [
  { key: 'dashboard', label: '概览', path: '/' },
  { key: 'sites', label: '网站', path: '/sites' },
  { key: 'environment', label: '环境', path: '/environment' },
  { key: 'more', label: '更多', path: '' },
]

const emit = defineEmits<{ (e: 'open-menu'): void }>()

function onTab(key: string, path: string) {
  if (key === 'more') {
    emit('open-menu')
    return
  }
  router.push(path)
}

function isActive(key: string) {
  if (key === 'dashboard') return route.name === 'dashboard'
  return route.name === key
}
</script>

<template>
  <nav class="tab-bar">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      class="tab-item"
      :class="{ active: isActive(tab.key) }"
      @click="onTab(tab.key, tab.path)"
    >
      {{ tab.label }}
    </button>
  </nav>
</template>

<style scoped>
.tab-bar {
  display: flex;
  border-top: 1px solid #e4e4e7;
  background: #fff;
}
.tab-item {
  flex: 1;
  border: none;
  background: transparent;
  padding: 10px 4px;
  font-size: 12px;
  color: #71717a;
  cursor: pointer;
}
.tab-item.active {
  color: #2563eb;
  font-weight: 600;
}
</style>
