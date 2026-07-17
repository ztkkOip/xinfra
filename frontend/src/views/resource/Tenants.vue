<template>
  <div>
    <div class="page-head">
      <div>
        <h1>业务配额</h1>
        <p>按业务线划分资源配额与隔离策略 · RBAC 权限随业务线自动具备</p>
      </div>
      <el-button type="primary">+ 新增业务线</el-button>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>业务线列表</h3>
        <span class="meta">按业务线隔离命名空间、节点标签与配额</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>业务线</th>
              <th>标识</th>
              <th>命名空间</th>
              <th>节点数</th>
              <th>CPU 配额</th>
              <th>内存配额</th>
              <th>状态</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tenant in tenants" :key="tenant.key" class="tr-hover">
              <td class="strong">{{ tenant.name }}</td>
              <td class="mono">{{ tenant.key }}</td>
              <td class="mono">{{ tenant.namespace }}</td>
              <td class="mono">{{ tenant.nodes }}</td>
              <td>
                <span class="bar-wrap">
                  <span class="bar-fill" :class="tenant.cpuClass" :style="{ width: tenant.cpuPct + '%' }"></span>
                </span>
                {{ tenant.cpuPct }}%
              </td>
              <td>
                <span class="bar-wrap">
                  <span class="bar-fill" :class="tenant.memClass" :style="{ width: tenant.memPct + '%' }"></span>
                </span>
                {{ tenant.memPct }}%
              </td>
              <td :class="['status-text', tenant.statusClass]">● {{ tenant.status }}</td>
              <td>
                <el-button text>详情</el-button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="tooltip-note" style="max-width: 720px;">
      <b>配额与隔离策略：</b>每个业务线绑定独立的 K8s namespace 与 node-label，资源配额通过 ResourceQuota 限制 CPU / 内存上限。
      RBAC 角色（viewer / operator / admin）随业务线自动具备，子系统（Wayne / CloudDM）凭 SSO 身份免登进入。
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const tenants = ref([
  {
    name: 'LAS',
    key: 'las',
    namespace: 'ns-las-prod',
    nodes: 12,
    cpuPct: 64,
    cpuClass: '',
    memPct: 71,
    memClass: '',
    status: '正常',
    statusClass: 'ok',
  },
  {
    name: 'KODO',
    key: 'kodo',
    namespace: 'ns-kodo-prod',
    nodes: 8,
    cpuPct: 45,
    cpuClass: '',
    memPct: 52,
    memClass: '',
    status: '正常',
    statusClass: 'ok',
  },
  {
    name: 'PILI',
    key: 'pili',
    namespace: 'ns-pili-prod',
    nodes: 6,
    cpuPct: 82,
    cpuClass: 'warn',
    memPct: 78,
    memClass: '',
    status: '资源告警',
    statusClass: 'warn',
  },
  {
    name: 'QVM',
    key: 'qvm',
    namespace: 'ns-qvm-prod',
    nodes: 4,
    cpuPct: 23,
    cpuClass: '',
    memPct: 31,
    memClass: '',
    status: '正常',
    statusClass: 'ok',
  },
])
</script>

<style scoped>
/* 公共样式已在 global.css 中定义 */
</style>
