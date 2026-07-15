<template>
  <div>
    <div class="page-head">
      <div>
        <h1>子系统导航</h1>
        <p>统一 LDAP 账号，点击即用免二次登录（同账号同密码）</p>
      </div>
    </div>
    <div class="ext-grid">
      <SubsystemCard
        v-for="system in subsystems"
        :key="system.id"
        :system="system"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { subsystemApi, type Subsystem } from '@/api/subsystem'
import SubsystemCard from '@/components/SubsystemCard.vue'

const subsystems = ref<Subsystem[]>([])

onMounted(async () => {
  try {
    const { data } = await subsystemApi.getSubsystems()
    subsystems.value = data
  } catch {
    // 错误已由 request 拦截器统一处理
  }
})
</script>

<style scoped>
/* .page-head, .ext-grid 已在 global.css 中定义 */
</style>
