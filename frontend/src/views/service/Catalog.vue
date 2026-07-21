<template>
  <div>
    <template v-if="!activeServiceKey">
    <div class="page-head">
      <div>
        <h1>基础服务目录</h1>
        <p>点击基础组件卡片进入对应部署配置页面，通过模板完成标准化交付</p>
      </div>
    </div>

    <div class="svc-grid">
      <div
        v-for="service in basicServices"
        :key="service.key"
        class="svc-card"
        :class="{ disabled: service.disabled }"
        :style="service.disabled ? { opacity: service.opacity } : {}"
        @click="handleCardClick(service)"
      >
        <div class="svc-top">
          <div class="svc-ic" :style="{ background: service.bgColor, color: service.iconColor }">
            {{ service.icon }}
          </div>
          <div>
            <h4>{{ service.name }}</h4>
            <span class="svc-status">{{ service.status }}</span>
          </div>
        </div>
        <div class="svc-desc">{{ service.description }}</div>
        <div class="svc-foot">
          <span class="playbook">{{ service.playbook }}</span>
          <span :style="service.statusColor ? { color: service.statusColor } : {}">{{ service.version }}</span>
        </div>
      </div>
    </div>
    </template>

    <template v-else-if="activeService">
    <div class="page-head">
      <div>
        <h1>{{ activeService.name }} 标准化交付</h1>
        <p>任务编号 {{ taskNo }} · 当前业务线：{{ currentName }}</p>
      </div>
      <div class="head-actions page-actions">
        <el-button :icon="Back" @click="backToCatalog">返回目录</el-button>
        <el-button :icon="Refresh" @click="resetWorkbench">重置工作台</el-button>
      </div>
    </div>

    <div class="panel delivery-panel">
      <div class="panel-head workbench-head">
        <div>
          <h3>{{ activeService.name }} 标准化交付</h3>
          <span class="meta">任务编号：{{ taskNo }} · 当前业务线：{{ currentName }}</span>
        </div>
        <div class="head-actions">
          <span class="tag" :class="deliveryStateClass">{{ deliveryStateText }}</span>
        </div>
      </div>

      <div class="panel-body padded">
        <div class="flow-tabs">
          <button
            v-for="tab in flowTabs"
            :key="tab.key"
            type="button"
            class="flow-tab"
            :class="{ active: activeView === tab.key }"
            @click="activeView = tab.key"
          >
            <strong>{{ tab.title }}</strong>
            <span>{{ tab.subtitle }}</span>
          </button>
        </div>

        <section v-if="activeView === 'config'" class="workbench-grid">
          <div class="config-panel">
            <div class="section-head">
              <div>
                <h3>{{ activeService.name }} 交付参数</h3>
                <p>{{ activeService.template }} · {{ activeService.runner }}</p>
              </div>
              <span class="tag tag-green">标准模板</span>
            </div>

            <div class="form-grid">
              <label class="form-field">
                实例名称
                <el-input v-model="deliveryForm.instanceName" />
              </label>
              <label class="form-field">
                组件版本
                <el-select v-model="deliveryForm.version">
                  <el-option v-for="version in activeService.versions" :key="version" :label="version" :value="version" />
                </el-select>
              </label>
            </div>

            <div class="config-list">
              <div class="config-row">
                <div class="config-label">实例架构</div>
                <div class="choice-row">
                  <button
                    v-for="mode in activeService.modes"
                    :key="mode.value"
                    type="button"
                    class="choice-btn"
                    :class="{ active: deliveryForm.mode === mode.value }"
                    @click="deliveryForm.mode = mode.value"
                  >
                    {{ mode.label }}
                  </button>
                </div>
              </div>
              <div class="config-row">
                <div class="config-label">安装路径</div>
                <div class="subform-grid">
                  <label class="form-field">
                    程序目录
                    <el-input v-model="deliveryForm.installPath" />
                  </label>
                  <label class="form-field">
                    数据目录
                    <el-input v-model="deliveryForm.dataPath" />
                  </label>
                  <label class="form-field">
                    日志目录
                    <el-input v-model="deliveryForm.logPath" />
                  </label>
                </div>
              </div>
              <div class="config-row">
                <div class="config-label">运行规格</div>
                <div class="choice-row">
                  <button
                    v-for="spec in activeService.specs"
                    :key="spec"
                    type="button"
                    class="choice-btn"
                    :class="{ active: deliveryForm.spec === spec }"
                    @click="deliveryForm.spec = spec"
                  >
                    {{ spec }}
                  </button>
                </div>
              </div>
              <div class="config-row">
                <div class="config-label">数据盘容量</div>
                <div class="choice-row">
                  <button
                    v-for="disk in activeService.disks"
                    :key="disk"
                    type="button"
                    class="choice-btn"
                    :class="{ active: deliveryForm.disk === disk }"
                    @click="deliveryForm.disk = disk"
                  >
                    {{ disk }}
                  </button>
                </div>
              </div>
              <div class="config-row">
                <div class="config-label">服务端口</div>
                <div class="subform-grid">
                  <label class="form-field">
                    监听端口
                    <el-input-number v-model="deliveryForm.port" :min="1024" :max="65535" controls-position="right" />
                  </label>
                  <label class="form-field">
                    字符集 / 协议
                    <el-select v-model="deliveryForm.charset">
                      <el-option v-for="item in activeService.charsets" :key="item" :label="item" :value="item" />
                    </el-select>
                  </label>
                  <label class="form-field">
                    资源池
                    <el-select v-model="deliveryForm.pool">
                      <el-option label="LAS-TEST-DB-Pool · 自动匹配" value="auto" />
                      <el-option label="由管理员手动指定" value="manual" />
                    </el-select>
                  </label>
                </div>
              </div>
            </div>
          </div>

          <div class="preview-panel">
            <div class="section-head">
              <div>
                <h3>架构预览</h3>
                <p>{{ topologySummary }}</p>
              </div>
            </div>

            <div class="topology" :class="`nodes-${topologyNodes.length}`">
              <div v-for="node in topologyNodes" :key="node.name" class="node-card">
                <div class="node-head">
                  <span class="node-icon">{{ activeService.icon }}</span>
                  <strong>{{ node.name }}</strong>
                </div>
                <div class="node-meta">
                  <span>{{ node.ip }}</span>
                  <span>{{ deliveryForm.spec }}</span>
                  <span>{{ node.role }}</span>
                  <span>可用 {{ node.freeDisk }}</span>
                </div>
              </div>
            </div>

            <pre class="yaml-preview"><code>{{ runnerPreview }}</code></pre>

            <div class="form-actions">
              <el-button :icon="CircleCheck" @click="precheck">执行预检查</el-button>
              <el-button type="primary" :icon="Promotion" :disabled="!precheckPassed" @click="createTask">创建交付任务</el-button>
            </div>
          </div>
        </section>

        <section v-else-if="activeView === 'execution'" class="execution-view">
          <div class="execution-head">
            <div>
              <h3>Ansible 自动交付</h3>
              <p>{{ deploymentId || '创建任务后生成 deployment_id' }} · {{ activeService.runner }} · {{ activeService.template }}</p>
            </div>
            <div class="head-actions">
              <el-button :disabled="!canCancelDeployment" :loading="canceling" @click="cancelDeployment">取消任务</el-button>
            </div>
          </div>

          <div class="step-list">
            <div v-for="step in steps" :key="step.name" class="delivery-step" :class="step.state">
              <div class="step-dot">{{ stepSymbol(step.state) }}</div>
              <div>
                <strong>{{ step.name }}</strong>
                <p>{{ step.desc }}</p>
              </div>
              <span class="tag">{{ stepStatusText(step.state) }}</span>
            </div>
          </div>

          <pre class="log-preview"><code>{{ deliveryLog }}</code></pre>
        </section>

        <section v-else class="result-view">
          <div class="result-head" :class="{ failed: deliveryFailed }">
            <div class="result-icon">{{ deliveryFailed ? '!' : '✓' }}</div>
            <div>
              <h3>{{ resultTitle }}</h3>
              <p>{{ resultSubtitle }}</p>
            </div>
          </div>

          <div class="result-grid">
            <div class="result-item">
              <span>访问地址</span>
              <strong>{{ resultAddress }}</strong>
            </div>
            <div class="result-item">
              <span>配置基线</span>
              <strong>{{ activeService.configId }}</strong>
            </div>
            <div class="result-item">
              <span>资产编号</span>
              <strong>{{ activeService.assetId }}</strong>
            </div>
            <div class="result-item">
              <span>健康状态</span>
              <strong>{{ deliveryFailed ? '健康检查失败 · 已回退' : activeService.healthText }}</strong>
            </div>
          </div>

          <div class="asset-table">
            <table>
              <thead>
                <tr>
                  <th>实例</th>
                  <th>组件</th>
                  <th>版本</th>
                  <th>架构</th>
                  <th>地址</th>
                  <th>状态</th>
                  <th>注册系统</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td class="mono">{{ deliveryForm.instanceName }}</td>
                  <td>{{ activeService.name }}</td>
                  <td>{{ deliveryForm.version }}</td>
                  <td>{{ currentModeLabel }}</td>
                  <td><code>{{ resultAddress }}</code></td>
                  <td><span class="tag" :class="deliveryFailed ? 'tag-red' : 'tag-green'">{{ deliveryFailed ? '已回退' : '正常' }}</span></td>
                  <td><code>{{ activeService.registerTo }}</code></td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </div>
    </template>
    <div v-else class="panel">
      <div class="panel-body padded empty-delivery">
        <h3>未找到对应基础服务</h3>
        <p>请返回基础服务目录选择可交付组件。</p>
        <el-button type="primary" :icon="Back" @click="backToCatalog">返回目录</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Back, CircleCheck, Promotion, Refresh } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { deploymentApi } from '@/api/deployment'
