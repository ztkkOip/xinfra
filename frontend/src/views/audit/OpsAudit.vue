<template>
  <div>
    <div class="page-head">
      <div>
        <h1>运维操作审计</h1>
        <p>查询运维操作记录，包括操作类型、目标、结果</p>
      </div>
      <el-button @click="refresh">⟳ 刷新</el-button>
    </div>
    <AuditLogTable
      title="运维操作审计记录"
      :columns="columns"
      :data="auditData"
      :total="total"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { auditApi, type OpsAudit } from '@/api/audit'
import AuditLogTable from '@/components/AuditLogTable.vue'

const columns = [
  { key: 'username', label: '用户名', class: 'mono' },
  { key: 'operation_type', label: '操作类型' },
  { key: 'operation', label: '操作内容' },
  { key: 'target', label: '目标' },
  { key: 'status', label: '状态' },
  { key: 'created_at', label: '时间' },
]

const auditData = ref<OpsAudit[]>([])
const total = ref(0)

const fetchData = async () => {
  const { data } = await auditApi.getOpsAudit()
  auditData.value = data.items
  total.value = data.total
}

const refresh = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
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
  letter-spacing: 0.2px;
}

.page-head p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12.5px;
}
</style>
