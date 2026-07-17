<template>
  <section class="page">
    <div class="page-head">
      <div>
        <h2>业务线管理</h2>
        <p>平台管理员维护全局业务线</p>
      </div>
    </div>

    <template v-if="isPlatformAdmin">
      <div class="toolbar">
        <el-input v-model="newName" class="name-input" placeholder="业务线名称" @keyup.enter="createBusinessLine" />
        <el-button type="primary" :icon="Plus" :loading="savingLine" @click="createBusinessLine">新增</el-button>
        <el-button :loading="loadingLines" @click="loadAllBusinessLines">刷新</el-button>
      </div>

      <el-table :data="allBusinessLines" v-loading="loadingLines" class="data-table" empty-text="没有业务线数据">
        <el-table-column prop="id" label="ID" width="90" />
        <el-table-column label="名称" min-width="220">
          <template #default="{ row }">
            <el-input v-if="editingId === row.id" v-model="editingName" size="small" @keyup.enter="saveBusinessLine(row.id)" />
            <span v-else>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" min-width="210" />
        <el-table-column label="操作" width="190" fixed="right">
          <template #default="{ row }">
            <template v-if="editingId === row.id">
              <el-button size="small" type="primary" @click="saveBusinessLine(row.id)">保存</el-button>
              <el-button size="small" @click="cancelEdit">取消</el-button>
            </template>
            <template v-else>
              <el-button size="small" :icon="Edit" @click="startEdit(row)">编辑</el-button>
              <el-button size="small" type="danger" :icon="Delete" @click="deleteBusinessLine(row)">删除</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </template>
    <div v-else class="empty-state">需要平台管理员权限</div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { businessLineApi, type BusinessLine } from '@/api/businessLine'
import { useAuthStore } from '@/stores/auth'
import { useBusinessLineStore } from '@/stores/businessLine'

const authStore = useAuthStore()
const businessLineStore = useBusinessLineStore()
const isPlatformAdmin = computed(() => authStore.user?.is_admin === true)

const allBusinessLines = ref<BusinessLine[]>([])
const loadingLines = ref(false)
const savingLine = ref(false)
const newName = ref('')
const editingId = ref<number | null>(null)
const editingName = ref('')

watch(
  isPlatformAdmin,
  async (admin) => {
    if (admin) {
      await loadAllBusinessLines()
    }
  },
  { immediate: true },
)

onMounted(async () => {
  try {
    await authStore.refreshUser()
  } catch {
    return
  }
  if (isPlatformAdmin.value) {
    await loadAllBusinessLines()
  }
})

async function loadAllBusinessLines() {
  loadingLines.value = true
  try {
    allBusinessLines.value = await businessLineApi.listAll()
  } catch (error) {
    allBusinessLines.value = []
    ElMessage.error(error instanceof Error ? error.message : '查询业务线失败')
  } finally {
    loadingLines.value = false
  }
}

async function createBusinessLine() {
  const name = newName.value.trim()
  if (!name) {
    ElMessage.warning('请输入业务线名称')
    return
  }
  savingLine.value = true
  try {
    await businessLineApi.create(name)
    newName.value = ''
    await loadAllBusinessLines()
    await businessLineStore.loadMine()
    ElMessage.success('已新增业务线')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '新增业务线失败')
  } finally {
    savingLine.value = false
  }
}

function startEdit(row: BusinessLine) {
  editingId.value = row.id
  editingName.value = row.name
}

function cancelEdit() {
  editingId.value = null
  editingName.value = ''
}

async function saveBusinessLine(id: number) {
  const name = editingName.value.trim()
  if (!name) {
    ElMessage.warning('请输入业务线名称')
    return
  }
  try {
    await businessLineApi.update(id, name)
    cancelEdit()
    await loadAllBusinessLines()
    await businessLineStore.loadMine()
    ElMessage.success('已更新业务线')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '更新业务线失败')
  }
}

async function deleteBusinessLine(row: BusinessLine) {
  try {
    await ElMessageBox.confirm(`确认删除业务线 ${row.name}？`, '删除业务线', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await businessLineApi.remove(row.id)
    await loadAllBusinessLines()
    await businessLineStore.loadMine()
    ElMessage.success('已删除业务线')
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : '删除业务线失败')
    }
  }
}
</script>

<style scoped>
.page {
  min-height: 100%;
}

.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--line);
  padding-bottom: 18px;
}

.toolbar {
  display: flex;
  gap: 10px;
  margin: 18px 0 14px;
}

.name-input {
  width: 260px;
}

.data-table {
  width: 100%;
}

.empty-state {
  padding: 28px 0;
  color: var(--text-dim);
}

h2 {
  margin: 0 0 6px;
  font-size: 20px;
}

p {
  margin: 0;
  color: var(--text-dim);
  font-size: 14px;
}
</style>
