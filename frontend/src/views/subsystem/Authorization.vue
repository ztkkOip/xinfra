<template>
  <div>
    <div class="page-head">
      <div>
        <h1>子系统赋权</h1>
        <p>{{ currentBusinessLineName }} · Wayne namespace 角色授权</p>
      </div>
      <div class="head-actions">
        <el-button :loading="loading" @click="reloadAll">刷新</el-button>
        <el-button v-if="hasWayneRoleBindingPermission" type="primary" :disabled="!canOperate || !selectedUserId || !selectedNamespaceId" :loading="saving" @click="saveRoles">
          保存授权
        </el-button>
      </div>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">接入子系统</div>
        <div class="value">{{ enabledSystemCount }}</div>
        <div class="delta">{{ systemSummary }}</div>
      </div>
      <div class="stat-card">
        <div class="label">Wayne Namespace</div>
        <div class="value">{{ wayneNamespaces.length }}</div>
        <div class="delta">当前业务线映射</div>
      </div>
      <div class="stat-card">
        <div class="label">可选角色</div>
        <div class="value">{{ wayneRoles.length }}</div>
        <div class="delta">{{ roleSummary }}</div>
      </div>
      <div class="stat-card">
        <div class="label">当前操作权限</div>
        <div class="value state-value" :class="{ ok: canOperate, warn: !canOperate }">● {{ operatorStateText }}</div>
        <div class="delta">{{ operatorStateDetail }}</div>
      </div>
    </div>

    <div class="matrix">
      <div class="system-card">
        <div class="system-top">
          <div class="logo wayne">W</div>
          <div>
            <h3>Wayne</h3>
            <p>业务线 namespace 角色绑定，默认新用户初始化为访客</p>
          </div>
        </div>
        <div class="role-grid">
          <div v-for="role in wayneRoles" :key="role.id" class="role-cell">
            <span>{{ role.name }}</span>
            <strong>#{{ role.id }}</strong>
          </div>
          <div v-if="!wayneRoles.length" class="role-cell empty-cell">暂无角色</div>
        </div>
      </div>

      <div class="system-card muted-card">
        <div class="system-top">
          <div class="logo clouddm">DM</div>
          <div>
            <h3>CloudDM</h3>
            <p>接口预留，当前不开放授权操作</p>
          </div>
        </div>
        <div class="role-grid">
          <div class="role-cell">
            <span>状态</span>
            <strong>未启用</strong>
          </div>
        </div>
      </div>
    </div>

    <div v-if="hasWayneRoleBindingPermission" class="panel auth-panel">
      <div class="panel-head">
        <h3>Wayne 授权操作</h3>
        <span class="meta">数据源：AuthServer · Wayne internal API</span>
      </div>
      <div class="form-grid">
        <el-form label-position="top">
          <el-form-item label="用户">
            <el-select v-model="selectedUserId" filterable placeholder="选择用户" :loading="loadingUsers" @change="loadSelectedUserRoles">
              <el-option v-for="user in users" :key="user.uid" :label="user.username" :value="user.uid" />
            </el-select>
          </el-form-item>
          <el-form-item label="Wayne Namespace">
            <el-select v-model="selectedNamespaceId" filterable placeholder="选择 namespace" :loading="loadingNamespaces" @change="syncSelectedRoleIds">
              <el-option
                v-for="namespace in wayneNamespaces"
                :key="namespace.id"
                :disabled="!namespace.can_bind && !namespace.can_unbind"
                :label="namespaceLabel(namespace)"
                :value="namespace.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Wayne 角色">
            <el-select v-model="selectedRoleIds" multiple filterable collapse-tags collapse-tags-tooltip placeholder="选择角色">
              <el-option v-for="role in wayneRoles" :key="role.id" :label="role.name" :value="role.id" />
            </el-select>
          </el-form-item>
          <div class="button-row">
            <el-button type="primary" :disabled="!canOperate || !selectedUserId || !selectedNamespaceId || !selectedRoleIds.length" :loading="saving" @click="saveRoles">
              保存角色
            </el-button>
            <el-button :disabled="!canUnbindSelected || !selectedUserId || !selectedNamespaceId" :loading="clearing" @click="clearRoles">
              清空角色
            </el-button>
            <el-button :disabled="!canOperate || !selectedUserId" :loading="initializing" @click="initVisitor">
              初始化访客
            </el-button>
          </div>
        </el-form>

        <div class="hint-panel">
          <div class="hint-title">授权规则</div>
          <p>当前账号必须是平台管理员或当前业务线管理员。</p>
          <p>保存前会再次校验 Wayne namespace 的授权能力。</p>
          <p>用户加入业务线时后端会自动初始化 Wayne 访客角色。</p>
        </div>
      </div>
    </div>

    <div v-else class="panel readonly-panel">
      <div class="panel-head">
        <h3>当前权限</h3>
        <span class="meta">只读模式</span>
      </div>
      <div class="readonly-body">
        当前账号没有 Wayne namespace 角色绑定权限，只展示现有权限。
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>当前用户 Wayne 角色</h3>
        <span class="meta">{{ selectedUsername || '未选择用户' }}</span>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr>
              <th>Namespace</th>
              <th>Kube Namespace</th>
              <th>当前角色</th>
              <th>授权能力</th>
              <th v-if="hasWayneRoleBindingPermission">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="namespace in wayneNamespaces" :key="namespace.id" class="tr-hover">
              <td>
                <div class="strong">{{ namespace.name || '-' }}</div>
                <div class="sub mono">id={{ namespace.id }}</div>
              </td>
              <td class="mono">{{ namespace.kubeNamespace || '-' }}</td>
              <td>
                <span v-if="roleNamesForNamespace(namespace.id).length" class="role-list">
                  <span v-for="role in roleNamesForNamespace(namespace.id)" :key="role" class="role">{{ role }}</span>
                </span>
                <span v-else class="scope">未绑定</span>
              </td>
              <td>
                <span v-if="namespace.permission_error" class="status-text warn">● {{ namespace.permission_error }}</span>
                <span v-else-if="namespace.can_bind" class="status-text ok">● 可授权</span>
                <span v-else class="status-text idle">● 无授权权限</span>
              </td>
              <td v-if="hasWayneRoleBindingPermission">
                <div class="actions">
                  <button type="button" @click="chooseNamespace(namespace.id)">选择</button>
                  <button type="button" class="danger" :disabled="!namespace.can_unbind || !selectedUserId" @click="clearNamespaceRoles(namespace.id)">清空</button>
                </div>
              </td>
            </tr>
            <tr v-if="!wayneNamespaces.length">
              <td :colspan="hasWayneRoleBindingPermission ? 5 : 4" class="empty-row">当前业务线没有绑定 Wayne namespace</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { subsystemAuthApi, type SubsystemAuthSystem, type WayneBusinessLineNamespace, type WayneRole, type WayneUserRoles } from '@/api/subsystemAuth'
