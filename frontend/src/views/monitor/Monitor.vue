<template>
  <div>
    <div class="page-head">
      <div>
        <h1>监控</h1>
        <p>VictoriaMetrics 指标 · Zabbix 硬件状态</p>
      </div>
    </div>

    <div class="cols-2-equal">
      <div class="panel">
        <div class="panel-head">
          <h3>VictoriaMetrics · 指标概览</h3>
          <span class="meta">K8s / 容器 / 业务指标</span>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
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
            </tbody>
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
            <tbody>
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
            </tbody>
          </table>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
/* .page-head, .cols-2-equal, .panel, table, .status-text 已在 global.css 中定义 */
</style>
