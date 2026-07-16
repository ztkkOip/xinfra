<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="brand">
          <div class="mark">xi</div>
          xinfra
        </div>
        <h2>{{ title }}</h2>
        <p>{{ subtitle }}</p>
      </div>
      <div v-if="!ssoEnabled" class="local-login">
        <el-input
          v-model="username"
          size="large"
          placeholder="输入用户名"
          @keyup.enter="handleLocalLogin"
        />
      </div>
      <el-button type="primary" size="large" class="login-btn" :loading="loading" @click="handleLoginClick">
        {{ buttonText }}
      </el-button>
      <div class="login-footer">
        <p>{{ footerText }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { redirectToSSO } from '@/utils/sso'
import { useAuth } from '@/composables/useAuth'
import { authApi } from '@/api/auth'

const route = useRoute()
const router = useRouter()
const { login } = useAuth()
const loading = ref(false)
const username = ref('')
const ssoEnabled = ref(true)
const loggedOut = computed(() => route.query.logged_out === '1')
const title = computed(() => loggedOut.value ? '已退出登录' : '统一基础设施平台')
const subtitle = computed(() => {
  if (!ssoEnabled.value) {
    return '开发测试模式，输入用户名登录'
  }
  return loggedOut.value ? '本地登录态已清除' : '正在跳转到 SSO 登录'
})
const buttonText = computed(() => {
  if (!ssoEnabled.value) {
    return '登录'
  }
  return loggedOut.value ? '重新登录' : '重新跳转'
})
const footerText = computed(() => ssoEnabled.value ? '统一 LDAP 账号，同账号同密码' : 'SSO 已关闭，仅用于开发测试')

onMounted(async () => {
  const { data } = await authApi.getConfig()
  ssoEnabled.value = data.sso_enabled
  if (ssoEnabled.value && !loggedOut.value) {
    redirectToSSO('', '/')
  }
})

const handleLoginClick = () => {
  if (ssoEnabled.value) {
    redirectToSSO('', '/')
    return
  }
  handleLocalLogin()
}

const handleLocalLogin = async () => {
  const value = username.value.trim()
  if (!value) {
    return
  }
  loading.value = true
  try {
    await login(value, '')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-deep);
}

.login-card {
  width: 400px;
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: var(--radius-xl);
  padding: 40px;
}

.login-btn {
  width: 100%;
}

.local-login {
  margin-bottom: 14px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header .brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 12px;
}

.login-header .brand .mark {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, var(--accent), var(--accent-end));
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 16px;
  font-weight: 600;
  color: var(--bg-deep);
}

.login-header h2 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
}

.login-header p {
  margin: 0;
  color: var(--text-dim);
  font-size: 14px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--line-soft);
}

.login-footer p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12px;
}

:deep(.el-form-item__label) {
  color: var(--text-mid);
}

:deep(.el-input__wrapper) {
  background: var(--bg-panel-2);
  border: 1px solid var(--line);
  box-shadow: none;
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--hover-border);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--accent);
}

:deep(.el-input__inner) {
  color: var(--text-hi);
}

:deep(.el-button--primary) {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--bg-deep);
  font-weight: 600;
}

:deep(.el-button--primary:hover) {
  opacity: 0.9;
}
</style>
