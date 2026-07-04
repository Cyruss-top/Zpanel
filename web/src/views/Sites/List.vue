<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NButton, NTag } from 'naive-ui'
import { fetchSites, deleteSite, type Site } from '@/api/site'
import { APIError } from '@/api/client'

const router = useRouter()
const message = useMessage()
const loading = ref(true)
const sites = ref<Site[]>([])

const typeColor: Record<string, 'default' | 'info' | 'success'> = {
  html: 'default',
  php: 'info',
  go: 'success',
}

const columns = [
  { title: '名称', key: 'name' },
  {
    title: '类型',
    key: 'type',
    render: (row: Site) =>
      h(NTag, { bordered: true, size: 'small', type: typeColor[row.type] || 'default' }, () =>
        row.type.toUpperCase()
      ),
  },
  {
    title: '域名',
    key: 'domains',
    render: (row: Site) => row.domains.join(', '),
  },
  {
    title: '状态',
    key: 'status',
    render: (row: Site) =>
      h(
        NTag,
        {
          bordered: true,
          size: 'small',
          type: row.status === 'running' ? 'success' : 'warning',
        },
        () => row.status
      ),
  },
  {
    title: '操作',
    key: 'actions',
    render: (row: Site) =>
      h(
        NButton,
        { size: 'small', type: 'error', quaternary: true, onClick: () => onDelete(row) },
        () => '删除'
      ),
  },
]

async function load() {
  loading.value = true
  try {
    sites.value = await fetchSites()
  } catch (e) {
    message.error(e instanceof APIError ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

async function onDelete(site: Site) {
  try {
    await deleteSite(site.id)
    message.success('已删除')
    await load()
  } catch (e) {
    message.error(e instanceof APIError ? e.message : '删除失败')
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-head">
      <h1 class="page-title">网站管理</h1>
      <n-button type="primary" @click="router.push('/sites/create')">添加站点</n-button>
    </div>
    <n-spin :show="loading">
      <n-card :bordered="true">
        <n-data-table :columns="columns" :data="sites" :bordered="false" size="small" />
        <p v-if="!loading && sites.length === 0" class="empty">暂无站点</p>
      </n-card>
    </n-spin>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.page-head .page-title {
  margin: 0;
}
.empty {
  margin: 16px 0 0;
  color: #71717a;
  font-size: 14px;
}
</style>
