<template>
  <div>
    <div class="page-head">
      <div>
        <h1>配置中心管理</h1>
        <p>Apollo 配置中心（携程开源）· 按机房部署 Config Service Cluster，Portal 统一入口管理</p>
      </div>
      <span class="env-pill">● 运行中</span>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">接入机房数</div>
        <div class="value">3 / 7</div>
        <div class="delta">国内 3 已接入 · 海外 4 建设中</div>
      </div>
      <div class="stat-card">
        <div class="label">Apollo Cluster</div>
        <div class="value">7</div>
        <div class="delta">按机房 1:1 划分</div>
      </div>
      <div class="stat-card">
        <div class="label">配置项总数</div>
        <div class="value">1,260</div>
        <div class="delta">Namespace 38 个</div>
      </div>
      <div class="stat-card">
        <div class="label">接入业务线</div>
        <div class="value">2 / 5</div>
        <div class="delta">kodo · las 已接入</div>
      </div>
      <div class="stat-card">
        <div class="label">同步状态</div>
        <div class="value accent-pill">● 正常</div>
        <div class="delta">最近发布 09:12:30</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>统一入口</h3>
        <span class="meta">Apollo Portal · LDAP 单点登录</span>
      </div>
      <div class="panel-body padded">
        <div class="ext-card">
          <div class="ext-top">
            <div class="ext-logo ap-logo">AP</div>
            <div>
              <h4>Apollo Portal</h4>
              <div class="ext-sub">apollo.xinfra.internal</div>
            </div>
          </div>
          <p>统一管理 7 个机房 Config Service Cluster 的配置发布、灰度、回滚与变更审计，所有机房配置项均可在此单一入口检索与编辑，发布后由各机房 Config Service 就近下发。</p>
          <div class="sso-row">
            <span class="sso-dot"></span>LDAP 原生配置接入 · 已上线
          </div>
          <el-button type="primary" style="margin-top: 8px">打开 Apollo Portal ↗</el-button>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>机房配置中心列表</h3>
        <span class="meta">Cluster 维度 · Config Service 部署明细</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>机房 / 区域</th>
              <th>Apollo Cluster</th>
              <th>Config Service 地址</th>
              <th>承载方式</th>
              <th>接入业务线</th>
              <th>配置项数</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cluster in clusters" :key="cluster.name" class="tr-hover">
              <td><span class="tag" :class="cluster.zoneClass">{{ cluster.name }}</span></td>
              <td class="mono">cluster={{ cluster.cluster }}</td>
              <td class="mono text-xs">{{ cluster.address }}</td>
              <td class="mono">{{ cluster.type }}</td>
              <td class="mono">{{ cluster.biz }}</td>
              <td class="mono">{{ cluster.items }}</td>
              <td :class="['status-text', cluster.statusClass]">● {{ cluster.status }}</td>
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
  { name: 'YZH 机房', zoneClass: 'zone-a', cluster: 'yzh', address: 'config-yzh.xinfra.internal:8080', type: '容器化 · K8s Service', biz: 'kodo', items: 486, status: '运行中', statusClass: 'ok' },
  { name: 'XS 机房', zoneClass: 'zone-b', cluster: 'xs', address: 'config-xs.xinfra.internal:8080', type: '容器化 · K8s Service', biz: 'las', items: 398, status: '运行中', statusClass: 'ok' },
  { name: 'JF 机房', zoneClass: 'zone-c', cluster: 'jf', address: 'config-jf.xinfra.internal:8080', type: '容器化 · K8s Service', biz: '灵矽', items: 376, status: '灰度接入中', statusClass: 'warn' },
  { name: '达拉斯 IDC', zoneClass: '', cluster: 'dallas', address: '—', type: '容器化 · K8s Service', biz: '—', items: 0, status: '建设中', statusClass: 'idle' },
  { name: '新加坡 IDC', zoneClass: '', cluster: 'singapore', address: '—', type: '容器化 · K8s Service', biz: '—', items: 0, status: '建设中', statusClass: 'idle' },
  { name: '香港 IDC', zoneClass: '', cluster: 'hk', address: '—', type: '容器化 · K8s Service', biz: '—', items: 0, status: '待启动', statusClass: 'idle' },
  { name: '东南亚 IDC', zoneClass: '', cluster: 'sea', address: '—', type: '容器化 · K8s Service', biz: '—', items: 0, status: '待启动', statusClass: 'idle' },
])
</script>

<style scoped>
/* 公共样式已在 global.css 中定义 */

h4 {
  margin: 0;
  font-size: 14px;
}

p {
  font-size: 12px;
  color: var(--text-mid);
  line-height: 1.6;
  margin: 0 0 10px;
}
</style>
