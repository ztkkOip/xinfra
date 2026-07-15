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
              <td class="mono text-xs">{{ node.taint }}</td>
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
/* 公共样式已在 global.css 中定义 */

.tag.biz {
  color: var(--text-mid);
}
</style>
