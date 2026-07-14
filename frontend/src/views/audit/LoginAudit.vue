<template>
  <div>
    <div class="page-head">
      <div>
        <h1>登录审计</h1>
        <p>查询所有用户的登录记录，包括时间、IP、目标系统、结果</p>
      </div>
      <el-button @click="refresh">⟳ 刷新</el-button>
    </div>
    <AuditLogTable
      title="登录审计记录"
      :columns="columns"
      :data="auditData"
      :total="total"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { auditApi, type LoginAudit } from '@/api/audit'
import AuditLogTable from '@/components/AuditLogTable.vue'

const columns = [
  { key: 'username', label: '用户名', class: 'mono' },
  { key: 'login_time', label: '登录时间' },
  { key: 'source_ip', label: '来源 IP', class: 'mono' },
  { key: 'target_system', label: '目标系统' },
  { key: 'status', label: '状态' },
]

const auditData = ref<LoginAudit[]>([])
const total = ref(0)

const fetchData = async () => {
  const { data } = await auditApi.getLoginAudit()
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
