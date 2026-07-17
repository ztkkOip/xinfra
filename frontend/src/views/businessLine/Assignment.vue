<template>
  <section class="page">
    <div class="page-head">
      <div>
        <h2>业务线分配</h2>
        <p>{{ currentName }}</p>
      </div>
    </div>

    <div v-if="isCurrentBusinessLineAdmin" class="assignment-panel">
      <el-form label-position="top">
        <el-form-item label="用户">
          <el-select v-model="grantForm.target_user_id" filterable placeholder="选择用户">
            <el-option v-for="user in users" :key="user.uid" :label="user.username" :value="user.uid" />
          </el-select>
        </el-form-item>
        <el-form-item label="业务线">
          <el-select v-model="grantForm.target_business_line_id" filterable placeholder="选择业务线">
            <el-option v-for="item in businessLineStore.items" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限">
          <el-segmented v-model="grantForm.permission" :options="permissionOptions" />
        </el-form-item>
        <el-button type="primary" :loading="granting" @click="grantPermission">保存分配</el-button>
      </el-form>

      <div class="section-divider"></div>

      <el-form label-position="top">
        <el-form-item label="Wayne namespace">
          <el-select
            v-model="selectedWayneNamespaceIds"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            :loading="loadingWayneNamespaces"
            placeholder="选择 Wayne namespace"
          >
            <el-option
              v-for="item in wayneNamespaces"
              :key="item.id"
              :label="`${item.name} / ${item.kubeNamespace}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-button type="primary" :loading="savingWayneNamespaces" @click="saveWayneNamespaceMapping">保存 Wayne namespace 映射</el-button>
      </el-form>
    </div>
    <div v-else class="empty-state">需要当前业务线管理员权限</div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { businessLineApi, type WayneNamespace } from '@/api/businessLine'
import { userApi, type UserOption } from '@/api/user'
import { useBusinessLineStore } from '@/stores/businessLine'

const businessLineStore = useBusinessLineStore()
const isCurrentBusinessLineAdmin = computed(() => businessLineStore.isCurrentAdmin)
const currentName = computed(() => businessLineStore.current?.name || '未选择业务线')
const users = ref<UserOption[]>([])
const wayneNamespaces = ref<WayneNamespace[]>([])
const selectedWayneNamespaceIds = ref<number[]>([])
const granting = ref(false)
const loadingWayneNamespaces = ref(false)
const savingWayneNamespaces = ref(false)
const grantForm = reactive<{
  target_user_id: number | null
  target_business_line_id: number | null
  permission: 0 | 1
}>({
  target_user_id: null,
  target_business_line_id: null,
  permission: 1,
})

const permissionOptions = [
  { label: '管理员', value: 0 },
  { label: '普通用户', value: 1 },
]

watch(
  () => businessLineStore.current?.id,
  (id) => {
    if (id && !grantForm.target_business_line_id) {
      grantForm.target_business_line_id = id
    }
    if (id && isCurrentBusinessLineAdmin.value) {
      loadWayneNamespaceMapping(id)
    }
  },
  { immediate: true },
)

watch(
  isCurrentBusinessLineAdmin,
  async (admin) => {
    if (admin && !users.value.length) {
      users.value = await userApi.list()
    }
    if (admin) {
      await loadWayneNamespaces()
      const businessLineID = businessLineStore.current?.id
      if (businessLineID) {
        await loadWayneNamespaceMapping(businessLineID)
      }
    }
  },
  { immediate: true },
)

onMounted(async () => {
  if (isCurrentBusinessLineAdmin.value && !users.value.length) {
    users.value = await userApi.list()
  }
  if (isCurrentBusinessLineAdmin.value) {
    await loadWayneNamespaces()
    const businessLineID = businessLineStore.current?.id
    if (businessLineID) {
      await loadWayneNamespaceMapping(businessLineID)
    }
  }
})

async function grantPermission() {
  const businessLineID = businessLineStore.current?.id
  if (!businessLineID || !grantForm.target_user_id || !grantForm.target_business_line_id) {
    ElMessage.warning('请选择用户和业务线')
    return
  }
  granting.value = true
  try {
    await businessLineApi.grant({
      business_line_id: businessLineID,
      target_user_id: grantForm.target_user_id,
      target_business_line_id: grantForm.target_business_line_id,
      permission: grantForm.permission,
    })
    ElMessage.success('已保存业务线分配')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存业务线分配失败')
  } finally {
    granting.value = false
  }
}

async function loadWayneNamespaces() {
  if (wayneNamespaces.value.length) {
    return
  }
  loadingWayneNamespaces.value = true
  try {
    wayneNamespaces.value = await businessLineApi.listWayneNamespaces()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '查询 Wayne namespace 失败')
  } finally {
    loadingWayneNamespaces.value = false
  }
}

async function loadWayneNamespaceMapping(businessLineID: number) {
  loadingWayneNamespaces.value = true
  try {
    const mapped = await businessLineApi.listMappedWayneNamespaces(businessLineID)
    selectedWayneNamespaceIds.value = mapped.map((item) => item.id)
  } catch (error) {
    selectedWayneNamespaceIds.value = []
    ElMessage.error(error instanceof Error ? error.message : '查询 Wayne namespace 映射失败')
  } finally {
    loadingWayneNamespaces.value = false
  }
}

async function saveWayneNamespaceMapping() {
  const businessLineID = businessLineStore.current?.id
  if (!businessLineID) {
    ElMessage.warning('请选择当前业务线')
    return
  }
  const selected = wayneNamespaces.value.filter((item) => selectedWayneNamespaceIds.value.includes(item.id))
  savingWayneNamespaces.value = true
  try {
    await businessLineApi.replaceMappedWayneNamespaces(businessLineID, selected)
    ElMessage.success('已保存 Wayne namespace 映射')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存 Wayne namespace 映射失败')
  } finally {
    savingWayneNamespaces.value = false
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

.assignment-panel {
  max-width: 520px;
  margin-top: 18px;
}

.section-divider {
  border-top: 1px solid var(--line);
  margin: 22px 0 18px;
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
