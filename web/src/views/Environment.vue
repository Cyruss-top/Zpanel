<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { fetchLNMPStatus, installLNMP, type LNMPStatus } from '@/api/lnmp'
import { APIError } from '@/api/client'

const message = useMessage()
const loading = ref(true)
const installing = ref(false)
const status = ref<LNMPStatus | null>(null)

const rows = [
  { key: 'nginx', label: 'Nginx' },
  { key: 'mysql', label: 'MySQL' },
  { key: 'php', label: 'PHP-FPM' },
]

async function load() {
  loading.value = true
  try {
    status.value = await fetchLNMPStatus()
  } catch (e) {
    message.error(e instanceof APIError ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

async function onInstall() {
  installing.value = true
  try {
    await installLNMP()
    message.success('LNMP 安装完成')
    await load()
  } catch (e) {
    message.error(e instanceof APIError ? e.message : '安装失败')
  } finally {
    installing.value = false
  }
}

function statusTag(comp?: { installed: boolean; running: boolean }) {
  if (!comp?.installed) return { type: 'default' as const, text: '未安装' }
  if (comp.running) return { type: 'success' as const, text: '运行中' }
  return { type: 'warning' as const, text: '已停止' }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-head">
      <h1 class="page-title">环境管理</h1>
      <n-button type="primary" :loading="installing" @click="onInstall">一键安装 LNMP</n-button>
    </div>

    <n-spin :show="loading">
      <n-card v-if="status" :bordered="true" title="LNMP 状态">
        <p class="platform mono">平台: {{ status.platform }}</p>
        <n-table :bordered="false" size="small">
          <thead>
            <tr>
              <th>组件</th>
              <th>版本</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.key">
              <td>{{ row.label }}</td>
              <td class="mono">{{ status.components[row.key]?.version || '-' }}</td>
              <td>
                <n-tag
                  size="small"
                  :bordered="true"
                  :type="statusTag(status.components[row.key]).type"
                >
                  {{ statusTag(status.components[row.key]).text }}
                </n-tag>
              </td>
            </tr>
          </tbody>
        </n-table>
        <p v-if="status.platform !== 'linux'" class="hint">
          LNMP 安装与检测需要在 Linux 服务器上运行。当前环境仅可查看状态。
        </p>
      </n-card>
    </n-spin>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}
.page-head .page-title {
  margin: 0;
}
.platform {
  margin: 0 0 12px;
  font-size: 13px;
  color: #71717a;
}
.hint {
  margin: 16px 0 0;
  font-size: 13px;
  color: #71717a;
}
</style>
