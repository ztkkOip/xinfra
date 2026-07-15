<template>
  <div>
    <div class="page-head">
      <div>
        <h1>服务管理</h1>
        <p>同步全部机房 Consul 注册中心服务目录，统一查看服务名 / 服务 IP / 业务标签 / 实例数</p>
      </div>
      <el-button>⟳ 立即同步</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">接入 Consul Datacenter</div>
        <div class="value">7</div>
        <div class="delta">国内 3 · 海外 4</div>
      </div>
      <div class="stat-card">
        <div class="label">服务总数</div>
        <div class="value">186</div>
        <div class="delta">去重后唯一服务名</div>
      </div>
      <div class="stat-card">
        <div class="label">服务实例总数</div>
        <div class="value">842</div>
        <div class="delta">所有机房注册 IP 汇总</div>
      </div>
      <div class="stat-card">
        <div class="label">健康实例占比</div>
        <div class="value accent">98.6%</div>
        <div class="delta warn">12 个实例 critical</div>
      </div>
      <div class="stat-card">
        <div class="label">最近同步</div>
        <div class="value accent-pill">● 正常</div>
        <div class="delta">07:41:05 · 间隔 30s</div>
      </div>
    </div>

    <div class="las-toolbar">
      <el-select placeholder="全部机房">
        <el-option label="全部机房" value="" />
        <el-option label="YZH" value="yzh" />
        <el-option label="XS" value="xs" />
        <el-option label="JF" value="jf" />
        <el-option label="达拉斯" value="dallas" />
        <el-option label="新加坡" value="singapore" />
        <el-option label="香港" value="hk" />
        <el-option label="东南亚" value="sea" />
      </el-select>
      <el-select placeholder="全部业务标签">
        <el-option label="全部业务标签" value="" />
        <el-option label="kodo" value="kodo" />
        <el-option label="las" value="las" />
        <el-option label="lingxi" value="lingxi" />
        <el-option label="ltoken" value="ltoken" />
        <el-option label="maas" value="maas" />
      </el-select>
      <el-select placeholder="全部状态">
        <el-option label="全部状态" value="" />
        <el-option label="健康" value="healthy" />
        <el-option label="部分异常" value="partial" />
        <el-option label="全部下线" value="offline" />
      </el-select>
      <el-input placeholder="搜索服务名 / IP" class="search-input" />
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>服务台账 · 按 Consul Datacenter 同步</h3>
        <span class="meta">数据源：Consul Catalog API</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>服务名</th>
              <th>机房 / Datacenter</th>
              <th>业务标签</th>
              <th>实例数</th>
              <th>健康实例</th>
              <th>示例 IP</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="service in services" :key="service.name" class="tr-hover">
              <td class="strong mono">{{ service.name }}</td>
              <td><span class="tag" :class="service.dcClass">{{ service.dc }}</span></td>
              <td class="mono">{{ service.biz }}</td>
              <td class="mono">{{ service.instances }}</td>
              <td class="mono">{{ service.healthy }}</td>
              <td class="mono">{{ service.ip }}</td>
              <td :class="['status-text', service.statusClass]">● {{ service.status }}</td>
            </tr>
          </tbody>
        </table>
        <div class="pagination">
          <span>共 186 条 · 每页 7 条</span>
          <div class="pg-btns">
            <span class="pg-btn disabled">‹</span>
            <span class="pg-btn active">1</span>
            <span class="pg-btn">2</span>
            <span class="pg-btn">3</span>
            <span class="pg-sep">…</span>
            <span class="pg-btn">27</span>
            <span class="pg-btn">›</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const services = ref([
  { name: 'kodo-gateway-svc', dc: 'YZH', dcClass: 'zone-a', biz: 'kodo', instances: 12, healthy: '12 / 12', ip: '10.21.4.51', status: '健康', statusClass: 'ok' },
  { name: 'kodo-upload-svc', dc: 'XS', dcClass: 'zone-b', biz: 'kodo', instances: 10, healthy: '10 / 10', ip: '10.34.37.66', status: '健康', statusClass: 'ok' },
  { name: 'las-search-api', dc: 'XS', dcClass: 'zone-b', biz: 'las', instances: 18, healthy: '17 / 18', ip: '10.34.37.20', status: '部分异常', statusClass: 'warn' },
  { name: 'las-order-svc', dc: '达拉斯 IDC', dcClass: '', biz: 'las', instances: 5, healthy: '5 / 5', ip: '10.66.2.20', status: '健康', statusClass: 'ok' },
  { name: 'lingxi-render-svc', dc: 'JF', dcClass: 'zone-c', biz: 'lingxi', instances: 6, healthy: '6 / 6', ip: '10.45.2.30', status: '健康', statusClass: 'ok' },
  { name: 'ltoken-wallet-svc', dc: 'YZH', dcClass: 'zone-a', biz: 'ltoken', instances: 8, healthy: '8 / 8', ip: '10.21.4.88', status: '健康', statusClass: 'ok' },
  { name: 'maas-infer-svc', dc: '新加坡 IDC', dcClass: '', biz: 'maas', instances: 4, healthy: '3 / 4', ip: '10.88.1.40', status: '部分异常', statusClass: 'warn' },
])
</script>

<style scoped>
/* 所有样式已在 global.css 中定义 */
</style>
