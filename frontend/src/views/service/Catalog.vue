<template>
  <div>
    <!-- 基础服务目录 -->
    <div class="page-head">
      <div>
        <h1>基础服务目录</h1>
        <p>点击卡片，通过封装好的 ansible-playbook 一键标准化交付，完成后自动注册进对应子系统</p>
      </div>
    </div>

    <div class="svc-grid">
      <div
        v-for="service in basicServices"
        :key="service.name"
        class="svc-card"
        :class="{ disabled: service.disabled }"
        :style="service.disabled ? { opacity: service.opacity } : {}"
        @click="handleCardClick(service)"
      >
        <div class="svc-top">
          <div class="svc-ic" :style="{ background: service.bgColor, color: service.iconColor }">
            {{ service.icon }}
          </div>
          <h4 :style="service.disabled && service.icon === '+' ? { color: 'var(--text-dim)' } : {}">{{ service.name }}</h4>
        </div>
        <div class="svc-desc">{{ service.description }}</div>
        <div class="svc-foot">
          <span class="playbook">{{ service.playbook }}</span>
          <span :style="service.statusColor ? { color: service.statusColor } : {}">{{ service.version }}</span>
        </div>
      </div>
    </div>

    <!-- 容器服务部署 -->
    <div class="page-head" style="margin-top: 24px;">
      <div>
        <h1>容器服务部署</h1>
        <p>容器化服务统一管理与发布</p>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <h3>统一入口</h3>
        <span class="meta">Wayne · LDAP 单点登录</span>
      </div>
      <div class="panel-body" style="padding: 16px;">
        <div class="ext-card" style="max-width: 560px;">
          <div class="ext-top">
            <div class="ext-logo" style="background: var(--logo-w-bg); color: var(--tag-blue-text);">W</div>
            <div>
              <h4>Wayne</h4>
              <div class="ext-sub">wayne.xinfra.internal</div>
            </div>
          </div>
          <p>业务容器发布、命名空间与配额管理，复用 Wayne 原生多租户能力。所有集群的容器发布均可通过此入口统一操作。</p>
          <div class="sso-row">
            <span class="sso-dot"></span> LDAP 原生配置接入 · 在线
          </div>
          <button class="btn btn-primary" style="margin-top: 8px; align-self: flex-start;" @click="openWayne">
            打开 Wayne ↗
          </button>
        </div>
      </div>
    </div>

    <!-- 弹窗：数据库标准化交付 · MySQL -->
    <Teleport to="body">
      <div v-if="showMysqlModal" class="overlay" @click.self="showMysqlModal = false">
        <div class="modal">
          <div class="modal-head">
            <h3>数据库标准化交付 · MySQL</h3>
            <button class="close" @click="showMysqlModal = false">&times;</button>
          </div>
          <div class="modal-body">
            <div class="field-row">
              <div class="field">
                <label>业务线</label>
                <select v-model="mysqlForm.bl">
                  <option value="kodo">kodo</option>
                  <option value="linxi">linxi</option>
                  <option value="xinfra">xinfra</option>
                  <option value="las">las</option>
                </select>
              </div>
              <div class="field">
                <label>架构模式</label>
                <select v-model="mysqlForm.topology">
                  <option value="1m2r">一主两从</option>
                  <option value="1m1r">一主一从</option>
                </select>
              </div>
            </div>
            <div class="field-row">
              <div class="field">
                <label>规格</label>
                <select v-model="mysqlForm.spec">
                  <option value="8C32G">8C 32G · 500G SSD（Ceph RBD）</option>
                  <option value="16C64G">16C 64G · 1T SSD</option>
                </select>
              </div>
              <div class="field">
                <label>版本</label>
                <select v-model="mysqlForm.version">
                  <option value="8.0.36">MySQL 8.0.36</option>
                  <option value="5.7.44">MySQL 5.7.44</option>
                </select>
              </div>
            </div>
            <div class="field">
              <label>实例名称</label>
              <input v-model="mysqlForm.instanceName" placeholder="mysql-las-billing-02" />
            </div>
            <div class="yaml-preview"><span class="c"># ansible-runner 调用预览</span>