import { useBusinessLineStore } from '@/stores/businessLine'
import { useBusinessLineMockProfile } from '@/utils/businessLineMock'

type FlowView = 'config' | 'execution' | 'result'
type StepState = 'pending' | 'running' | 'done' | 'failed'

interface ServiceMode {
  value: string
  label: string
}

interface Service {
  key: string
  name: string
  icon: string
  bgColor: string
  iconColor: string
  description: string
  playbook: string
  version: string
  status: string
  template: string
  runner: string
  versions: string[]
  modes: ServiceMode[]
  specs: string[]
  disks: string[]
  charsets: string[]
  defaultPort: number
  defaultPaths: {
    install: string
    data: string
    log: string
  }
  registerTo: string
  configId: string
  assetId: string
  healthText: string
  statusColor?: string
  disabled?: boolean
  opacity?: number
}

interface DeliveryStep {
  name: string
  desc: string
  state: StepState
}

const { currentName } = useBusinessLineMockProfile()
const route = useRoute()
const router = useRouter()
const businessLineStore = useBusinessLineStore()

const basicServices = ref<Service[]>([
  {
    key: 'mysql',
    name: 'MySQL',
    icon: 'My',
    bgColor: 'var(--logo-w-bg)',
    iconColor: 'var(--tag-blue-text)',
    description: '一主多从标准化交付，完成后自动注册进 CloudDM / open-cdm 数据源，统一审核 SQL 上线。',
    playbook: 'roles/mysql-deploy',
    version: 'v8.0.36',
    status: '可交付',
    template: 'mysql-delivery@v1.8.3',
    runner: 'runner-02',
    versions: ['MySQL 8.0.36', 'MySQL 8.0.40', 'MySQL 5.7.44'],
    modes: [
      { value: 'single', label: '单实例' },
      { value: 'replica', label: '一主一从' },
      { value: 'mgr', label: 'MGR 三节点' },
    ],
    specs: ['4C / 16G', '8C / 32G', '16C / 64G'],
    disks: ['200 GB', '500 GB', '1 TB', '2 TB'],
    charsets: ['utf8mb4', 'utf8'],
    defaultPort: 3306,
    defaultPaths: { install: '/opt/mysql', data: '/data/mysql', log: '/data/mysql-log' },
    registerTo: 'CloudDM / CMDB / Prometheus',
    configId: 'CFG-MY-10872',
    assetId: 'CDM-MY-10872',
    healthText: '数据库正常 · 延迟 0 秒',
  },
  {
    key: 'openresty',
    name: 'openresty',
    icon: 'Or',
    bgColor: 'var(--logo-dm-bg)',
    iconColor: 'var(--tag-green-text)',
    description: '容器化部署在 K8s 内，作为业务边缘网关，支持灰度路由配置。',
    playbook: 'roles/openresty-deploy',
    version: 'v1.25',
    status: '可交付',
    template: 'nginx-delivery@v1.3.0',
    runner: 'runner-02',
    versions: ['OpenResty 1.25.3', 'OpenResty 1.21.4'],
    modes: [
      { value: 'single', label: '单节点' },
      { value: 'replica', label: '双节点 + VIP' },
    ],
    specs: ['2C / 4G', '4C / 8G', '8C / 16G'],
    disks: ['50 GB', '100 GB', '200 GB'],
    charsets: ['HTTP / HTTPS', 'HTTP only'],
    defaultPort: 443,
    defaultPaths: { install: '/opt/openresty', data: '/data/openresty', log: '/data/openresty-log' },
    registerTo: 'CMDB / Prometheus / 日志平台',
    configId: 'CFG-NGX-02841',
    assetId: 'CMP-NGX-1042',
    healthText: 'HTTP 200 · 23 ms',
  },
  {
    key: 'pgsql',
    name: 'PgSQL',
    icon: 'Pg',
    bgColor: 'var(--logo-ap-bg)',
    iconColor: 'var(--tag-purple-text)',
    description: '流复制主从架构，规划中 · 接入 CloudDM 统一审核体系。',
    playbook: 'roles/pgsql-deploy',
    version: '规划中',
    status: '规划中',
    template: '',
    runner: '',
    versions: [],
    modes: [],
    specs: [],
    disks: [],
    charsets: [],
    defaultPort: 5432,
    defaultPaths: { install: '', data: '', log: '' },
    registerTo: '',
    configId: '',
    assetId: '',
    healthText: '',
    statusColor: 'var(--text-dim)',
    disabled: true,
    opacity: 0.75,
  },
  {
    key: 'redis',
    name: 'Redis',
    icon: 'Rd',
    bgColor: 'var(--logo-cc-bg)',
    iconColor: 'var(--err)',
    description: 'Redis 管理由 CacheCloud 承载，规划接入后开放标准化交付。',
    playbook: 'CacheCloud',
    version: '规划接入',
    status: '规划接入',
    template: '',
    runner: '',
    versions: [],
    modes: [],
    specs: [],
    disks: [],
    charsets: [],
    defaultPort: 6379,
    defaultPaths: { install: '', data: '', log: '' },
    registerTo: '',
    configId: '',
    assetId: '',
    healthText: '',
    disabled: true,
    opacity: 0.6,
  },
  {
    key: 'new',
    name: '接入新服务',
    icon: '+',
    bgColor: 'var(--bg-panel-2)',
    iconColor: 'var(--text-dim)',
    description: '封装新的 ansible-playbook，注册进基础服务目录。',
    playbook: '',
    version: '',
    status: '待接入',
    template: '',
    runner: '',
    versions: [],
    modes: [],
    specs: [],
    disks: [],
    charsets: [],
    defaultPort: 0,
    defaultPaths: { install: '', data: '', log: '' },
    registerTo: '',
    configId: '',
    assetId: '',
    healthText: '',
    disabled: true,
    opacity: 0.5,
  },
])

