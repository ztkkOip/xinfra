<template>
  <div class="ext-card" @click="handleClick">
    <div class="ext-top">
      <div class="ext-logo" :style="{ background: bgColor, color: iconColor }">
        {{ system.icon }}
      </div>
      <div>
        <h4>{{ system.label }} {{ system.name }}</h4>
        <div class="ext-sub">{{ system.domain }}</div>
      </div>
    </div>
    <p>{{ system.description }}</p>
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
    W: 'var(--logo-w-bg)',
    DM: 'var(--logo-dm-bg)',
    CC: 'var(--logo-cc-bg)',
    AP: 'var(--logo-ap-bg)',
    QP: 'var(--logo-qp-bg)',
    GF: 'var(--logo-gf-bg)',
    SS: 'var(--logo-ss-bg)',
  }
  return colors[props.system.icon] || 'var(--bg-panel-2)'
})

const iconColor = computed(() => {
  const colors: Record<string, string> = {
    W: 'var(--tag-blue-text)',
    DM: 'var(--tag-green-text)',
    CC: 'var(--err)',
    AP: 'var(--tag-blue-text)',
    QP: 'var(--tag-amber-text)',
    GF: 'var(--warn)',
    SS: 'var(--tag-cyan-text)',
  }
  return colors[props.system.icon] || 'var(--text-mid)'
})

const handleClick = async () => {
  if (props.system.sso_enabled) {
    const { data } = await subsystemApi.getSSOUrl(props.system.id)
    window.location.assign(data.sso_url)
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
  border-color: var(--hover-border);
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
