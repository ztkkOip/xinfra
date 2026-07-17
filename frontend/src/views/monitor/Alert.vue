<template>
  <div>
    <div class="page-head">
      <div>
        <h1>告警信息</h1>
        <p>夜莺（Nightingale）统一告警引擎对接 Zabbix + VictoriaMetrics · qpass 推送</p>
      </div>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">P0 / Disaster 告警</div>
        <div class="value" style="color: var(--err)">{{ profile.alertsP0 }}</div>
        <div class="delta" style="color: var(--err)">来自 VictoriaMetrics</div>
      </div>
      <div class="stat-card">
        <div class="label">P1 / High 告警</div>
        <div class="value" style="color: var(--warn)">{{ profile.alertsP1 }}</div>
        <div class="delta warn">来自 Zabbix · {{ Math.max(0, profile.alertsP1 - 1) }} 条 / VM · {{ profile.alertsP1 ? 1 : 0 }} 条</div>
      </div>
      <div class="stat-card">
        <div class="label">今日已推送</div>
        <div class="value">{{ profile.alertsP0 + profile.alertsP1 }}</div>
        <div class="delta">qpass 统一告警通道</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>告警来源接入 · 夜莺统一告警</h3>
        <span class="meta">Nightingale ← Zabbix / VictoriaMetrics</span>
      </div>
      <div class="panel-body" style="padding: 14px 16px">
        <div class="dc-row dc-row-3" style="margin-bottom: 0">
          <div class="dc-card dc-card-flat">
            <div class="dc-name">Zabbix <span class="dc-role">物理机 / 硬件层</span></div>
            <div class="dc-item">IPMI · 温度 · 电源 · 风扇</div>
            <div class="dc-item">网络设备 · 存储健康度</div>
            <div class="dc-item">事件级别：disaster / high / average</div>
          </div>
          <div class="dc-card dc-card-flat">
            <div class="dc-name">VictoriaMetrics <span class="dc-role">虚机 / 容器 / 业务层</span></div>
            <div class="dc-item">K8s / 容器 / 主机指标</div>
            <div class="dc-item">基础服务可用性探活</div>
            <div class="dc-item">告警级别：P0 / P1 / P2</div>
          </div>
          <div class="dc-card dc-card-flat">
            <div class="dc-name">夜莺 Nightingale <span class="dc-role">统一告警引擎</span></div>
            <div class="dc-item">按 P0/P1 或 disaster/high/average 过滤聚合</div>
            <div class="dc-item">去重 · 收敛 · 分级推送</div>
            <div class="dc-item">下游：qpass 统一告警通道</div>
          </div>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>当前告警 · P0/P1（disaster / high / average）</h3>
        <span class="meta">夜莺聚合 · 已去重</span>
      </div>
      <div class="panel-body">
        <table>
          <tr>
            <th>级别</th>
            <th>来源</th>
            <th>原始级别</th>
            <th>对象</th>
            <th>告警内容</th>
            <th>时间</th>
          </tr>
          <tr v-for="alert in alertRows" :key="alert.object" class="tr-hover">
            <td><span :class="['tag', alert.sevClass]">{{ alert.level }}</span></td>
            <td class="mono">{{ alert.source }}</td>
            <td class="mono">{{ alert.rawLevel }}</td>
            <td class="strong mono">{{ alert.object }}</td>
            <td>{{ alert.message }}</td>
            <td class="mono">{{ alert.time }}</td>
          </tr>
        </table>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>告警推送记录</h3>
        <span class="meta">近 24h · qpass</span>
      </div>
      <div class="panel-body log-stream">
          <div v-for="(alert, index) in alerts" :key="index" :class="['task-log-line', alert.class]">
          <span class="t">{{ alert.time }}</span>{{ alert.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBusinessLineMockProfile } from '@/utils/businessLineMock'

const { currentName, profile } = useBusinessLineMockProfile()

const alertRows = computed(() => {
  const rows = []
  if (profile.value.alertsP0) {
    rows.push({ level: 'P0', sevClass: 'sev-p0', source: 'VictoriaMetrics', rawLevel: 'disaster', object: `${currentName.value}-node-critical-007`, message: '磁盘只读 / IO 错误，节点不可写', time: '07:38:55' })
  }
  for (let i = 0; i < profile.value.alertsP1; i += 1) {
    rows.push({ level: 'P1', sevClass: 'sev-p1', source: i % 2 ? 'VictoriaMetrics' : 'Zabbix', rawLevel: 'high', object: `${currentName.value}-svc-${String(i + 1).padStart(2, '0')}`, message: i % 2 ? '慢查询比例超阈值（>5%，持续5分钟）' : '节点硬件健康指标异常', time: `07:${35 - i * 4}:12` })
  }
  if (!rows.length) {
    rows.push({ level: 'average', sevClass: 'sev-avg', source: 'Zabbix', rawLevel: 'average', object: `${currentName.value}-baseline`, message: '当前无 P0/P1，仅保留基线观察项', time: '07:10:22' })
  }
  return rows
})

const alerts = computed(() => alertRows.value.map((alert) => ({
  time: alert.time,
  message: `[qpass] ${alert.level} 告警推送：${alert.object} ${alert.message}`,
  class: alert.level === 'P0' ? 'err' : '',
})))
</script>

<style scoped>
.dc-row {
  display: flex;
  gap: 16px;
}

.dc-row-3 {
  flex-wrap: wrap;
}

.dc-card {
  flex: 1;
  min-width: 200px;
  padding: 16px;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 6px;
}

.dc-card-flat {
  background: var(--bg);
}

.dc-name {
  font-weight: 600;
  margin-bottom: 8px;
}

.dc-role {
  font-size: 12px;
  color: var(--text-dim);
  margin-left: 8px;
}

.dc-item {
  font-size: 13px;
  color: var(--text-dim);
  line-height: 1.8;
}

.sev-p0 {
  background: var(--err) !important;
  color: #fff !important;
  border-color: var(--err) !important;
}

.sev-p1 {
  background: var(--warn) !important;
  color: #000 !important;
  border-color: var(--warn) !important;
}

.sev-avg {
  background: var(--card-bg) !important;
  color: var(--warn) !important;
  border-color: var(--warn) !important;
}
</style>