const flowTabs = [
  { key: 'config' as FlowView, title: '01 参数配置', subtitle: '模板与资源规划' },
  { key: 'execution' as FlowView, title: '02 自动交付', subtitle: 'Ansible + 健康探测' },
  { key: 'result' as FlowView, title: '03 交付结果', subtitle: '资产与访问地址' },
]

const activeServiceKey = ref(componentParam())
const activeView = ref<FlowView>('config')
const precheckPassed = ref(false)
const running = ref(false)
const deliveryDone = ref(false)
const deliveryFailed = ref(false)
const canceling = ref(false)
const deploymentId = ref('')
const eventSource = ref<EventSource | null>(null)
const deliveryLog = ref('[ready] 等待创建交付任务...')

const deliveryForm = reactive({
  instanceName: `mysql-${currentName.value}-billing-02`,
  version: 'MySQL 8.0.36',
  mode: 'replica',
  spec: '8C / 32G',
  disk: '500 GB',
  port: 3306,
  charset: 'utf8mb4',
  pool: 'auto',
  installPath: '/opt/mysql',
  dataPath: '/data/mysql',
  logPath: '/data/mysql-log',
})

const steps = ref<DeliveryStep[]>([])

const activeService = computed(() => basicServices.value.find((service) => service.key === activeServiceKey.value && !service.disabled))
const taskNo = computed(() => `CMP-20260721-${activeServiceKey.value === 'mysql' ? '0024' : '0023'}`)
const currentModeLabel = computed(() => activeService.value?.modes.find((mode) => mode.value === deliveryForm.mode)?.label || '-')
const topologySummary = computed(() => `${currentModeLabel.value} · ${deliveryForm.spec} · ${deliveryForm.disk}`)
const topologyNodes = computed(() => {
  if (deliveryForm.mode === 'single') {
    return [{ name: activeServiceKey.value === 'mysql' ? 'db-tx-021' : 'ngx-tx-011', ip: activeServiceKey.value === 'mysql' ? '10.24.18.21' : '10.24.31.11', role: '主节点', freeDisk: '820 GB' }]
  }
  if (deliveryForm.mode === 'mgr') {
    return [
      { name: 'db-tx-021', ip: '10.24.18.21', role: 'MGR 成员', freeDisk: '820 GB' },
      { name: 'db-tx-022', ip: '10.24.18.22', role: 'MGR 成员', freeDisk: '790 GB' },
      { name: 'db-tx-023', ip: '10.24.18.23', role: 'MGR 成员', freeDisk: '820 GB' },
    ]
  }
  return activeServiceKey.value === 'mysql'
    ? [
        { name: 'db-tx-021', ip: '10.24.18.21', role: '主节点', freeDisk: '820 GB' },
        { name: 'db-tx-022', ip: '10.24.18.22', role: '从节点', freeDisk: '790 GB' },
      ]
    : [
        { name: 'ngx-tx-011', ip: '10.24.31.11', role: '主节点', freeDisk: '180 GB' },
        { name: 'ngx-tx-012', ip: '10.24.31.12', role: '备用节点', freeDisk: '180 GB' },
      ]
})
const runnerPreview = computed(() => {
  const playbook = activeService.value?.playbook || '-'
  return [
    '# ansible-runner 调用预览',
    `playbook: ${playbook}`,
    `tenant: ${currentName.value}`,
    `component: ${activeService.value?.key || '-'}`,
    `instance: ${deliveryForm.instanceName}`,
    `topology: ${currentModeLabel.value}`,
    `node_selector: business-line=${currentName.value}`,
    `install_path: ${deliveryForm.installPath}`,
    `data_path: ${deliveryForm.dataPath}`,
    `register_to: ${activeService.value?.registerTo || '-'}`,
  ].join('\n')
})
const resultAddress = computed(() => `${topologyNodes.value[0]?.ip || '10.24.18.21'}:${deliveryForm.port}`)
const resultTitle = computed(() => {
  if (deliveryFailed.value) return '交付失败，已自动回退'
  return activeServiceKey.value === 'mysql' ? 'MySQL 实例已交付' : 'OpenResty 集群已交付'
})
const resultSubtitle = computed(() => {
  if (deliveryFailed.value) return '健康检查未通过 · 资源已释放 · 变更未交付'
  return activeServiceKey.value === 'mysql' ? '全部步骤执行成功 · 用时 06:42' : '全部步骤执行成功 · 用时 02:18'
})
const deliveryStateText = computed(() => {
  if (deliveryFailed.value) return '已回退'
  if (deliveryDone.value) return '已交付'
  if (running.value) return '交付中'
  if (precheckPassed.value) return '待执行'
  return '配置中'
})
const canCancelDeployment = computed(() => Boolean(deploymentId.value) && running.value && !deliveryDone.value && !deliveryFailed.value)
const deliveryStateClass = computed(() => {
  if (deliveryFailed.value) return 'tag-red'
  if (deliveryDone.value) return 'tag-green'
  if (precheckPassed.value || running.value) return 'tag-amber'
  return ''
})

