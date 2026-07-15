<template>
  <div>
    <div class="page-head">
      <div>
        <h1>集群 / 节点</h1>
        <p>按机房切分的 RKE2 集群，Calico BGP 扁平网络</p>
      </div>
      <el-button type="primary">+ 加入新节点</el-button>
    </div>

    <div class="panel">
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>集群</th>
              <th>机房 / 区域</th>
              <th>状态</th>
              <th>节点</th>
              <th>CPU 使用</th>
              <th>RKE2 版本</th>
              <th>Calico</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cluster in clusters" :key="cluster.name" class="tr-hover">
              <td class="strong">{{ cluster.name }}</td>
              <td><span class="tag" :class="cluster.zoneClass">{{ cluster.zone }}</span></td>
              <td :class="['status-text', cluster.statusClass]">● {{ cluster.status }}</td>
              <td class="mono">{{ cluster.nodes }}</td>
              <td>
                <span class="bar-wrap"><span class="bar-fill" :class="cluster.cpuClass" :style="{ width: cluster.cpu + '%' }"></span></span>{{ cluster.cpu }}%
              </td>
              <td class="mono">{{ cluster.version }}</td>
              <td class="mono">{{ cluster.calico }}</td>
              <td><el-button text>详情</el-button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>rke2-bj-prod-01 · 节点列表（节选）</h3>
        <span class="meta">node-label 多租户隔离</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>节点</th>
              <th>IP</th>
              <th>业务线标签</th>
              <th>污点 Taint</th>
              <th>规格</th>
              <th>CPU / Mem</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="node in nodes" :key="node.name" class="tr-hover">
              <td class="strong mono">{{ node.name }}</td>
              <td class="mono">{{ node.ip }}</td>
              <td><span class="tag biz" :style="{ color: node.labelColor, borderColor: node.labelBorder }">{{ node.label }}</span></td>
              <td class="mono" style="font-size: 11px">{{ node.taint }}</td>
              <td class="mono">{{ node.spec }}</td>
              <td><span class="bar-wrap"><span class="bar-fill" :class="node.cpuClass" :style="{ width: node.cpu + '%' }"></span></span>{{ node.cpu }}%</td>
              <td :class="['status-text', node.statusClass]">● {{ node.status }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const clusters = ref([
  { name: 'rke2-bj-prod-01', zone: '华北 · IDC', zoneClass: 'zone-a', status: '健康', statusClass: 'ok', nodes: 46, cpu: 64, cpuClass: '', version: 'v1.28.9+rke2r1', calico: 'BGP AS64512' },
  { name: 'rke2-sh-prod-01', zone: '华东 · IDC', zoneClass: 'zone-b', status: '健康', statusClass: 'ok', nodes: 52, cpu: 71, cpuClass: '', version: 'v1.28.9+rke2r1', calico: 'BGP AS64513' },
  { name: 'rke2-ali-sz-01', zone: '阿里云 · 华南', zoneClass: 'zone-c', status: '1 节点告警', statusClass: 'warn', nodes: 30, cpu: 83, cpuClass: 'warn', version: 'v1.28.6+rke2r1', calico: 'BGP AS64514' },
])

const nodes = ref([
  { name: 'bj-node-014', ip: '10.21.4.14', label: 'business-line=kodo', labelColor: 'var(--tag-blue-text)', labelBorder: 'var(--tag-blue-border)', taint: 'kodo:NoSchedule', spec: '64C/256G', cpu: 58, cpuClass: '', status: 'Ready', statusClass: 'ok' },
  { name: 'bj-node-015', ip: '10.21.4.15', label: 'business-line=kodo', labelColor: 'var(--tag-blue-text)', labelBorder: 'var(--tag-blue-border)', taint: 'kodo:NoSchedule', spec: '64C/256G', cpu: 62, cpuClass: '', status: 'Ready', statusClass: 'ok' },
  { name: 'bj-node-031', ip: '10.21.4.31', label: 'business-line=las', labelColor: 'var(--tag-purple-text)', labelBorder: 'var(--tag-purple-border)', taint: 'las:NoSchedule', spec: '32C/128G', cpu: 88, cpuClass: 'warn', status: '资源告警', statusClass: 'warn' },
  { name: 'bj-node-048', ip: '10.21.4.48', label: '未分配', labelColor: 'var(--text-dim)', labelBorder: 'var(--line)', taint: '—', spec: '32C/128G', cpu: 4, cpuClass: '', status: 'Ready · 空闲', statusClass: 'ok' },
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

.tag.biz {
  color: var(--text-mid);
}

.bar-wrap {
  width: 90px;
  height: 5px;
  background: var(--bar-bg);
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
</style>
