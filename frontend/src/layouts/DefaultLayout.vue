<template>
  <div class="shell">
    <AppHeader />
    <AppSidebar />
    <AppMain />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/Layout/AppHeader.vue'
import AppSidebar from '@/components/Layout/AppSidebar.vue'
import AppMain from '@/components/Layout/AppMain.vue'
</script>

<style scoped>
.shell {
  display: grid;
  grid-template-columns: 226px 1fr;
  grid-template-rows: 56px 1fr;
  height: 100vh;
}
</style>
