<template>
  <aside class="sidebar">
    <div class="nav-group">
      <div class="nav-label">概览</div>
      <router-link to="/dashboard" class="nav-item" active-class="active">
        <span class="ic">▣</span>资源大盘
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">业务交付</div>
      <router-link to="/service/catalog" class="nav-item" active-class="active">
        <span class="ic">▤</span>基础服务
      </router-link>
      <router-link to="/task/log" class="nav-item" active-class="active">
        <span class="ic">▧</span>任务日志<span class="badge">2</span>
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">资源纳管</div>
      <router-link to="/machine/management" class="nav-item" active-class="active">
        <span class="ic">▥</span>机器管理<span class="badge">CMDB</span>
      </router-link>
      <router-link to="/service/management" class="nav-item" active-class="active">
        <span class="ic">◈</span>服务管理<span class="badge">Consul</span>
      </router-link>
      <router-link to="/infrastructure/cluster" class="nav-item" active-class="active">
        <span class="ic">◆</span>集群与容器<span class="badge">3</span>
      </router-link>
      <router-link to="/business/quota" class="nav-item" active-class="active">
        <span class="ic">▦</span>业务配额
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">可观测性</div>
      <router-link to="/observable/monitoring" class="nav-item" active-class="active">
        <span class="ic">◎</span>监控看板
      </router-link>
      <router-link to="/observable/alert" class="nav-item" active-class="active">
        <span class="ic">⚡</span>告警信息
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">统一入口</div>
      <div class="nav-item nav-expand" :class="{ active: isPortalActive }">
        <router-link to="/subsystem" class="expand-label">
          <span class="ic">↗</span>子系统导航
        </router-link>
        <span class="expand-toggle" @click="portalExpanded = !portalExpanded">
          {{ portalExpanded ? '▾' : '▸' }}
        </span>
      </div>
      <div v-show="portalExpanded" class="nav-children">
        <router-link
          v-for="sys in subsystems"
          :key="sys.name"
          :to="`/subsystem/detail/${sys.name.toLowerCase()}`"
          class="nav-item nav-child"
          active-class="active"
        >
          <span class="child-dot" :style="{ background: dotColor(sys.icon) }"></span>{{ sys.label }} {{ sys.name }}
        </router-link>
        <router-link to="/subsystem" class="nav-item nav-child" active-class="active">
          <span class="ic" style="font-size:12px;">▦</span>全部总览
        </router-link>
      </div>
      <router-link to="/subsystem/authorization" class="nav-item" active-class="active">
        <span class="ic">▦</span>子系统赋权
      </router-link>
      <router-link v-if="isPlatformAdmin" to="/business-line/manage" class="nav-item" active-class="active">
        <span class="ic">▤</span>业务线管理
      </router-link>
      <router-link v-if="isBusinessLineAdmin" to="/business-line/assignment" class="nav-item" active-class="active">
        <span class="ic">▥</span>业务线分配
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">审计</div>
      <router-link to="/audit/login" class="nav-item" active-class="active">
        <span class="ic">📋</span>登录审计
      </router-link>
      <router-link to="/audit/operation" class="nav-item" active-class="active">
        <span class="ic">📝</span>运维操作审计
      </router-link>
    </div>
    <div class="sidebar-foot">
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBusinessLineStore } from '@/stores/businessLine'

const route = useRoute()
const authStore = useAuthStore()
const businessLineStore = useBusinessLineStore()
const portalExpanded = ref(true)
const isPlatformAdmin = computed(() => authStore.isAdmin)
const isBusinessLineAdmin = computed(() => businessLineStore.isCurrentAdmin)

const subsystems = [
  { name: 'Wayne', label: '多集群发布平台', icon: 'W' },
  { name: 'CloudDM', label: 'SQL 审核平台', icon: 'DM' },
  { name: 'CacheCloud', label: 'Redis 管理平台', icon: 'CC' },
  { name: 'Apollo', label: '配置中心', icon: 'AP' },
  { name: 'qpass', label: '七牛统一运维平台', icon: 'QP' },
  { name: 'Grafana', label: '指标可视化平台', icon: 'GF' },
  { name: 'Superset', label: '日志数据可视化', icon: 'SS' },
]

const isPortalActive = computed(() => {
  return route.path.startsWith('/subsystem/detail/') || route.path === '/subsystem'
})

function dotColor(icon: string): string {
  const colors: Record<string, string> = {
    W: 'var(--tag-blue-text)',
    DM: 'var(--tag-green-text)',
    CC: 'var(--err)',
    AP: 'var(--tag-blue-text)',
    QP: 'var(--tag-amber-text)',
    GF: 'var(--warn)',
    SS: 'var(--tag-cyan-text)',
  }
  return colors[icon] || 'var(--text-dim)'
}
</script>

<style scoped>
.sidebar {
  background: var(--bg-panel);
  border-right: 1px solid var(--line);
  padding: 14px 10px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.nav-group {
  margin-bottom: 18px;
}

.nav-label {
  font-size: 13px;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 0 10px;
  margin-bottom: 6px;
  font-weight: 600;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px;
  border-radius: 6px;
  color: var(--text-mid);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 2px;
  border: 1px solid transparent;
  text-decoration: none;
}

.nav-item .ic {
  width: 16px;
  text-align: center;
  font-size: 16px;
  opacity: 0.85;
}

.nav-item:hover {
  background: var(--bg-panel-2);
  color: var(--text-hi);
}

.nav-item.active {
  background: var(--accent-dim);
  color: var(--accent);
  border-color: var(--accent);
}

.nav-item.active .expand-label {
  color: var(--accent);
}

.nav-item .badge {
  margin-left: auto;
  font-family: var(--mono);
  font-size: 12px;
  background: var(--bg-panel-2);
  color: var(--text-dim);
  padding: 1px 6px;
  border-radius: 10px;
  border: 1px solid var(--line);
}

.nav-item.active .badge {
  background: var(--accent);
  color: #FFFFFF;
  border-color: var(--accent);
}

.sidebar-foot {
  margin-top: auto;
  padding: 10px;
  border-top: 1px solid var(--line-soft);
  font-size: 14px;
  color: var(--text-dim);
  line-height: 1.6;
}

.nav-expand {
  justify-content: flex-start;
  padding: 0;
  user-select: none;
}

.expand-label {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  padding: 10px;
  text-decoration: none;
  color: inherit;
}

.expand-toggle {
  padding: 10px;
  font-size: 12px;
  opacity: 0.6;
  cursor: pointer;
  border-radius: 0 6px 6px 0;
}

.expand-toggle:hover {
  opacity: 1;
  background: var(--bg-panel-2);
}

.nav-children {
  padding-left: 8px;
}

.nav-child {
  padding: 7px 10px 7px 26px;
  font-size: 14px;
  gap: 8px;
}

.child-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>
