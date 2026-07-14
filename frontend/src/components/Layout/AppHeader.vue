<template>
  <header class="topbar">
    <div class="brand">
      <div class="mark">xi</div>
      xinfra <span class="sub">统一基础设施平台</span>
    </div>
    <span class="env-pill">● 生产环境</span>
    <div class="spacer"></div>
    <div class="search-box">
      <el-icon><Search /></el-icon>
      搜索集群 / 节点 / 服务实例
      <kbd>⌘K</kbd>
    </div>
    <div class="identity">
      <div class="avatar">{{ userInitial }}</div>
      <div>
        <div>{{ user?.display_name || '未登录' }}</div>
        <div class="ldap-tag">LDAP · {{ user?.business_line || '' }}业务线</div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const userInitial = computed(() => {
  return user.value?.display_name?.charAt(0) || 'U'
})
</script>

<style scoped>
.topbar {
  grid-column: 1/3;
  display: flex;
  align-items: center;
  padding: 0 18px 0 16px;
  background: var(--bg-panel);
  border-bottom: 1px solid var(--line);
  gap: 14px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 9px;
  font-weight: 600;
  letter-spacing: 0.3px;
  font-size: 15px;
}

.brand .mark {
  width: 22px;
  height: 22px;
  border-radius: 5px;
  background: linear-gradient(135deg, var(--accent), #1C8F69);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 600;
  color: #08140F;
}

.brand .sub {
  color: var(--text-dim);
  font-weight: 400;
  font-size: 12px;
  margin-left: 4px;
}

.spacer {
  flex: 1;
}

.env-pill {
  font-family: var(--mono);
  font-size: 11px;
  color: var(--accent);
  border: 1px solid var(--accent-dim);
  background: #0F261D;
  padding: 3px 9px;
  border-radius: 20px;
  letter-spacing: 0.4px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 7px;
  background: var(--bg-panel-2);
  border: 1px solid var(--line);
  border-radius: 6px;
  padding: 6px 10px;
  color: var(--text-dim);
  font-size: 12.5px;
  width: 240px;
}

.search-box kbd {
  margin-left: auto;
  font-family: var(--mono);
  font-size: 10px;
  background: #262C38;
  padding: 1px 5px;
  border-radius: 3px;
  color: var(--text-dim);
}

.identity {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  color: var(--text-mid);
}

.identity .avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: #2A3142;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  color: var(--text-hi);
}

.identity .ldap-tag {
  font-family: var(--mono);
  font-size: 10px;
  color: var(--text-dim);
}
</style>