watch(currentName, (name) => {
  deliveryForm.instanceName = `${activeServiceKey.value === 'mysql' ? 'mysql' : 'nginx'}-${name}-billing-02`
})

watch(activeServiceKey, () => {
  hydrateServiceDefaults()
})

watch(
  () => route.params.component,
  () => {
    activeServiceKey.value = componentParam()
    activeView.value = 'config'
  },
)

hydrateServiceDefaults()

function handleCardClick(service: Service) {
  if (service.disabled) {
    if (service.key !== 'new') ElMessage.warning(`${service.name} 功能规划中，敬请期待`)
    return
  }
  router.push({ name: 'ServiceDelivery', params: { component: service.key } })
}

function componentParam() {
  const param = route.params.component
  return typeof param === 'string' ? param : ''
}

function backToCatalog() {
  router.push({ name: 'ServiceCatalog' })
}

function hydrateServiceDefaults() {
  const service = activeService.value
  if (!service) return
  const prefix = service.key === 'mysql' ? 'mysql' : 'nginx'
  deliveryForm.instanceName = `${prefix}-${currentName.value}-billing-02`
  deliveryForm.version = service.versions[0] || ''
  deliveryForm.mode = service.modes[1]?.value || service.modes[0]?.value || 'single'
  deliveryForm.spec = service.specs[1] || service.specs[0] || ''
  deliveryForm.disk = service.disks[1] || service.disks[0] || ''
  deliveryForm.port = service.defaultPort
  deliveryForm.charset = service.charsets[0] || ''
  deliveryForm.installPath = service.defaultPaths.install
  deliveryForm.dataPath = service.defaultPaths.data
  deliveryForm.logPath = service.defaultPaths.log
  resetExecutionState()
}

