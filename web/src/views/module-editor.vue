<template>
  <a-page-header
    :style="{ background: 'var(--color-bg-2)' }"
    title="ArcoDesign"
    subtitle="ArcoDesign Vue 2.0"
    :show-back="false"
  >
    <template #extra></template>
  </a-page-header>
</template>

<script setup lang="ts">
  import { useRoute } from 'vue-router';

  const route = useRoute();
  const mode = route.name === 'new-module' ? 'new' : 'update';
</script>

<style scoped lang="less"></style>
