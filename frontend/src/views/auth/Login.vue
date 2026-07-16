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
      <el-button type="primary" size="large" style="width: 100%" @click="redirectToSSO('', '/')">{{ buttonText }}</el-button>
      <div class="login-footer">
        <p>统一 LDAP 账号，同账号同密码</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { redirectToSSO } from '@/utils/sso'

const route = useRoute()
const loggedOut = computed(() => route.query.logged_out === '1')
const title = computed(() => loggedOut.value ? '已退出登录' : '统一基础设施平台')
const subtitle = computed(() => loggedOut.value ? '本地登录态已清除' : '正在跳转到 SSO 登录')
const buttonText = computed(() => loggedOut.value ? '重新登录' : '重新跳转')

onMounted(() => {
  if (!loggedOut.value) {
    redirectToSSO('', '/')
  }
})
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
  border-radius: 12px;
  padding: 40px;
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
  border-radius: 8px;
  background: linear-gradient(135deg, var(--accent), #1C8F69);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 16px;
  font-weight: 600;
  color: #08140F;
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
  border-color: #3A4356;
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
  color: #06150F;
  font-weight: 600;
}

:deep(.el-button--primary:hover) {
  background: #55E8B3;
}
</style>
