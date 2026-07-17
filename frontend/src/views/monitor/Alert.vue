<template>
  <div>
    <div class="page-head">
      <div>
        <h1>告警</h1>
        <p>夜莺（Nightingale）统一告警引擎 · qpass 推送</p>
      </div>
    </div>

    <div class="cols-2-equal">
      <div class="panel">
        <div class="panel-head">
          <h3>当前告警</h3>
          <span class="meta">实时</span>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
              <tr class="tr-hover">
                <td>硬件告警</td>
                <td class="status-text warn">● 2 条</td>
              </tr>
              <tr class="tr-hover">
                <td>业务告警</td>
                <td class="status-text ok">● 0 条</td>
              </tr>
              <tr class="tr-hover">
                <td>今日已推送</td>
                <td class="mono strong">5</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>告警通道</h3>
          <span class="meta">qpass</span>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
              <tr class="tr-hover">
                <td>推送状态</td>
                <td class="status-text ok">● 正常</td>
              </tr>
              <tr class="tr-hover">
                <td>通道类型</td>
                <td>企业微信 · 钉钉 · 短信</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>告警记录</h3>
        <span class="meta">近 24h</span>
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
  { time: '07:38:55', message: '[Zabbix] yzh-las-014 风扇转速低于阈值（FAN_6: 1920 RPM）', class: '' },
  { time: '07:38:55', message: '[Zabbix] yzh-las-014 风扇转速低于阈值（FAN_5: 1880 RPM）', class: '' },
  { time: '07:30:00', message: '[qpass] 告警已合并推送：1 条硬件告警 + 0 条业务告警', class: '' },
  { time: '06:15:22', message: '[Nightingale] sg-las-007 CPU 使用率恢复（当前 42%）', class: 'ok' },
  { time: '05:48:10', message: '[Nightingale] sg-las-007 CPU 使用率超过 90% 持续 5 分钟', class: '' },
])
</script>

<style scoped>
</style>