function resetWorkbench() {
  hydrateServiceDefaults()
  activeView.value = 'config'
  ElMessage.success('交付工作台已重置')
}

function resetExecutionState() {
  closeEventSource()
  precheckPassed.value = false
  running.value = false
  deliveryDone.value = false
  deliveryFailed.value = false
  canceling.value = false
  deploymentId.value = ''
  deliveryLog.value = '[ready] 等待创建交付任务...'
  steps.value = defaultSteps().map((step) => ({ ...step, state: 'pending' }))
}

function defaultSteps(): DeliveryStep[] {
  const isMysql = activeServiceKey.value === 'mysql'
  const copy = isMysql
    ? [
        ['资源锁定与基线检查', '主机、端口、目录、旧版本配置'],
        ['安装 MySQL 与初始化目录', '软件包、运行用户、数据盘'],
        ['应用标准参数与复制配置', '基础参数、GTID、主从关系'],
        ['数据库健康检查', '端口、读写、复制延迟、参数基线'],
        ['注册资产与交付归档', 'CloudDM、CMDB、Prometheus'],
      ]
    : [
        ['资源锁定与配置检查', '节点、端口、证书、Ingress 冲突'],
        ['安装 OpenResty 与初始化目录', '软件包、运行用户、日志目录'],
        ['下发标准网关配置', 'upstream、灰度路由、限流参数'],
        ['网关健康检查', '端口、HTTP 状态、访问延迟'],
        ['注册资产与交付归档', 'CMDB、Prometheus、日志平台'],
      ]
  return copy.map(([name, desc]) => ({ name, desc, state: 'pending' as StepState }))
}