<span class="k">playbook</span>: roles/mysql-deploy
<span class="k">tenant</span>: <span class="s">{{ mysqlForm.bl }}</span>
<span class="k">topology</span>: <span class="s">{{ mysqlForm.topology === '1m2r' ? '1-master-2-replica' : '1-master-1-replica' }}</span>
<span class="k">node_selector</span>: <span class="s">business-line={{ mysqlForm.bl }}</span>
<span class="k">storage_class</span>: <span class="s">ceph-rbd-ssd</span>
<span class="k">register_to</span>: <span class="s">clouddm</span>  <span class="c"># 自动注册数据源，跳 CloudDM 上线</span></div>
          </div>
          <div class="modal-foot">
            <button class="btn btn-ghost" @click="showMysqlModal = false">取消</button>
            <button class="btn btn-primary" @click="handleMysqlSubmit">确认交付</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { subsystemApi } from '@/api/subsystem'
import { useBusinessLineMockProfile } from '@/utils/businessLineMock'

const router = useRouter()
const { currentName } = useBusinessLineMockProfile()

interface Service {
  name: string
  icon: string
  bgColor: string
  iconColor: string
  description: string
  playbook: string
  version: string
  statusColor?: string
  disabled?: boolean
  opacity?: number
  modal?: string
}

const basicServices = ref<Service[]>([
  {
    name: 'MySQL',
    icon: 'My',
    bgColor: '#1C2E4A',
    iconColor: '#7FB8FF',
    description: '一主多从标准化交付，完成后自动注册进 CloudDM / open-cdm 数据源，统一审核 SQL 上线。',
    playbook: 'roles/mysql-deploy',
    version: 'v8.0.36',
    modal: 'mysql',
  },
  {
    name: 'openresty',
    icon: 'Or',
    bgColor: '#1C3A2E',
    iconColor: '#7FFFC2',
    description: '容器化部署在 K8s 内，作为业务边缘网关，支持灰度路由配置。',
    playbook: 'roles/openresty-deploy',
    version: 'v1.25',
    modal: 'openresty',
  },
  {
    name: 'PgSQL',
    icon: 'Pg',
    bgColor: '#2E1C3A',
    iconColor: '#C9A6FF',
    description: '流复制主从架构，规划中 · 接入 CloudDM 统一审核体系。',
    playbook: 'roles/pgsql-deploy',
    version: '规划中',
    statusColor: 'var(--text-dim)',
    disabled: true,
    opacity: 0.8,
  },
  {
    name: 'Redis',
    icon: 'Rd',
    bgColor: '#3A1C28',
    iconColor: '#FF8E9C',
    description: 'Redis 管理由 CacheCloud 承载，规划接入后开放标准化交付。',
    playbook: 'CacheCloud',
    version: '规划接入',
    disabled: true,
    opacity: 0.6,
  },
  {
    name: '接入新服务',
    icon: '+',
    bgColor: '#262C38',
    iconColor: 'var(--text-dim)',
    description: '封装新的 ansible-playbook，注册进基础服务目录',
    playbook: '',
    version: '',
    disabled: true,
    opacity: 0.5,
  },
])

// 弹窗状态
const showMysqlModal = ref(false)

const mysqlForm = reactive({
  bl: currentName.value,
  topology: '1m2r',
  spec: '8C32G',
  version: '8.0.36',
  instanceName: `mysql-${currentName.value}-billing-02`,
})

watch(currentName, (name) => {
  mysqlForm.bl = name
  mysqlForm.instanceName = `mysql-${name}-billing-02`
})

const handleCardClick = (service: Service) => {
  if (service.disabled) {
    if (service.name === '接入新服务') return
    ElMessage.warning(`${service.name} 功能规划中，敬请期待`)
    return
  }
  if (service.modal === 'wayne') {
    openWayne()
  } else if (service.modal === 'mysql') {
    showMysqlModal.value = true
  } else if (service.modal === 'openresty') {
    ElMessage.info('openresty 部署功能对接中，敬请期待')
  }
}

