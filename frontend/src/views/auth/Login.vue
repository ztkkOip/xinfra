<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="brand">
          <div class="mark">xi</div>
          xinfra
        </div>
        <h2>统一基础设施平台</h2>
        <p>正在跳转到 SSO 登录</p>
      </div>
      <el-button type="primary" size="large" style="width: 100%" @click="redirectToSSO()">
        重新跳转
      </el-button>
      <div class="login-footer">
        <p>统一 LDAP 账号，同账号同密码</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { redirectToSSO } from '@/utils/sso'

onMounted(() => redirectToSSO())
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
