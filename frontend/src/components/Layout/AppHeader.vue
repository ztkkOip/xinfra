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
    <div class="theme-switch" @click="toggleTheme" :title="theme === 'light' ? '切换到深色模式' : '切换到浅色模式'">
      <el-icon v-if="theme === 'light'" :size="18">
        <Moon />
      </el-icon>
      <el-icon v-else :size="18">
        <Sunny />
      </el-icon>
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
import { useTheme } from '@/composables/useTheme'
import { Moon, Sunny } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const userInitial = computed(() => {
  return user.value?.display_name?.charAt(0) || 'U'
})

const { theme, toggleTheme } = useTheme()
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
  color: var(--bg-deep);
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
  background: var(--accent-dim);
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
  background: var(--bg-panel);
  padding: 1px 5px;
  border-radius: 3px;
  color: var(--text-dim);
  border: 1px solid var(--line);
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
  background: var(--bg-panel-2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  color: var(--text-hi);
  border: 1px solid var(--line);
}

.identity .ldap-tag {
  font-family: var(--mono);
  font-size: 10px;
  color: var(--text-dim);
}

.theme-switch {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-mid);
  background: var(--bg-panel-2);
  border: 1px solid var(--line);
}

.theme-switch:hover {
  background: var(--bg-panel);
  color: var(--text-hi);
  border-color: var(--accent);
}

.theme-switch:active {
  transform: scale(0.95);
}

.theme-switch .el-icon {
  transition: transform 0.3s ease;
}

.theme-switch:hover .el-icon {
  transform: rotate(15deg);
}
</style>
