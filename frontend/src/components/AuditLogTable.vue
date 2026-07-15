<template>
  <div class="panel">
    <div class="panel-head">
      <h3>{{ title }}</h3>
      <span class="meta">共 {{ total }} 条</span>
    </div>
    <div class="panel-body">
      <table>
        <thead>
          <tr>
            <th v-for="col in columns" :key="col.key">{{ col.label }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in data" :key="row.id" class="tr-hover">
            <td v-for="col in columns" :key="col.key" :class="col.class">
              <template v-if="col.key === 'status'">
                <span :class="['status-text', row.status]">
                  ● {{ statusMap[row.status] || row.status }}
                </span>
              </template>
              <template v-else-if="col.key === 'source_ip' || col.key === 'username'">
                <span class="mono">{{ row[col.key] }}</span>
              </template>
              <template v-else>
                {{ row[col.key] }}
              </template>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="pagination">
        <span>共 {{ total }} 条 · 每页 {{ pageSize }} 条</span>
        <div class="pg-btns">
          <span class="pg-btn disabled">‹</span>
          <span class="pg-btn active">1</span>
          <span class="pg-btn">2</span>
          <span class="pg-btn">3</span>
          <span class="pg-sep">…</span>
          <span class="pg-btn">{{ Math.ceil(total / pageSize) }}</span>
          <span class="pg-btn">›</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Column {
  key: string
  label: string
  class?: string
}

withDefaults(defineProps<{
  title: string
  columns: Column[]
  data: any[]
  total: number
  pageSize?: number
}>(), {
  pageSize: 10,
})

const statusMap: Record<string, string> = {
  success: '成功',
  failed: '失败',
  running: '执行中',
}
</script>

<style scoped>
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

.mono {
  font-family: var(--mono);
}

.status-text {
  font-size: 12px;
  display: flex;
  align-items: center;
}

.status-text.success {
  color: var(--accent);
}

.status-text.failed {
  color: var(--err);
}

.status-text.running {
  color: var(--warn);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-top: 1px solid var(--line-soft);
  font-size: 11.5px;
  color: var(--text-dim);
}

.pg-btns {
  display: flex;
  gap: 4px;
  align-items: center;
}

.pg-btn {
  min-width: 26px;
  height: 26px;
  padding: 0 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--line);
  border-radius: 5px;
  color: var(--text-mid);
  font-family: var(--mono);
  font-size: 11px;
  background: var(--bg-panel-2);
  cursor: pointer;
}

.pg-btn.active {
  background: var(--accent-dim);
  border-color: var(--accent-dim);
  color: #CFFCE9;
}

.pg-btn:hover {
  border-color: #3A4356;
  color: var(--text-hi);
}

.pg-btn.disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.pg-sep {
  color: var(--text-dim);
  padding: 0 2px;
}
</style>
