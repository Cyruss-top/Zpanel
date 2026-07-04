<script setup lang="ts">
defineProps<{ title: string }>()
</script>

<template>
  <div>
    <h1 class="page-title">{{ title }}</h1>
    <n-card :bordered="true">
      <p class="hint">功能开发中，将在后续版本开放。</p>
    </n-card>
  </div>
</template>

<style scoped>
.hint {
  margin: 0;
  color: #71717a;
  font-size: 14px;
}
</style>
