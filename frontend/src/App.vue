<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { authApi } from '@/api/auth'
import { getToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  const token = getToken()
  if (!token) {
    return
  }

  try {
    await authStore.refreshUser()
  } catch {
    authStore.clearAuth()
    const { data } = await authApi.getConfig()
    if (data.sso_enabled) {
      redirectToSSO()
    } else {
      window.location.assign('/login')
    }
  }
})
</script>

<style>
/* 全局样式已在 global.css 中定义 */
</style>