function precheck() {
  precheckPassed.value = true
  deliveryFailed.value = false
  deliveryLog.value = [
    '[precheck] resource lock acquired',
    `[precheck] ${deliveryForm.instanceName} naming rule passed`,
    `[precheck] ${resultAddress.value} port is available`,
    '[precheck] baseline checks passed',
  ].join('\n')
  ElMessage.success('预检查通过，可以创建交付任务')
}

async function createTask() {
  if (!precheckPassed.value) {
    ElMessage.warning('请先执行预检查')
    return
  }
  const businessLineId = businessLineStore.current?.id
  if (!businessLineId) {
    ElMessage.warning('请先选择业务线')
    return
  }
  if (!activeService.value) return
  running.value = true
  deliveryDone.value = false
  deliveryFailed.value = false
  activeView.value = 'execution'
  steps.value = defaultSteps().map((step) => ({ ...step, state: 'pending' }))
  try {
    const result = await deploymentApi.create({
      component: activeService.value.key,
      business_line_id: businessLineId,
      params: deploymentParams(),
    })
    deploymentId.value = result.deployment_id
    deliveryLog.value += `\n[task] ${result.deployment_id} created by ${currentName.value}`
    connectDeploymentEvents(result.deployment_id)
    ElMessage.success('交付任务已创建')
  } catch (error) {
    running.value = false
    deliveryFailed.value = true
    activeView.value = 'result'
    ElMessage.error(error instanceof Error ? error.message : '创建交付任务失败')
  }
}

async function cancelDeployment() {
  if (!deploymentId.value) return
  canceling.value = true
  try {
    await deploymentApi.cancel(deploymentId.value)
    appendLog('[cancel] cancel requested')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '取消任务失败')
  } finally {
    canceling.value = false
  }
}

function deploymentParams() {
  return {
    business_line: currentName.value,
    instance_name: deliveryForm.instanceName,
    version: deliveryForm.version,
    mode: deliveryForm.mode,
    mode_label: currentModeLabel.value,
    spec: deliveryForm.spec,
    disk: deliveryForm.disk,
    port: deliveryForm.port,
    charset: deliveryForm.charset,
    pool: deliveryForm.pool,
    install_path: deliveryForm.installPath,
    data_path: deliveryForm.dataPath,
    log_path: deliveryForm.logPath,
  }
}

function connectDeploymentEvents(id: string) {
  closeEventSource()
  const source = new EventSource(deploymentApi.eventsURL(id))
  eventSource.value = source
  source.addEventListener('log', (event) => {
    const data = parseSSEData(event)
    appendLog(data.message)
    markStepRunning()
  })
  source.addEventListener('status', (event) => {
    const data = parseSSEData(event)
    applyDeploymentStatus(String(data.status || ''), data.message ? String(data.message) : '')
  })
  source.addEventListener('error', (event) => {
    const data = parseSSEData(event)
    if (data.message) appendLog(String(data.message))
    applyDeploymentStatus(String(data.status || 'failed'), '')
  })
  source.addEventListener('result', (event) => {
    const data = parseSSEData(event)
    applyDeploymentStatus(String(data.status || 'success'), data.message ? String(data.message) : '')
  })
  source.addEventListener('done', (event) => {
    const data = parseSSEData(event)
    applyDeploymentStatus(String(data.status || ''), '')
    closeEventSource()
    activeView.value = 'result'
  })
  source.onerror = () => {
    appendLog('[sse] connection interrupted')
  }
}

function closeEventSource() {
  eventSource.value?.close()
  eventSource.value = null
}

function parseSSEData(event: Event): Record<string, unknown> {
  const message = event as MessageEvent<string>
  try {
    return JSON.parse(message.data || '{}')
  } catch {
    return {}
  }
}

function appendLog(message: unknown) {
  if (!message) return
  deliveryLog.value += `\n${String(message)}`
}

function applyDeploymentStatus(status: string, message: string) {
  if (message) appendLog(message)
  if (status === 'running') {
    running.value = true
    markStepRunning()
    return
  }
  if (status === 'success') {
    running.value = false
    deliveryDone.value = true
    deliveryFailed.value = false
    steps.value = steps.value.map((step) => ({ ...step, state: 'done' }))
    return
  }
  if (status === 'failed' || status === 'canceled') {
    running.value = false
    deliveryDone.value = false
    deliveryFailed.value = true
    markCurrentStepFailed()
  }
}

