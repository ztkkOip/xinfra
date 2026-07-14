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
        <div class="value" style="color: var(--accent); font-size: 16px">● 正常</div>
        <div class="delta">最近发布 09:12:30</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>统一入口</h3>
        <span class="meta">Apollo Portal · LDAP 单点登录</span>
      </div>
      <div class="panel-body" style="padding: 16px">
        <div class="ext-card">
          <div class="ext-top">
            <div class="ext-logo" style="background: #1C2A3A; color: #8EC8FF">AP</div>
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
              <td class="mono" style="font-size: 11px">{{ cluster.address }}</td>
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

.env-pill {
  font-family: var(--mono);
  font-size: 11px;
  color: var(--accent);
  border: 1px solid var(--accent-dim);
  background: #0F261D;
  padding: 3px 9px;
  border-radius: 20px;
  letter-spacing: 0.4px;
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

.ext-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 18px;
  max-width: 560px;
}

.ext-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.ext-logo {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-weight: 700;
  font-size: 13px;
}

h4 {
  margin: 0;
  font-size: 14px;
}

.ext-sub {
  font-size: 11px;
  color: var(--text-dim);
  font-family: var(--mono);
}

p {
  font-size: 12px;
  color: var(--text-mid);
  line-height: 1.6;
  margin: 0 0 10px;
}

.sso-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--text-dim);
}

.sso-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
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
  color: #8EC8FF;
  border-color: #1C3A5C;
  background: #0E1C2C;
}

.tag.zone-b {
  color: #C9A6FF;
  border-color: #3A2C5C;
  background: #1A1430;
}

.tag.zone-c {
  color: #FFC97A;
  border-color: #5C4A1C;
  background: #241B0A;
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

.status-text.idle {
  color: var(--text-dim);
}
</style>
