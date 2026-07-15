<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { authApi } from '@/api/auth'
import { getToken, removeToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  const token = getToken()
  if (!token) {
    return
  }

  try {
    const { data } = await authApi.getUserInfo()
    authStore.setAuth(token, data)
  } catch {
    authStore.clearAuth()
    removeToken()
    redirectToSSO()
  }
})
</script>

<style>
/* 全局样式已在 global.css 中定义 */
</style>