function markStepRunning() {
  const index = steps.value.findIndex((step) => step.state === 'pending' || step.state === 'running')
  if (index < 0) return
  steps.value = steps.value.map((step, stepIndex) => {
    if (stepIndex < index) return { ...step, state: 'done' }
    if (stepIndex === index) return { ...step, state: 'running' }
    return step
  })
}

function markCurrentStepFailed() {
  const index = steps.value.findIndex((step) => step.state === 'running' || step.state === 'pending')
  if (index < 0) return
  steps.value = steps.value.map((step, stepIndex) => {
    if (stepIndex < index) return { ...step, state: 'done' }
    if (stepIndex === index) return { ...step, state: 'failed' }
    return step
  })
}

function stepStatusText(state: StepState) {
  const map: Record<StepState, string> = {
    pending: '等待',
    running: '执行中',
    done: '成功',
    failed: '失败',
  }
  return map[state]
}

function stepSymbol(state: StepState) {
  const map: Record<StepState, string> = {
    pending: '○',
    running: '…',
    done: '✓',
    failed: '×',
  }
  return map[state]
}

</script>

<style scoped>
.svc-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.svc-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: border-color 0.15s, transform 0.15s;
}

.svc-card:hover:not(.disabled),
.svc-card.active {
  border-color: var(--accent);
}

.svc-card.active {
  box-shadow: 0 0 0 1px var(--accent-dim) inset;
}

.svc-card.disabled {
  cursor: not-allowed;
}

.svc-top,
.ext-top,
.section-head,
.execution-head,
.result-head,
.head-actions,
.form-actions {
  display: flex;
  align-items: center;
}

.svc-top {
  gap: 10px;
  margin-bottom: 10px;
}

.svc-ic,
.ext-logo,
.node-icon,
.result-icon,
.step-dot {
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-weight: 600;
}

.svc-ic {
  width: 34px;
  height: 34px;
  border-radius: 7px;
  font-size: 13px;
}

h4 {
  margin: 0;
  font-size: 13.5px;
}

.svc-status {
  display: block;
  margin-top: 2px;
  font-size: 11px;
  color: var(--text-dim);
}

.svc-desc {
  min-height: 36px;
  margin-bottom: 12px;
  color: var(--text-dim);
  font-size: 12px;
  line-height: 1.5;
}

.svc-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-dim);
  font-size: 11px;
}

.playbook {
  color: var(--text-mid);
  font-family: var(--mono);
}

.delivery-panel {
  margin-top: 6px;
}

.workbench-head {
  align-items: center;
}

.workbench-head h3 {
  margin-bottom: 4px;
}

.head-actions {
  gap: 8px;
  margin-left: auto;
}

.flow-tabs {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 16px;
}

.flow-tab {
  min-height: 58px;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--bg-panel-2);
  color: var(--text-mid);
  text-align: left;
}

.flow-tab strong,
.flow-tab span {
  display: block;
}

.flow-tab strong {
  color: var(--text-hi);
  font-size: 13px;
}

.flow-tab span {
  margin-top: 3px;
  color: var(--text-dim);
  font-size: 11.5px;
}

.flow-tab.active {
  border-color: var(--accent);
  background: var(--accent-dim);
}

.workbench-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(340px, 0.9fr);
  gap: 14px;
}

.config-panel,
.preview-panel,
.execution-view,
.result-view {
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--bg-panel);
  padding: 16px;
}

.section-head,
.execution-head {
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.section-head h3,
.execution-head h3 {
  margin: 0 0 4px;
  font-size: 14px;
}

.section-head p,
.execution-head p,
.result-head p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12px;
}

.form-grid,
.subform-grid {
  display: grid;
  gap: 12px;
}

.form-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-bottom: 12px;
}

.subform-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.form-field {
  display: grid;
  gap: 6px;
  color: var(--text-dim);
  font-size: 11.5px;
}

.config-list {
  border-top: 1px solid var(--line-soft);
}

.config-row {
  display: grid;
  grid-template-columns: 92px minmax(0, 1fr);
  gap: 12px;
  padding: 13px 0;
  border-bottom: 1px solid var(--line-soft);
}

.config-row:last-child {
  border-bottom: 0;
}

.config-label {
  padding-top: 8px;
  color: var(--text-dim);
  font-size: 12px;
}

.choice-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.choice-btn {
  height: 32px;
  padding: 0 12px;
  border: 1px solid var(--line);
  border-radius: 6px;
  background: var(--bg-panel-2);
  color: var(--text-mid);
  font-size: 12px;
}

.choice-btn.active {
  border-color: var(--accent);
  background: var(--accent);
  color: #06150F;
  font-weight: 600;
}

.topology {
  display: grid;
  gap: 10px;
  margin-bottom: 14px;
}

