<template>
  <div v-if="system">
    <div class="page-head">
      <div>
        <h1>{{ system.name }}</h1>
        <p>{{ system.description }} · {{ system.domain }}</p>
      </div>
      <span class="env-pill" :class="system.status === 'active' ? '' : 'warn'">
        {{ system.status === 'active' ? '● 在线' : '● 接入中' }}
      </span>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">接入状态</div>
        <div class="value" :class="system.status === 'active' ? 'accent' : 'warn'">
          {{ system.status === 'active' ? '在线' : '接入中' }}
        </div>
        <div class="delta">{{ system.category }}</div>
      </div>
      <div class="stat-card">
        <div class="label">SSO 登录</div>
        <div class="value" :class="system.sso_enabled ? 'accent' : ''">
          {{ system.sso_enabled ? '已启用' : '未启用' }}
        </div>
        <div class="delta">LDAP 原生配置 · 同账号同密码</div>
      </div>
      <div class="stat-card">
        <div class="label">访问域名</div>
        <div class="value" style="font-size:14px;">{{ system.domain }}</div>
        <div class="delta">内网 DNS 解析</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>统一入口</h3>
        <span class="meta">{{ system.name }} · LDAP 单点登录</span>
      </div>
      <div class="panel-body" style="padding:16px;">
        <div class="ext-card" style="max-width:560px;">
          <div class="ext-top">
            <div class="ext-logo" :style="{ background: bgColor, color: iconColor }">
              {{ system.icon }}
            </div>
            <div>
              <h4>{{ system.name }}</h4>
              <div class="ext-sub">{{ system.domain }}</div>
            </div>
          </div>
          <p>{{ detailDesc }}</p>
          <div class="sso-row">
            <span class="sso-dot" :class="{ warn: system.status !== 'active' }"></span>
            {{ system.status === 'active' ? 'LDAP 原生配置接入 · 在线' : 'LDAP 接入改造中' }}
          </div>
          <button
            class="btn btn-primary"
            :class="{ 'btn-disabled': system.status !== 'active' }"
            :disabled="system.status !== 'active'"
            style="margin-top:8px;align-self:flex-start;"
            @click="openPortal"
          >
            {{ system.status !== 'active' ? '暂未开放' : '打开 ' + system.name + ' ↗' }}
          </button>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="page-head">
    <div><h1>子系统未找到</h1><p>请从左侧导航选择一个子系统</p></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { subsystemApi, type Subsystem } from '@/api/subsystem'
import { getToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'

const route = useRoute()

const system = computed<Subsystem | undefined>(() => {
  const name = route.params.name as string
  return subsystemApi.getSubsystemByName(name)
})

const bgColor = computed(() => {
  if (!system.value) return 'var(--bg-panel-2)'
  const colors: Record<string, string> = {
    W: 'var(--logo-w-bg)',
    DM: 'var(--logo-dm-bg)',
    CC: 'var(--logo-cc-bg)',
    AP: 'var(--logo-ap-bg)',
    QP: 'var(--logo-qp-bg)',
    GF: 'var(--logo-gf-bg)',
    SS: 'var(--logo-ss-bg)',
  }
  return colors[system.value.icon] || 'var(--bg-panel-2)'
})

const iconColor = computed(() => {
  if (!system.value) return 'var(--text-mid)'
  const colors: Record<string, string> = {
    W: 'var(--tag-blue-text)',
    DM: 'var(--tag-green-text)',
    CC: 'var(--err)',
    AP: 'var(--tag-blue-text)',
    QP: 'var(--tag-amber-text)',
    GF: 'var(--warn)',
    SS: 'var(--tag-cyan-text)',
  }
  return colors[system.value.icon] || 'var(--text-mid)'
})

const detailDesc = computed(() => {
  if (!system.value) return ''
  const descs: Record<string, string> = {
    Wayne: '业务容器发布、命名空间与配额管理，复用 Wayne 原生多租户能力。所有集群的容器发布均可通过此入口统一操作。',
    CloudDM: '数据库 SQL 上线统一审核，支持 LDAP 用户组到角色的自动映射。变更记录可追溯，审计合规。',
    CacheCloud: 'Redis 实例全生命周期管理，登录与监控运维统一入口。按业务线隔离实例，统一密码管理。',
    Apollo: '统一管理多机房 Config Service Cluster 的配置发布、灰度、回滚与变更审计，所有机房配置项均可在此单一入口检索与编辑。',
    qpass: '七牛云内部系统。',
    Grafana: '对接 VictoriaMetrics 数据源，K8s / 容器 / 业务指标统一仪表盘展示。预置核心看板，支持自定义。',
    Superset: '日志数据探索与可视化分析平台，支持多数据源接入和丰富的图表展示。用于日志数据分析和可视化报表。',
  }
  return descs[system.value.name] || system.value.description
})

const openPortal = async () => {
  if (!system.value) return
  if (system.value.sso_enabled) {
    try {
      const { data } = await subsystemApi.getSSOUrl(system.value.id)
      window.location.assign(data.sso_url)
    } catch {
      // 错误已由 request 拦截器统一处理
    }
  } else {
    window.open(system.value.url, '_blank')
  }
}
</script>

<style scoped>
.env-pill {
  font-size: 13px;
  color: var(--accent);
  background: var(--accent-dim);
  padding: 4px 14px;
  border-radius: 20px;
  white-space: nowrap;
}

.env-pill.warn {
  color: var(--warn);
  background: rgba(255, 180, 0, 0.1);
}

.panel-body .ext-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.btn-primary {
  background: var(--accent);
  color: #fff;
  border: none;
  padding: 8px 18px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.15s;
}

.btn-primary:hover {
  opacity: 0.85;
}

.btn-disabled {
  background: var(--text-dim);
  cursor: not-allowed;
  opacity: 0.6;
}

.btn-disabled:hover {
  opacity: 0.6;
}
</style>