import { userApi, type UserOption } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import { useBusinessLineStore } from '@/stores/businessLine'

const authStore = useAuthStore()
const businessLineStore = useBusinessLineStore()

const systems = ref<SubsystemAuthSystem[]>([])
const users = ref<UserOption[]>([])
const wayneRoles = ref<WayneRole[]>([])
const wayneNamespaces = ref<WayneBusinessLineNamespace[]>([])
const userRoles = ref<WayneUserRoles | null>(null)
const selectedUserId = ref<number | null>(null)
const selectedNamespaceId = ref<number | null>(null)
const selectedRoleIds = ref<number[]>([])
const loading = ref(false)
const loadingUsers = ref(false)
const loadingNamespaces = ref(false)
const saving = ref(false)
const clearing = ref(false)
const initializing = ref(false)

const isPlatformAdmin = computed(() => authStore.isAdmin)
const canManageBusinessLine = computed(() => isPlatformAdmin.value || businessLineStore.isCurrentAdmin)
const currentBusinessLineId = computed(() => businessLineStore.current?.id || null)
const currentBusinessLineName = computed(() => businessLineStore.current?.name || '未选择业务线')
const selectedUser = computed(() => users.value.find((user) => user.uid === selectedUserId.value))
const selectedUsername = computed(() => selectedUser.value?.username || '')
const selectedNamespace = computed(() => wayneNamespaces.value.find((item) => item.id === selectedNamespaceId.value))
const canOperate = computed(() => canManageBusinessLine.value && Boolean(selectedNamespace.value?.can_bind))
const canUnbindSelected = computed(() => canManageBusinessLine.value && Boolean(selectedNamespace.value?.can_unbind))
const hasWayneRoleBindingPermission = computed(() =>
  canManageBusinessLine.value && wayneNamespaces.value.some((item) => item.can_bind || item.can_unbind),
)
const enabledSystemCount = computed(() => systems.value.filter((item) => item.enabled).length)
const systemSummary = computed(() => systems.value.map((item) => `${item.name}${item.enabled ? '' : '(未启用)'}`).join(' · ') || 'Wayne')
const roleSummary = computed(() => wayneRoles.value.map((item) => item.name).join(' · ') || '暂无')
const operatorStateText = computed(() => {
	if (!canManageBusinessLine.value) return '无业务线权限'
	if (!wayneNamespaces.value.length) return '未绑定 namespace'
	if (hasWayneRoleBindingPermission.value) return '可授权'
	return '只读'
})
const operatorStateDetail = computed(() => {
	if (!canManageBusinessLine.value) return '需要平台管理员或当前业务线管理员'
	if (!wayneNamespaces.value.length) return '先在业务线分配中绑定 Wayne namespace'
	if (!hasWayneRoleBindingPermission.value) return '没有 Wayne 角色绑定权限'
	return `${wayneNamespaces.value.filter((item) => item.can_bind).length} 个 namespace 可授权`
})

