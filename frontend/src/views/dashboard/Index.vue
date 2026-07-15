<template>
  <div>
    <div class="page-head">
      <div>
        <h1>资源大盘</h1>
        <p>跨机房 / 多云容器资源统一视图 · 最近更新于 12 秒前</p>
      </div>
      <el-button @click="refresh">⟳ 刷新</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">RKE2 集群</div>
        <div class="value">3</div>
        <div class="delta">3 机房在线</div>
      </div>
      <div class="stat-card">
        <div class="label">节点总数</div>
        <div class="value">128</div>
        <div class="delta up">↑ 6 本周新增</div>
      </div>
      <div class="stat-card">
        <div class="label">CPU 已分配</div>
        <div class="value">61%</div>
        <div class="delta">2,048 / 3,360 核</div>
      </div>
      <div class="stat-card">
        <div class="label">组件实例</div>
        <div class="value">214</div>
        <div class="delta">MySQL 38 · Redis 92 · 其他 84</div>
      </div>
      <div class="stat-card">
        <div class="label">进行中任务</div>
        <div class="value" style="color: var(--warn)">2</div>
        <div class="delta warn">1 个待关注</div>
      </div>
    </div>

    <div class="cols-2">
      <div class="panel">
        <div class="panel-head">
          <h3>集群拓扑 · 按机房 / 业务线标签</h3>
          <span class="meta">node-label 调度视图</span>
        </div>
        <div class="topo">
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-a">IDC-华北机房</span>
              <span class="mono" style="color: var(--text-dim)">46节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 24" :key="i" :class="['node-cell', i <= 18 ? 'a' : '']"></div>
            </div>
          </div>
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-b">IDC-华东机房</span>
              <span class="mono" style="color: var(--text-dim)">52节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 24" :key="i" :class="['node-cell', i <= 20 ? 'b' : '']"></div>
            </div>
          </div>
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-c">阿里云-华南</span>
              <span class="mono" style="color: var(--text-dim)">30节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 18" :key="i" :class="['node-cell', i <= 12 ? 'c' : '']"></div>
            </div>
          </div>
        </div>
        <div class="legend">
          <span><span class="node-cell a inline"></span>Kodo业务线</span>
          <span><span class="node-cell b inline"></span>LAS业务线</span>
          <span><span class="node-cell c inline"></span>灵矽业务线</span>
          <span><span class="node-cell empty inline"></span>空闲</span>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>最近任务</h3>
          <router-link to="/tasks" class="meta">查看全部 →</router-link>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">MySQL 主从部署</td>
                <td class="mono">2m14s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text warn">● 执行中</td>
                <td class="strong">RKE2 节点加入</td>
                <td class="mono">38s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">Redis Cluster 部署</td>
                <td class="mono">3m02s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text err">● 失败</td>
                <td class="strong">openresty 网关部署</td>
                <td class="mono">41s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">业务线标签同步</td>
                <td class="mono">5s</td>
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
  position: relative;
  overflow: hidden;
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
  grid-template-columns: 1.4fr 1fr;
  gap: 18px;
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
  text-decoration: none;
}

.panel-body {
  padding: 4px 0;
}

.topo {
  padding: 18px 16px;
  display: flex;
  gap: 14px;
}

.topo-cluster {
  flex: 1;
  border: 1px dashed var(--line);
  border-radius: 8px;
  padding: 12px;
}

.topo-cluster .tname {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.node-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 5px;
}

.node-cell {
  aspect-ratio: 1;
  border-radius: 3px;
  background: var(--node-empty-bg);
  border: 1px solid var(--line-soft);
}

.node-cell.a {
  background: var(--node-a-bg);
  border-color: var(--node-a-border);
}

.node-cell.b {
  background: var(--node-b-bg);
  border-color: var(--node-b-border);
}

.node-cell.c {
  background: var(--node-c-bg);
  border-color: var(--node-c-border);
}

.node-cell.empty {
  background: transparent;
  border-style: dashed;
}

.node-cell.inline {
  display: inline-block;
  width: 9px;
  height: 9px;
  margin-right: 5px;
}

.legend {
  padding: 0 16px 14px;
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: var(--text-dim);
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

td.strong {
  color: var(--text-hi);
  font-weight: 500;
}

.mono {
  font-family: var(--mono);
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-family: var(--mono);
  font-size: 10.5px;
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid var(--line);
  color: var(--text-mid);
}

.tag.zone-a {
  color: var(--tag-blue-text);
  border-color: var(--tag-blue-border);
  background: var(--tag-blue-bg);
}

.tag.zone-b {
  color: var(--tag-purple-text);
  border-color: var(--tag-purple-border);
  background: var(--tag-purple-bg);
}

.tag.zone-c {
  color: var(--tag-amber-text);
  border-color: var(--tag-amber-border);
  background: var(--tag-amber-bg);
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
