<template>
  <div>
    <div class="page-head">
      <div>
        <h1>监控 / 日志 / 告警</h1>
        <p>VictoriaMetrics 指标 · Zabbix 硬件状态 · ELK 日志中心 · 夜莺（Nightingale）统一告警引擎对接 qpass</p>
      </div>
    </div>

    <div class="cols-2">
      <div class="panel">
        <div class="panel-head">
          <h3>VictoriaMetrics · 指标概览</h3>
          <span class="meta">K8s / 容器 / 业务指标</span>
        </div>
        <div class="panel-body">
          <table>
            <tr class="tr-hover">
              <td>七机房采集节点数</td>
              <td class="mono strong">128</td>
            </tr>
            <tr class="tr-hover">
              <td>活跃时序数量</td>
              <td class="mono strong">42.6M</td>
            </tr>
            <tr class="tr-hover">
              <td>Prometheus 接口接入</td>
              <td class="status-text ok">● 正常</td>
            </tr>
            <tr class="tr-hover">
              <td>Grafana 仪表盘数</td>
              <td class="mono strong">37</td>
            </tr>
          </table>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>Zabbix · 硬件健康</h3>
          <span class="meta">IPMI / 温度 / 电源 / 风扇</span>
        </div>
        <div class="panel-body">
          <table>
            <tr class="tr-hover">
              <td>物理机监控覆盖</td>
              <td class="mono strong">186 / 186</td>
            </tr>
            <tr class="tr-hover">
              <td>当前告警</td>
              <td class="status-text warn">● 2 条（风扇转速异常）</td>
            </tr>
            <tr class="tr-hover">
              <td>网络设备健康</td>
              <td class="status-text ok">● 正常</td>
            </tr>
            <tr class="tr-hover">
              <td>存储健康度</td>
              <td class="status-text ok">● 正常</td>
            </tr>
          </table>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>日志 + 统一告警</h3>
        <span class="meta">Filebeat → Kafka → ES → Kibana · VM + Zabbix → qpass</span>
      </div>
      <div class="panel-body log-stream">
        <div v-for="(log, index) in logs" :key="index" :class="['task-log-line', log.class]">
          <span class="t">{{ log.time }}</span>{{ log.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const logs = ref([
  { time: '07:39:10', message: 'CMDB 采集器同步完成：xs291 · 华东-杭州下沙机房', class: '' },
  { time: '07:38:55', message: '[Zabbix] 告警：yzh-las-014 风扇转速低于阈值（FAN_6: 1920 RPM）', class: '' },
  { time: '07:38:02', message: '[VictoriaMetrics] 七机房抓取任务全部正常', class: 'ok' },
  { time: '07:35:41', message: '[ELK] sg-las-007 应用日志接入 Kafka topic：las-sg-app', class: '' },
  { time: '07:30:00', message: '[qpass] 告警已合并推送：1 条硬件告警 + 0 条业务告警', class: '' },
])
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 18px;
}

.page-head h1 {
  font-size: 19px;
  margin: 0 0 4px;
  font-weight: 700;
}

.page-head p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12.5px;
}

.cols-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  margin-bottom: 18px;
}

.panel {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  margin-bottom: 18px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 13px 16px;
  border-bottom: 1px solid var(--line-soft);
}

.panel-head h3 {
  margin: 0;
  font-size: 13.5px;
  font-weight: 600;
}

.panel-head .meta {
  font-size: 11.5px;
  color: var(--text-dim);
  font-family: var(--mono);
}

.panel-body {
  padding: 4px 0;
}

.log-stream {
  padding: 14px 16px;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12.5px;
}

td {
  padding: 11px 16px;
  border-bottom: 1px solid var(--line-soft);
  color: var(--text-mid);
}

tr:last-child td {
  border-bottom: none;
}

tr.tr-hover:hover td {
  background: var(--bg-panel-2);
}

.strong {
  color: var(--text-hi);
  font-weight: 500;
}

.mono {
  font-family: var(--mono);
}

.status-text {
  font-size: 12px;
  display: flex;
  align-items: center;
}

.status-text.ok {
  color: var(--accent);
}

.status-text.warn {
  color: var(--warn);
}

.task-log-line {
  font-family: var(--mono);
  font-size: 11.5px;
  padding: 2px 0;
  color: var(--text-mid);
}

.task-log-line .t {
  color: var(--text-dim);
  margin-right: 10px;
}

.task-log-line.ok {
  color: var(--accent);
}
</style>