watch(
  () => businessLineStore.current?.id,
  async () => {
    selectedUserId.value = null
    selectedNamespaceId.value = null
    selectedRoleIds.value = []
    userRoles.value = null
    await Promise.all([loadUsers(), loadBusinessLineData()])
  },
)

watch(selectedNamespaceId, () => {
  syncSelectedRoleIds()
})

onMounted(async () => {
  await reloadAll()
})

async function reloadAll() {
  loading.value = true
  try {
    await ensureBusinessLinesLoaded()
    await Promise.all([loadSystems(), loadUsers(), loadWayneRoles(), loadBusinessLineData()])
    if (selectedUserId.value) {
      await loadSelectedUserRoles()
    }
  } finally {
    loading.value = false
  }
}

async function ensureBusinessLinesLoaded() {
  if (businessLineStore.current?.id && businessLineStore.businessLines.length) {
    return
  }
  await businessLineStore.loadMine().catch(() => {})
}

async function loadSystems() {
  try {
    systems.value = await subsystemAuthApi.listSystems()
  } catch {
    systems.value = [
      { key: 'wayne', name: 'Wayne', enabled: true },
      { key: 'clouddm', name: 'CloudDM', enabled: false },
    ]
  }
}

async function loadUsers() {
  const businessLineId = currentBusinessLineId.value
  if (!businessLineId) {
    users.value = []
    selectedUserId.value = null
    return
  }
  loadingUsers.value = true
  try {
    users.value = await userApi.list({ businessLineId })
    if (!selectedUserId.value && users.value.length) {
      selectedUserId.value = users.value[0].uid
      await loadSelectedUserRoles()
    } else if (selectedUserId.value && !users.value.some((user) => user.uid === selectedUserId.value)) {
      selectedUserId.value = users.value[0]?.uid || null
      await loadSelectedUserRoles()
    }
  } catch (error) {
    users.value = []
    selectedUserId.value = null
    ElMessage.error(error instanceof Error ? error.message : '查询用户失败')
  } finally {
    loadingUsers.value = false
  }
}

async function loadWayneRoles() {
  try {
    wayneRoles.value = await subsystemAuthApi.listWayneRoles()
  } catch (error) {
    wayneRoles.value = []
    ElMessage.error(error instanceof Error ? error.message : '查询 Wayne 角色失败')
  }
}

async function loadBusinessLineData() {
  const businessLineId = currentBusinessLineId.value
  if (!businessLineId || !canManageBusinessLine.value) {
    wayneNamespaces.value = []
    return
  }
  loadingNamespaces.value = true
  try {
    wayneNamespaces.value = await subsystemAuthApi.listWayneNamespaces(businessLineId)
    if (!selectedNamespaceId.value && wayneNamespaces.value.length) {
      selectedNamespaceId.value = wayneNamespaces.value[0].id
    }
  } catch (error) {
    wayneNamespaces.value = []
    ElMessage.error(error instanceof Error ? error.message : '查询 Wayne namespace 失败')
  } finally {
    loadingNamespaces.value = false
  }
}

async function loadSelectedUserRoles() {
  if (!selectedUsername.value) {
    userRoles.value = null
    selectedRoleIds.value = []
    return
  }
  try {
    userRoles.value = await subsystemAuthApi.getWayneUserRoles(selectedUsername.value)
    syncSelectedRoleIds()
  } catch (error) {
    userRoles.value = null
    selectedRoleIds.value = []
    ElMessage.error(error instanceof Error ? error.message : '查询用户 Wayne 角色失败')
  }
}

function syncSelectedRoleIds() {
  const namespaceId = selectedNamespaceId.value
  if (!namespaceId || !userRoles.value?.namespaces) {
    selectedRoleIds.value = []
    return
  }
  const binding = userRoles.value.namespaces.find((item) => Number(item.namespace?.id) === namespaceId)
  selectedRoleIds.value = (binding?.groups || []).map((item) => Number(item.id)).filter(Boolean)
}

async function saveRoles() {
  const businessLineId = currentBusinessLineId.value
  const namespaceId = selectedNamespaceId.value
  const username = selectedUsername.value
  if (!businessLineId || !namespaceId || !username || !selectedRoleIds.value.length) {
    ElMessage.warning('请选择用户、namespace 和角色')
    return
  }
  saving.value = true
  try {
    await subsystemAuthApi.bindWayneNamespaceRoles(businessLineId, namespaceId, username, {
      groupIds: selectedRoleIds.value,
      replace: true,
      requestId: requestId('wayne-bind'),
      reason: '子系统赋权',
    })
    ElMessage.success('Wayne 角色已保存')
    await loadSelectedUserRoles()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存 Wayne 角色失败')
  } finally {
    saving.value = false
  }
}

