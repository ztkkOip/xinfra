<template>
  <div>
    <div class="page-head">
      <div>
        <h1>资源状态看板</h1>
        <p>物理机 / 虚机 / 基础服务状态汇总</p>
      </div>
      <el-button @click="refresh">⟳ 刷新</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">物理机总数</div>
        <div class="value">186</div>
        <div class="delta">186 / 186 监控覆盖</div>
      </div>
      <div class="stat-card">
        <div class="label">虚机总数（LAS）</div>
        <div class="value">512</div>
        <div class="delta up">↑ 14 本周新增</div>
      </div>
      <div class="stat-card">
        <div class="label">基础组件实例</div>
        <div class="value">214</div>
        <div class="delta">MySQL 38 · Redis 92 · 其他 84</div>
      </div>
      <div class="stat-card">
        <div class="label">P0 / Disaster 告警</div>
        <div class="value err">1</div>
        <div class="delta err">来自 VictoriaMetrics</div>
      </div>
      <div class="stat-card">
        <div class="label">P1 / High 告警</div>
        <div class="value warn">5</div>
        <div class="delta warn">来自 Zabbix · 4 条 / VM · 1 条</div>
      </div>
    </div>

    <div class="cols-2-equal">
      <div class="panel">
        <div class="panel-head">
          <h3>物理机状态 · 按机房</h3>
          <span class="meta">数据源：Zabbix</span>
        </div>
        <div class="panel-body">
          <table>
            <thead>
              <tr><th>机房</th><th>总数</th><th>在线</th><th>健康率</th><th>状态</th></tr>
            </thead>
            <tbody>
              <tr class="tr-hover">
                <td class="strong mono">YZH 机房</td>
                <td class="mono">78</td>
                <td class="mono">78</td>
                <td><span class="bar-wrap"><span class="bar-fill" style="width: 98%"></span></span>98%</td>
                <td class="status-text ok">● 正常</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">XS 机房</td>
                <td class="mono">64</td>
                <td class="mono">63</td>
                <td><span class="bar-wrap"><span class="bar-fill warn" style="width: 91%"></span></span>91%</td>
                <td class="status-text warn">● 1 条告警</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">JF 机房</td>
                <td class="mono">22</td>
                <td class="mono">22</td>
                <td><span class="bar-wrap"><span class="bar-fill" style="width: 100%"></span></span>100%</td>
                <td class="status-text ok">● 正常</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">海外 IDC（4）</td>
                <td class="mono">22</td>
                <td class="mono">21</td>
                <td><span class="bar-wrap"><span class="bar-fill err" style="width: 86%"></span></span>86%</td>
                <td class="status-text err">● 1 条严重</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>虚机状态 · LAS 资源池</h3>
          <span class="meta">数据源：VictoriaMetrics</span>
        </div>
        <div class="panel-body">
          <table>
            <thead>
              <tr><th>业务线</th><th>虚机数</th><th>CPU 均值</th><th>状态</th></tr>
            </thead>
            <tbody>
              <tr class="tr-hover">
                <td class="strong mono">Kodo</td>
                <td class="mono">218</td>
                <td><span class="bar-wrap"><span class="bar-fill" style="width: 64%"></span></span>64%</td>
                <td class="status-text ok">● 正常</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">LAS</td>
                <td class="mono">164</td>
                <td><span class="bar-wrap"><span class="bar-fill warn" style="width: 82%"></span></span>82%</td>
                <td class="status-text warn">● 1 条告警</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">灵矽</td>
                <td class="mono">96</td>
                <td><span class="bar-wrap"><span class="bar-fill" style="width: 57%"></span></span>57%</td>
                <td class="status-text ok">● 正常</td>
              </tr>
              <tr class="tr-hover">
                <td class="strong mono">共享池 / 未分配</td>
                <td class="mono">34</td>
                <td><span class="bar-wrap"><span class="bar-fill" style="width: 31%"></span></span>31%</td>
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
const refresh = () => {
  // 刷新数据
}
</script>

<style scoped>
/* 所有样式已在 global.css 中定义 */
</style>
