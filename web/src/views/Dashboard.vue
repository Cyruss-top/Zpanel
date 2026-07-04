<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import StatCard from '@/components/StatCard.vue'
import { fetchOverview, type MonitorOverview } from '@/api/monitor'
import { APIError } from '@/api/client'

const message = useMessage()
const loading = ref(true)
const data = ref<MonitorOverview | null>(null)
let timer: number | undefined

function formatBytes(n: number) {
  if (!n) return '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let v = n
  let i = 0
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(1)} ${units[i]}`
}

function formatUptime(sec: number) {
  const d = Math.floor(sec / 86400)
  const h = Math.floor((sec % 86400) / 3600)
  const m = Math.floor((sec % 3600) / 60)
  if (d > 0) return `${d}天 ${h}小时`
  if (h > 0) return `${h}小时 ${m}分钟`
  return `${m}分钟`
}

async function load() {
  try {
    data.value = await fetchOverview()
  } catch (e) {
    if (e instanceof APIError && e.code !== 'UNAUTHORIZED') {
      message.error(e.message)
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
  timer = window.setInterval(load, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div>
    <h1 class="page-title">系统概览</h1>

    <n-spin :show="loading && !data">
      <div v-if="data" class="stat-grid">
        <StatCard
          title="CPU 使用率"
          :value="`${data.cpu.usage_percent.toFixed(1)}%`"
          :sub="`${data.cpu.cores} 核心`"
          :warn="data.cpu.usage_percent > 80"
        />
        <StatCard
          title="内存使用率"
          :value="`${data.memory.used_percent.toFixed(1)}%`"
          :sub="`${formatBytes(data.memory.used)} / ${formatBytes(data.memory.total)}`"
          :warn="data.memory.used_percent > 80"
        />
        <StatCard
          title="磁盘使用率"
          :value="data.disk.total ? `${data.disk.used_percent.toFixed(1)}%` : '-'"
          :sub="data.disk.total ? formatBytes(data.disk.used) + ' / ' + formatBytes(data.disk.total) : data.disk.path"
          :warn="data.disk.used_percent > 80"
        />
        <StatCard
          title="运行时间"
          :value="formatUptime(data.uptime_seconds)"
          :sub="`${data.os} / ${data.arch}`"
        />
      </div>

      <n-card v-if="data" :bordered="true" size="small" title="负载">
        <div class="load-row mono">
          <span>1m: {{ data.load.load1.toFixed(2) }}</span>
          <span>5m: {{ data.load.load5.toFixed(2) }}</span>
          <span>15m: {{ data.load.load15.toFixed(2) }}</span>
        </div>
        <div class="platform mono">{{ data.platform }}</div>
      </n-card>
    </n-spin>
  </div>
</template>

<style scoped>
.load-row {
  display: flex;
  gap: 24px;
  font-size: 14px;
}
.platform {
  margin-top: 12px;
  font-size: 12px;
  color: #71717a;
}
</style>
