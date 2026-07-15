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
import { consumeOpenApp } from '@/utils/sso'

const subsystems = ref<Subsystem[]>([])

onMounted(async () => {
  try {
    const { data } = await subsystemApi.getSubsystems()
    subsystems.value = data
    const openApp = consumeOpenApp()
    if (openApp) {
      const targetName = openApp === 'clouddm' ? 'CloudDM' : 'Wayne'
      const target = data.find((item) => item.name === targetName)
      if (target) {
        const response = await subsystemApi.getSSOUrl(target.id)
        window.location.assign(response.data.sso_url)
      }
    }
  } catch {
    // 错误已由 request 拦截器统一处理
  }
})
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 18px;
}

.page-head h1 {
  font-size: 19px;
  margin: 0 0 4px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.page-head p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12.5px;
}

.ext-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
}
</style>