.topology.nodes-1 {
  grid-template-columns: minmax(0, 1fr);
}

.topology.nodes-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.topology.nodes-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.node-card {
  min-width: 0;
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--bg-panel-2);
}

.node-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.node-icon {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  background: var(--bg-panel);
  color: var(--accent);
  font-size: 11px;
}

.node-meta {
  display: grid;
  gap: 4px;
  color: var(--text-dim);
  font-family: var(--mono);
  font-size: 11px;
}

.yaml-preview,
.log-preview {
  overflow-x: auto;
  border: 1px solid var(--line-soft);
  border-radius: 6px;
  background: #0B0D12;
  color: #8EC8FF;
  font-family: var(--mono);
  font-size: 11.5px;
  line-height: 1.7;
  white-space: pre;
}

.yaml-preview {
  min-height: 180px;
  margin: 0 0 12px;
  padding: 12px 14px;
}

.log-preview {
  min-height: 170px;
  margin: 14px 0 0;
  padding: 12px 14px;
}

.form-actions {
  justify-content: flex-end;
  gap: 8px;
}

.step-list {
  display: grid;
}

.delivery-step {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
  padding: 13px 0;
  border-bottom: 1px solid var(--line-soft);
}

.delivery-step:last-child {
  border-bottom: 0;
}

.step-dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bg-panel-2);
  color: var(--text-dim);
  font-size: 13px;
}

.delivery-step.running .step-dot {
  color: var(--accent);
  background: var(--accent-dim);
}

.delivery-step.done .step-dot {
  color: #06150F;
  background: var(--accent);
}

.delivery-step.failed .step-dot {
  color: #FFFFFF;
  background: var(--err);
}

.delivery-step strong {
  display: block;
  margin-bottom: 3px;
  font-size: 13px;
}

.delivery-step p {
  margin: 0;
  color: var(--text-dim);
  font-size: 12px;
}

.result-head {
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--line-soft);
}

.result-icon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: var(--accent);
  color: #06150F;
  font-size: 22px;
}

.result-head.failed .result-icon {
  background: var(--err);
  color: #FFFFFF;
}

.result-head h3 {
  margin: 0 0 4px;
  font-size: 16px;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  padding: 16px 0;
}

.result-item {
  min-width: 0;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--bg-panel-2);
}

.result-item span {
  display: block;
  margin-bottom: 5px;
  color: var(--text-dim);
  font-size: 11.5px;
}

.result-item strong {
  overflow-wrap: anywhere;
  font-family: var(--mono);
  font-size: 12.5px;
}

.asset-table {
  overflow-x: auto;
}

.asset-table table {
  width: 100%;
  border-collapse: collapse;
}

.asset-table th,
.asset-table td {
  padding: 10px 8px;
  border-top: 1px solid var(--line-soft);
  color: var(--text-mid);
  font-size: 12px;
  text-align: left;
  white-space: nowrap;
}

.asset-table th {
  color: var(--text-dim);
  font-weight: 600;
}

.tag-green {
  color: var(--tag-green-text);
  border-color: var(--tag-green-border);
  background: var(--tag-green-bg);
}

.tag-amber {
  color: var(--tag-amber-text);
  border-color: var(--tag-amber-border);
  background: var(--tag-amber-bg);
}

.tag-red {
  color: var(--err);
  border-color: var(--err);
  background: var(--logo-cc-bg);
}

.container-head {
  margin-top: 24px;
}

.ext-card {
  max-width: 560px;
  display: grid;
  gap: 12px;
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--bg-panel-2);
}

.ext-logo {
  width: 38px;
  height: 38px;
  border-radius: 8px;
}

.wayne-logo {
  background: var(--logo-w-bg);
  color: var(--tag-blue-text);
}

.ext-sub,
.sso-row {
  color: var(--text-dim);
  font-size: 12px;
}

.ext-card p {
  margin: 0;
  color: var(--text-mid);
  font-size: 12.5px;
  line-height: 1.6;
}

.sso-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.sso-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 6px var(--accent-glow);
}

code {
  font-family: var(--mono);
}

@media (max-width: 1180px) {
  .workbench-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .page-head,
  .section-head,
  .execution-head,
  .workbench-head {
    align-items: stretch;
    flex-direction: column;
  }

  .head-actions {
    margin-left: 0;
  }

  .flow-tabs,
  .form-grid,
  .subform-grid,
  .topology,
  .result-grid {
    grid-template-columns: 1fr;
  }

  .config-row {
    grid-template-columns: 1fr;
  }

  .config-label {
    padding-top: 0;
  }

  .delivery-step {
    grid-template-columns: 28px minmax(0, 1fr);
  }

  .delivery-step .tag {
    grid-column: 2;
    justify-self: start;
  }
}
</style>
