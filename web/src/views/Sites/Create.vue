<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { createSite, type SiteType } from '@/api/site'
import { APIError } from '@/api/client'

const router = useRouter()
const message = useMessage()

const form = ref({
  name: '',
  type: 'html' as SiteType,
  domains: '',
  php_version: '8.2',
  go_port: 8080,
  go_binary: '',
})
const loading = ref(false)

const typeOptions = [
  { label: 'HTML 静态站', value: 'html' },
  { label: 'PHP 站点', value: 'php' },
  { label: 'Go 反向代理', value: 'go' },
]

async function onSubmit() {
  const domains = form.value.domains
    .split(/[,\s]+/)
    .map((d) => d.trim())
    .filter(Boolean)
  if (!form.value.name || domains.length === 0) {
    message.warning('请填写站点名称和域名')
    return
  }
  loading.value = true
  try {
    await createSite({
      name: form.value.name,
      type: form.value.type,
      domains,
      php_version: form.value.type === 'php' ? form.value.php_version : undefined,
      go_port: form.value.type === 'go' ? form.value.go_port : undefined,
      go_binary: form.value.type === 'go' ? form.value.go_binary : undefined,
    })
    message.success('站点创建成功')
    router.push('/sites')
  } catch (e) {
    message.error(e instanceof APIError ? e.message : '创建失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="page-title">添加站点</h1>
    <n-card :bordered="true" style="max-width: 560px">
      <n-form label-placement="top">
        <n-form-item label="站点名称（目录名）">
          <n-input v-model:value="form.name" placeholder="example.com" />
        </n-form-item>
        <n-form-item label="站点类型">
          <n-select v-model:value="form.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="域名（多个用逗号或空格分隔）">
          <n-input v-model:value="form.domains" placeholder="example.com www.example.com" />
        </n-form-item>
        <n-form-item v-if="form.type === 'php'" label="PHP 版本">
          <n-input v-model:value="form.php_version" placeholder="8.2" />
        </n-form-item>
        <template v-if="form.type === 'go'">
          <n-form-item label="监听端口">
            <n-input-number v-model:value="form.go_port" :min="1024" :max="65535" style="width: 100%" />
          </n-form-item>
          <n-form-item label="二进制路径（可选）">
            <n-input v-model:value="form.go_binary" placeholder="/var/www/example/app" />
          </n-form-item>
        </template>
        <n-space>
          <n-button type="primary" :loading="loading" @click="onSubmit">创建</n-button>
          <n-button @click="router.push('/sites')">取消</n-button>
        </n-space>
      </n-form>
    </n-card>
  </div>
</template>
