<template>
  <div>
    <div class="page-head">
      <div>
        <h1>子系统赋权</h1>
        <p>Wayne / CloudDM 入口权限、默认角色与授权状态</p>
      </div>
      <el-button type="primary">新增授权</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">接入子系统</div>
        <div class="value">2</div>
        <div class="delta">Wayne · CloudDM</div>
      </div>
      <div class="stat-card">
        <div class="label">授权主体</div>
        <div class="value">6</div>
        <div class="delta">用户 3 · 用户组 3</div>
      </div>
      <div class="stat-card">
        <div class="label">待审批</div>
        <div class="value" style="color: var(--warn)">2</div>
        <div class="delta">最近提交 10:18</div>
      </div>
      <div class="stat-card">
        <div class="label">默认授权</div>
        <div class="value" style="font-size: 16px; color: var(--accent)">● 生效</div>
        <div class="delta">新用户默认只读</div>
      </div>
    </div>

    <div class="toolbar">
      <el-select v-model="filters.system" placeholder="全部子系统" style="width: 150px">
        <el-option label="全部子系统" value="" />
        <el-option label="Wayne" value="Wayne" />
        <el-option label="CloudDM" value="CloudDM" />
      </el-select>
      <el-select v-model="filters.type" placeholder="全部主体" style="width: 150px">
        <el-option label="全部主体" value="" />
        <el-option label="用户" value="user" />
        <el-option label="用户组" value="group" />
      </el-select>
      <el-select v-model="filters.status" placeholder="全部状态" style="width: 150px">
        <el-option label="全部状态" value="" />
        <el-option label="已生效" value="active" />
        <el-option label="待审批" value="pending" />
        <el-option label="已停用" value="disabled" />
      </el-select>
      <el-input v-model="filters.keyword" placeholder="搜索账号 / 用户组 / 角色" style="flex: 1" />
    </div>

    <div class="matrix">
      <div v-for="system in systems" :key="system.name" class="system-card">
        <div class="system-top">
          <div class="logo" :class="system.className">{{ system.icon }}</div>
          <div>
            <h3>{{ system.name }}</h3>
            <p>{{ system.defaultPolicy }}</p>
          </div>
        </div>
        <div class="role-grid">
          <div v-for="role in system.roles" :key="role.name" class="role-cell">
            <span>{{ role.name }}</span>
            <strong>{{ role.count }}</strong>
          </div>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>授权列表</h3>
        <span class="meta">数据源：AuthServer · 子系统授权</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>授权主体</th>
              <th>类型</th>
              <th>子系统</th>
              <th>角色 / 范围</th>
              <th>来源</th>
              <th>状态</th>
              <th>最近变更</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredGrants" :key="item.id" class="tr-hover">
              <td>
                <div class="principal">
                  <span class="avatar">{{ item.initial }}</span>
                  <div>
                    <div class="strong">{{ item.principal }}</div>
                    <div class="sub mono">{{ item.detail }}</div>
                  </div>
                </div>
              </td>
              <td><span class="tag">{{ item.type === 'user' ? '用户' : '用户组' }}</span></td>
              <td class="mono">{{ item.system }}</td>
              <td>
                <span class="role">{{ item.role }}</span>
                <span class="scope mono">{{ item.scope }}</span>
              </td>
              <td class="mono">{{ item.source }}</td>
              <td :class="['status-text', item.statusClass]">● {{ item.statusText }}</td>
              <td class="mono">{{ item.updatedAt }}</td>
              <td>
                <div class="actions">
                  <button type="button">编辑</button>
                  <button type="button" class="danger">停用</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'

const filters = reactive({
  system: '',
  type: '',
  status: '',
  keyword: '',
})

const systems = [
  {
    name: 'Wayne',
    icon: 'W',
    className: 'wayne',
    defaultPolicy: '默认 namespace 只读 · DemoGroupId=23',
    roles: [
      { name: '只读', count: 18 },
      { name: '发布', count: 6 },
      { name: '管理员', count: 2 },
    ],
  },
  {
    name: 'CloudDM',
    icon: 'DM',
    className: 'clouddm',
    defaultPolicy: 'OIDC 登录 · SQL 审核角色映射',
    roles: [
      { name: '查询', count: 21 },
      { name: '审核', count: 5 },
      { name: '管理员', count: 1 },
    ],
  },
]

