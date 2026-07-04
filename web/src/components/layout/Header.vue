<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

defineProps<{ mobile?: boolean }>()
const emit = defineEmits<{ (e: 'open-menu'): void }>()

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="header-inner">
    <div class="left">
      <n-button v-if="mobile" quaternary @click="emit('open-menu')">菜单</n-button>
      <span v-else class="header-title">控制台</span>
    </div>
    <div class="right">
      <span class="username">{{ auth.username }}</span>
      <n-button size="small" @click="logout">退出</n-button>
    </div>
  </div>
</template>

<style scoped>
.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 16px;
}
.header-title {
  font-size: 14px;
  font-weight: 600;
  color: #18181b;
}
.username {
  margin-right: 12px;
  font-size: 13px;
  color: #71717a;
}
</style>