const openWayne = async () => {
  try {
    const wayne = subsystemApi.getSubsystemByName('wayne')
    if (!wayne) {
      ElMessage.error('Wayne 子系统未找到')
      return
    }
    const { data } = await subsystemApi.getSSOUrl(wayne.id)
    window.location.assign(data.sso_url)
  } catch {
    // 错误已由 request 拦截器统一处理
  }
}

const handleMysqlSubmit = () => {
  showMysqlModal.value = false
  ElMessage.success('MySQL 标准化交付已提交，完成后自动注册 CloudDM')
  router.push({ name: 'TaskLog' })
}
</script>

<style scoped>
.svc-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.svc-card:hover:not(.disabled) {
  border-color: #3A4356;
}

.svc-card.disabled {
  cursor: not-allowed;
}

.svc-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.svc-ic {
  width: 34px;
  height: 34px;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-weight: 600;
  font-size: 13px;
}

h4 {
  margin: 0;
  font-size: 13.5px;
}

.svc-desc {
  font-size: 12px;
  color: var(--text-dim);
  line-height: 1.5;
  margin-bottom: 12px;
  min-height: 36px;
}

.svc-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: var(--text-dim);
}

.playbook {
  font-family: var(--mono);
  color: var(--text-mid);
}

/* ========== overlay / modal ========== */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(7, 9, 13, 0.72);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  z-index: 50;
  padding-top: 6vh;
}

.modal {
  width: 560px;
  max-width: 92vw;
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.5);
}

.modal-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 18px;
  border-bottom: 1px solid var(--line-soft, var(--line));
}

.modal-head h3 {
  margin: 0;
  font-size: 14.5px;
}

.modal-head .close {
  margin-left: auto;
  color: var(--text-dim);
  font-size: 18px;
  border: none;
  background: none;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}

.modal-body {
  padding: 18px;
}

.modal-foot {
  padding: 14px 18px;
  border-top: 1px solid var(--line-soft, var(--line));
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* ========== form fields ========== */
.field {
  margin-bottom: 14px;
}

.field label {
  display: block;
  font-size: 11.5px;
  color: var(--text-dim);
  margin-bottom: 6px;
}

.field select,
.field input {
  width: 100%;
  background: var(--bg-panel-2, var(--bg-panel));
  border: 1px solid var(--line);
  border-radius: 6px;
  padding: 8px 10px;
  color: var(--text-hi);
  font-size: 12.5px;
  font-family: var(--sans, inherit);
  outline: none;
}

.field select:focus,
.field input:focus {
  border-color: var(--accent-dim, var(--accent));
}

.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

/* ========== yaml preview ========== */
.yaml-preview {
  background: #0B0D12;
  border: 1px solid var(--line-soft, var(--line));
  border-radius: 6px;
  padding: 12px 14px;
  font-family: var(--mono);
  font-size: 11.5px;
  color: #8EC8FF;
  line-height: 1.7;
  white-space: pre;
  overflow-x: auto;
}

.yaml-preview .k {
  color: #7FB8FF;
}

.yaml-preview .s {
  color: #FFC97A;
}

.yaml-preview .c {
  color: var(--text-dim);
}

/* ========== tooltip note ========== */
.tooltip-note {
  background: #0F261D;
  border: 1px solid var(--accent-dim, var(--accent));
  border-radius: 8px;
  padding: 12px 14px;
  font-size: 12px;
  color: var(--text-mid);
  line-height: 1.6;
  margin-top: 4px;
}

.tooltip-note b {
  color: var(--accent);
}

/* ========== buttons ========== */
.btn {
  border-radius: 6px;
  border: 1px solid var(--line);
  background: var(--bg-panel-2, var(--bg-panel));
  color: var(--text-hi);
  padding: 8px 14px;
  font-size: 12.5px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
}

.btn:hover {
  border-color: #3A4356;
}

.btn-primary {
  background: var(--accent);
  border-color: var(--accent);
  color: #06150F;
  font-weight: 600;
}

.btn-primary:hover {
  background: #55E8B3;
}

.btn-ghost {
  background: transparent;
  color: var(--text-mid);
  border-color: transparent;
}

.btn-ghost:hover {
  color: var(--text-hi);
  background: var(--bg-panel-2, var(--bg-panel));
}
</style>