const grants = [
  { id: 1, principal: 'eastsales@qiniu.com', initial: 'E', detail: 'eastsales', type: 'user', system: 'Wayne', role: '默认只读', scope: 'namespace=demo', source: 'SSO 自动初始化', status: 'active', statusText: '已生效', statusClass: 'ok', updatedAt: '2026-07-15 10:12' },
  { id: 2, principal: 'platform-admin', initial: 'P', detail: 'LDAP group', type: 'group', system: 'Wayne', role: '管理员', scope: 'all namespaces', source: '手动授权', status: 'active', statusText: '已生效', statusClass: 'ok', updatedAt: '2026-07-14 18:40' },
  { id: 3, principal: 'dba-reviewers', initial: 'D', detail: 'LDAP group', type: 'group', system: 'CloudDM', role: 'SQL 审核', scope: 'prod / staging', source: '手动授权', status: 'active', statusText: '已生效', statusClass: 'ok', updatedAt: '2026-07-14 16:05' },
  { id: 4, principal: 'las-dev', initial: 'L', detail: 'LDAP group', type: 'group', system: 'CloudDM', role: '查询', scope: 'las schemas', source: '审批流', status: 'pending', statusText: '待审批', statusClass: 'warn', updatedAt: '2026-07-15 10:18' },
  { id: 5, principal: 'ops-user1@qiniu.com', initial: 'O', detail: 'ops-user1', type: 'user', system: 'Wayne', role: '发布', scope: 'namespace=demo', source: '审批流', status: 'pending', statusText: '待审批', statusClass: 'warn', updatedAt: '2026-07-15 09:55' },
  { id: 6, principal: 'temp-sql@qiniu.com', initial: 'T', detail: 'temp-sql', type: 'user', system: 'CloudDM', role: '查询', scope: 'expired', source: '临时授权', status: 'disabled', statusText: '已停用', statusClass: 'idle', updatedAt: '2026-07-13 20:30' },
]

const filteredGrants = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()
  return grants.filter((item) => {
    if (filters.system && item.system !== filters.system) return false
    if (filters.type && item.type !== filters.type) return false
    if (filters.status && item.status !== filters.status) return false
    if (!keyword) return true
    return [item.principal, item.detail, item.system, item.role, item.scope]
      .join(' ')
      .toLowerCase()
      .includes(keyword)
  })
})
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

.stat-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 18px;
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

.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
  align-items: center;
}

.matrix {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 16px;
}

.system-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 16px;
}

.system-top {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 14px;
}

.logo {
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

.logo.wayne {
  background: #1C2A3A;
  color: #7FB8FF;
}

.logo.clouddm {
  background: #1C3A2E;
  color: #7FFFC2;
}

.system-top h3 {
  margin: 0 0 4px;
  font-size: 14px;
}

.system-top p {
  margin: 0;
  color: var(--text-dim);
  font-size: 11.5px;
}

.role-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.role-cell {
  border: 1px solid var(--line-soft);
  background: var(--bg-panel-2);
  border-radius: 6px;
  padding: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text-mid);
}

.role-cell strong {
  font-family: var(--mono);
  color: var(--text-hi);
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
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 980px;
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

.tr-hover:hover {
  background: #171D28;
}

.strong {
  color: var(--text-hi);
  font-weight: 600;
}

.sub {
  color: var(--text-dim);
  font-size: 11px;
  margin-top: 2px;
}

.principal {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #2A3142;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  color: var(--text-hi);
  font-size: 11px;
}

.role {
  color: var(--text-hi);
  margin-right: 8px;
}

.scope {
  color: var(--text-dim);
  font-size: 11px;
}

.actions {
  display: flex;
  gap: 8px;
}

.actions button {
  border: 1px solid var(--line);
  background: var(--bg-panel-2);
  color: var(--text-mid);
  border-radius: 5px;
  padding: 4px 8px;
  font-size: 12px;
}

.actions button:hover {
  color: var(--text-hi);
  border-color: #3A4356;
}

.actions button.danger {
  color: #FF9A95;
}
</style>
