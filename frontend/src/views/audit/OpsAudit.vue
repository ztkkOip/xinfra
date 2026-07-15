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
  try {
    const { data } = await auditApi.getOpsAudit()
    auditData.value = data.items
    total.value = data.total
  } catch {
    // 错误已由 request 拦截器统一处理
  }
}

const refresh = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 所有样式已在 global.css 中定义 */
</style>
