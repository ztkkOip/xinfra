<template>
  <div>
    <div class="page-head">
      <div>
        <h1>资源管理</h1>
        <p>统一纳管物理机与虚机资源 · 基础数据来自 SINA CMDB</p>
      </div>
      <el-button type="primary">同步</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">资源总数</div>
        <div class="value">658</div>
        <div class="delta">物理机 146 · 虚机 512</div>
      </div>
      <div class="stat-card">
        <div class="label">SINA CMDB</div>
        <div class="value">368</div>
        <div class="delta">物理机 146 · 虚机 222</div>
      </div>
      <div class="stat-card">
        <div class="label">阿里云同步</div>
        <div class="value">154</div>
        <div class="delta">ECS · 增量同步</div>
      </div>
      <div class="stat-card">
        <div class="label">AWS / 七牛 LAS 同步</div>
        <div class="value">136</div>
        <div class="delta">AWS 58 · 七牛 LAS 78</div>
      </div>
      <div class="stat-card">
        <div class="label">最近一次同步</div>
        <div class="value accent-pill">● 正常</div>
        <div class="delta">07:39:10 · 新增 3 条</div>
      </div>
    </div>

    <div class="las-toolbar">
      <el-select placeholder="全部资源类型">
        <el-option label="全部资源类型" value="" />
        <el-option label="物理机" value="physical" />
        <el-option label="虚机" value="vm" />
      </el-select>
      <el-select placeholder="全部数据来源">
        <el-option label="全部数据来源" value="" />
        <el-option label="SINA CMDB" value="cmdb" />
        <el-option label="阿里云同步" value="aliyun" />
        <el-option label="AWS 同步" value="aws" />
        <el-option label="七牛 LAS 同步" value="qiniu" />
      </el-select>
      <el-select placeholder="全部业务线">
        <el-option label="全部业务线" value="" />
        <el-option label="kodo" value="kodo" />
        <el-option label="las" value="las" />
        <el-option label="lingxi" value="lingxi" />
        <el-option label="ltoken" value="ltoken" />
        <el-option label="maas" value="maas" />
        <el-option label="未分配" value="unassigned" />
      </el-select>
      <el-input placeholder="搜索 hostname / asset_number / IP" class="search-input" />
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>资源台账 · 物理机 / 虚机</h3>
        <span class="meta">SINA CMDB asset items + 云厂商同步增量</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>主机名</th>
              <th>资产编号</th>
              <th>类型</th>
              <th>机房 / 区域</th>
              <th>内网 IP</th>
              <th>规格</th>
              <th>业务线</th>
              <th>数据来源</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in resources" :key="item.hostname" class="tr-hover">
              <td class="strong mono">{{ item.hostname }}</td>
              <td class="mono text-xs">{{ item.asset_number }}</td>
              <td><span class="tag" :class="item.type === '物理机' ? '' : 'vm'">{{ item.type }}</span></td>
              <td><span class="tag zone-a">{{ item.location }}</span></td>
              <td class="mono">{{ item.ip }}</td>
              <td class="mono text-xs">{{ item.spec }}</td>
              <td class="mono">{{ item.business_line }}</td>
              <td><span class="tag src-cmdb">{{ item.source }}</span></td>
              <td class="status-text ok">● {{ item.status }}</td>
            </tr>
          </tbody>
        </table>
        <div class="pagination">
          <span>共 658 条 · 每页 7 条</span>
          <div class="pg-btns">
            <span class="pg-btn disabled">‹</span>
            <span class="pg-btn active">1</span>
            <span class="pg-btn">2</span>
            <span class="pg-btn">3</span>
            <span class="pg-sep">…</span>
            <span class="pg-btn">94</span>
            <span class="pg-btn">›</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const resources = ref([
  { hostname: 'xs291', asset_number: 'SERV00003502', type: '物理机', location: '华东·杭州下沙', ip: '10.34.37.52', spec: '20C/192G/13.52T', business_line: 'kodo', source: 'SINA CMDB', status: 'production' },
  { hostname: 'yzh-las-014', asset_number: 'SERV00004187', type: '虚机', location: '华北·YZH', ip: '10.21.4.214', spec: '32C/128G/4.0T', business_line: 'las', source: 'SINA CMDB', status: 'production' },
  { hostname: 'jf-bm-118', asset_number: 'SERV00004290', type: '物理机', location: '华南·JF', ip: '10.45.2.18', spec: '16C/64G/2.0T', business_line: 'lingxi', source: 'SINA CMDB', status: 'production' },
  { hostname: 'ali-ecs-sz-0231', asset_number: '—', type: '虚机', location: '阿里云·华南', ip: '172.18.4.31', spec: '16C/64G/500G', business_line: 'ltoken', source: '阿里云同步', status: 'production' },
  { hostname: 'aws-us-i-0a13fe2', asset_number: '—', type: '虚机', location: 'AWS · 美国', ip: '10.66.2.12', spec: '8C/32G/200G', business_line: 'maas', source: 'AWS 同步', status: 'production' },
  { hostname: 'sg-las-007', asset_number: 'SERV00004511', type: '虚机', location: '七牛 LAS · 新加坡', ip: '10.88.1.07', spec: '16C/64G/2.0T', business_line: 'las', source: '七牛 LAS 同步', status: 'idle' },
  { hostname: 'hk-bm-001', asset_number: 'SERV00004602', type: '物理机', location: '香港 IDC', ip: '10.90.0.11', spec: '8C/32G/1.0T', business_line: '未分配', source: 'SINA CMDB', status: 'production' },
])
</script>

<style scoped>
/* 公共样式已在 global.css 中定义 */

.tag.vm {
  color: var(--tag-blue-text);
  border-color: var(--tag-blue-border);
  background: var(--tag-blue-bg);
}

.tag.src-cmdb {
  color: var(--tag-green-text);
  border-color: var(--tag-green-border);
  background: var(--tag-green-bg);
}
</style>