async function clearRoles() {
  const namespaceId = selectedNamespaceId.value
  if (!namespaceId) return
  await clearNamespaceRoles(namespaceId)
}

async function clearNamespaceRoles(namespaceId: number) {
  const businessLineId = currentBusinessLineId.value
  const username = selectedUsername.value
  if (!businessLineId || !username) {
    ElMessage.warning('请选择用户')
    return
  }
  await ElMessageBox.confirm('确认清空该用户在此 Wayne namespace 下的角色？', '清空角色', {
    type: 'warning',
    confirmButtonText: '清空',
    cancelButtonText: '取消',
  })
  clearing.value = true
  try {
    await subsystemAuthApi.unbindWayneNamespaceRoles(businessLineId, namespaceId, username, {
      requestId: requestId('wayne-clear'),
      reason: '子系统赋权清空角色',
    })
    ElMessage.success('Wayne 角色已清空')
    await loadSelectedUserRoles()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '清空 Wayne 角色失败')
  } finally {
    clearing.value = false
  }
}

async function initVisitor() {
  const businessLineId = currentBusinessLineId.value
  const userId = selectedUserId.value
  if (!businessLineId || !userId) {
    ElMessage.warning('请选择用户')
    return
  }
  initializing.value = true
  try {
    await subsystemAuthApi.initWayneBusinessLineUser(businessLineId, userId)
    ElMessage.success('已初始化 Wayne 访客角色')
    await loadSelectedUserRoles()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '初始化 Wayne 访客角色失败')
  } finally {
    initializing.value = false
  }
}

function chooseNamespace(namespaceId: number) {
  selectedNamespaceId.value = namespaceId
}

function roleNamesForNamespace(namespaceId: number): string[] {
  const binding = userRoles.value?.namespaces?.find((item) => Number(item.namespace?.id) === namespaceId)
  return (binding?.groups || []).map((item) => item.name).filter(Boolean)
}

function namespaceLabel(namespace: WayneBusinessLineNamespace) {
  const ability = namespace.can_bind ? '可授权' : namespace.can_unbind ? '可清空' : '无权限'
  return `${namespace.name || namespace.id} / ${namespace.kubeNamespace || '-'} · ${ability}`
}

function requestId(prefix: string) {
  return `${prefix}-${Date.now()}`
}
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

.head-actions,
.button-row {
  display: flex;
  gap: 8px;
  align-items: center;
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

.state-value {
  font-size: 16px !important;
}

.state-value.ok {
  color: var(--ok);
}

.state-value.warn {
  color: var(--warn);
}

.stat-card .delta {
  font-size: 11px;
  color: var(--text-dim);
  margin-top: 4px;
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

.muted-card {
  opacity: 0.76;
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
  background: #1c2a3a;
  color: #7fb8ff;
}

.logo.clouddm {
  background: #1c3a2e;
  color: #7fffc2;
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

.empty-cell {
  justify-content: center;
  color: var(--text-dim);
}

.panel {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  margin-bottom: 16px;
}

.auth-panel {
  padding-bottom: 4px;
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

.form-grid {
  display: grid;
  grid-template-columns: minmax(360px, 520px) 1fr;
  gap: 20px;
  padding: 16px;
}

.hint-panel {
  border: 1px solid var(--line-soft);
  background: var(--bg-panel-2);
  border-radius: 8px;
  padding: 14px;
  color: var(--text-dim);
  font-size: 12px;
}

.hint-title {
  color: var(--text-hi);
  font-weight: 600;
  margin-bottom: 8px;
}

.hint-panel p {
  margin: 6px 0;
}

.readonly-panel {
  margin-bottom: 16px;
}

.readonly-body {
  padding: 16px;
  color: var(--text-dim);
  font-size: 12.5px;
}

.panel-body {
  padding: 4px 0;
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 920px;
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
  background: #171d28;
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

.role-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.role {
  color: var(--text-hi);
  background: var(--bg-panel-2);
  border: 1px solid var(--line-soft);
  border-radius: 5px;
  padding: 3px 7px;
}

.scope {
  color: var(--text-dim);
  font-size: 11px;
}

.status-text.ok {
  color: var(--ok);
}

.status-text.warn {
  color: var(--warn);
}

.status-text.idle {
  color: var(--text-dim);
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

.actions button:hover:not(:disabled) {
  color: var(--text-hi);
  border-color: #3a4356;
}

.actions button:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.actions button.danger {
  color: #ff9a95;
}

.empty-row {
  text-align: center;
  color: var(--text-dim);
  padding: 26px 16px;
}

@media (max-width: 1180px) {
  .stat-row,
  .matrix,
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
