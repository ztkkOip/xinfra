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
        <div class="value" style="color: var(--err)">1</div>
        <div class="delta" style="color: var(--err)">来自 VictoriaMetrics</div>
      </div>
      <div class="stat-card">
        <div class="label">P1 / High 告警</div>
        <div class="value" style="color: var(--warn)">5</div>
        <div class="delta warn">来自 Zabbix · 4 条 / VM · 1 条</div>
      </div>
    </div>

    <div class="cols-2">
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
  letter-spacing: 0.2px;
}

.page-head p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12.5px;
}

.stat-row {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 14px 16px;
}

.stat-card .label {
  font-size: 11px;
  color: var(--text-dim);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-card .value {
  font-family: var(--mono);
  font-size: 24px;
  font-weight: 600;
}

.stat-card .delta {
  font-size: 11px;
  color: var(--text-dim);
  margin-top: 4px;
}

.stat-card .delta.up {
  color: var(--accent);
}

.stat-card .delta.warn {
  color: var(--warn);
}

.cols-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.panel {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
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

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12.5px;
}

th {
  text-align: left;
  color: var(--text-dim);
  font-weight: 500;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  padding: 9px 16px;
  border-bottom: 1px solid var(--line-soft);
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

.bar-wrap {
  width: 90px;
  height: 5px;
  background: #262C38;
  border-radius: 3px;
  overflow: hidden;
  display: inline-block;
  vertical-align: middle;
  margin-right: 8px;
}

.bar-fill {
  height: 100%;
  border-radius: 3px;
  background: var(--accent);
}

.bar-fill.warn {
  background: var(--warn);
}

.bar-fill.err {
  background: var(--err);
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

.status-text.err {
  color: var(--err);
}
</style>
