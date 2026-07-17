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
        <div class="value" style="color: var(--err)">1</div>
        <div class="delta" style="color: var(--err)">来自 VictoriaMetrics</div>
      </div>
      <div class="stat-card">
        <div class="label">P1 / High 告警</div>
        <div class="value" style="color: var(--warn)">5</div>
        <div class="delta warn">来自 Zabbix · 4 条 / VM · 1 条</div>
      </div>
      <div class="stat-card">
        <div class="label">今日已推送</div>
        <div class="value">5</div>
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
          <tr class="tr-hover">
            <td><span class="tag sev-p0">P0</span></td>
            <td class="mono">VictoriaMetrics</td>
            <td class="mono">disaster</td>
            <td class="strong mono">sg-las-overseas-007</td>
            <td>磁盘只读 / IO 错误，节点不可写</td>
            <td class="mono">07:38:55</td>
          </tr>
          <tr class="tr-hover">
            <td><span class="tag sev-p1">P1</span></td>
            <td class="mono">Zabbix</td>
            <td class="mono">high</td>
            <td class="strong mono">yzh-las-014</td>
            <td>风扇转速低于阈值（FAN_6: 1920 RPM）</td>
            <td class="mono">07:35:12</td>
          </tr>
          <tr class="tr-hover">
            <td><span class="tag sev-p1">P1</span></td>
            <td class="mono">Zabbix</td>
            <td class="mono">high</td>
            <td class="strong mono">xs-bm-291</td>
            <td>电源冗余丢失（PSU2 离线）</td>
            <td class="mono">07:21:40</td>
          </tr>
          <tr class="tr-hover">
            <td><span class="tag sev-p1">P1</span></td>
            <td class="mono">VictoriaMetrics</td>
            <td class="mono">high</td>
            <td class="strong mono">redis-ads-feature-cluster-02</td>
            <td>慢查询比例超阈值（>5%，持续5分钟）</td>
            <td class="mono">07:18:03</td>
          </tr>
          <tr class="tr-hover">
            <td><span class="tag sev-avg">average</span></td>
            <td class="mono">Zabbix</td>
            <td class="mono">average</td>
            <td class="strong mono">overseas-dallas-idc</td>
            <td>专线延迟抖动超基线（WireGuard ↔ AWS US）</td>
            <td class="mono">07:10:22</td>
          </tr>
          <tr class="tr-hover">
            <td><span class="tag sev-avg">average</span></td>
            <td class="mono">Zabbix</td>
            <td class="mono">average</td>
            <td class="strong mono">jf-bm-118</td>
            <td>CPU 温度接近阈值（76℃ / 80℃）</td>
            <td class="mono">06:58:47</td>
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
import { ref } from 'vue'

const alerts = ref([
  { time: '07:38:55', message: '[qpass] P0 告警推送：sg-las-overseas-007 磁盘只读 / IO 错误', class: 'err' },
  { time: '07:35:12', message: '[qpass] P1 告警推送：yzh-las-014 风扇转速低于阈值（FAN_6: 1920 RPM）', class: '' },
  { time: '07:21:40', message: '[qpass] P1 告警推送：xs-bm-291 电源冗余丢失（PSU2 离线）', class: '' },
  { time: '07:18:03', message: '[qpass] P1 告警推送：redis-ads-feature-cluster-02 慢查询比例超阈值', class: '' },
  { time: '07:10:22', message: '[qpass] average 告警推送：overseas-dallas-idc 专线延迟抖动超基线', class: '' },
  { time: '06:58:47', message: '[qpass] average 告警推送：jf-bm-118 CPU 温度接近阈值（76℃ / 80℃）', class: '' },
  { time: '06:15:22', message: '[Nightingale] sg-las-007 CPU 使用率恢复（当前 42%）', class: 'ok' },
  { time: '05:48:10', message: '[Nightingale] sg-las-007 CPU 使用率超过 90% 持续 5 分钟', class: '' },
])
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
