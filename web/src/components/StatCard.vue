<script setup lang="ts">
defineProps<{
  title: string
  value: string | number
  sub?: string
  warn?: boolean
}>()
</script>

<template>
  <n-card :bordered="true" size="small">
    <div class="stat-label">{{ title }}</div>
    <div class="stat-value" :class="{ warn }">{{ value }}</div>
    <div v-if="sub" class="stat-sub mono">{{ sub }}</div>
  </n-card>
</template>

<style scoped>
.stat-label {
  font-size: 12px;
  color: #71717a;
  margin-bottom: 8px;
}
.stat-value {
  font-size: 28px;
  font-weight: 600;
  line-height: 1.2;
  color: #18181b;
}
.stat-value.warn {
  color: #d97706;
}
.stat-sub {
  margin-top: 8px;
  font-size: 12px;
  color: #71717a;
}
</style>
