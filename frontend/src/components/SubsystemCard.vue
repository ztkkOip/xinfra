<template>
  <div class="ext-card" @click="handleClick">
    <div class="ext-top">
      <div class="ext-logo" :style="{ background: bgColor, color: iconColor }">
        {{ system.icon }}
      </div>
      <div>
        <h4>{{ system.name }}</h4>
        <div class="ext-sub">{{ system.description }}</div>
      </div>
    </div>
    <p>{{ description }}</p>
    <div class="sso-row">
      <span class="sso-dot" :class="{ warn: system.status === 'integrating' }"></span>
      {{ system.status === 'integrated' ? 'LDAP 原生配置接入 · 零代码' : 'LDAP 接入改造中' }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { subsystemApi, type Subsystem } from '@/api/subsystem'

const props = defineProps<{
  system: Subsystem
}>()

const bgColor = computed(() => {
  const colors: Record<string, string> = {
    W: '#1C2A3A',
    DM: '#1C3A2E',
    CC: '#3A1C28',
    AP: '#1C2A3A',
    QP: '#241B0A',
    GF: '#1B202B',
  }
  return colors[props.system.icon] || '#262C38'
})

const iconColor = computed(() => {
  const colors: Record<string, string> = {
    W: '#7FB8FF',
    DM: '#7FFFC2',
    CC: '#FF8E9C',
    AP: '#8EC8FF',
    QP: '#FFC97A',
    GF: '#F2B544',
  }
  return colors[props.system.icon] || '#9AA3B5'
})

const description = computed(() => {
  const descs: Record<string, string> = {
    Wayne: '业务容器发布、命名空间与配额管理，复用 Wayne 原生多租户能力。',
    CloudDM: '数据库 SQL 上线统一审核，支持 LDAP 用户组到角色的自动映射。',
    CacheCloud: 'Redis 实例全生命周期管理，登录与监控运维统一入口。',
    Apollo: '统一配置管理，按机房 Cluster 隔离，支持灰度发布、版本回滚与变更审计。',
    qpass: '夜莺、Zabbix、VictoriaMetrics 告警统一收敛与值班通知，按 P0/P1 分级推送。',
    Grafana: '对接 VictoriaMetrics 数据源，K8s / 容器 / 业务指标统一仪表盘展示。',
  }
  return descs[props.system.name] || ''
})

const handleClick = async () => {
  if (props.system.sso_enabled) {
    const { data } = await subsystemApi.getSSOUrl(props.system.id)
    window.open(data.sso_url, '_blank')
  } else {
    window.open(props.system.url, '_blank')
  }
}
</script>

<style scoped>
.ext-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.ext-card:hover {
  border-color: #3A4356;
}

.ext-top {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ext-logo {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-weight: 700;
  font-size: 13px;
}

h4 {
  margin: 0;
  font-size: 14px;
}

.ext-sub {
  font-size: 11px;
  color: var(--text-dim);
  font-family: var(--mono);
}

p {
  font-size: 12px;
  color: var(--text-mid);
  line-height: 1.6;
  margin: 0;
}

.sso-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--text-dim);
  margin-top: auto;
}

.sso-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
}

.sso-dot.warn {
  background: var(--warn);
}
</style>
