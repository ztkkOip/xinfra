<template>
  <div>
    <div class="page-head">
      <div>
        <h1>资源大盘</h1>
        <p>跨机房 / 多云容器资源统一视图 · 最近更新于 12 秒前</p>
      </div>
      <el-button @click="refresh">⟳ 刷新</el-button>
    </div>

    <div class="stat-row">
      <div class="stat-card">
        <div class="label">RKE2 集群</div>
        <div class="value">{{ profile.clusters }}</div>
        <div class="delta">{{ profile.clusters }} 机房在线</div>
      </div>
      <div class="stat-card">
        <div class="label">节点总数</div>
        <div class="value">{{ profile.nodes }}</div>
        <div class="delta up">↑ {{ Math.max(1, Math.round(profile.nodes / 24)) }} 本周新增</div>
      </div>
      <div class="stat-card">
        <div class="label">CPU 已分配</div>
        <div class="value">{{ profile.cpuAllocated }}%</div>
        <div class="delta">{{ Math.round(profile.nodes * 16 * profile.cpuAllocated / 100) }} / {{ profile.nodes * 16 }} 核</div>
      </div>
      <div class="stat-card">
        <div class="label">组件实例</div>
        <div class="value">{{ profile.components }}</div>
        <div class="delta">MySQL {{ profile.mysql }} · Redis {{ profile.redis }} · 其他 {{ profile.components - profile.mysql - profile.redis }}</div>
      </div>
      <div class="stat-card">
        <div class="label">进行中任务</div>
        <div class="value warn">{{ profile.tasksRunning }}</div>
        <div class="delta warn">{{ Math.min(profile.tasksRunning, 1) }} 个待关注</div>
      </div>
    </div>

    <div class="cols-2">
      <div class="panel">
        <div class="panel-head">
          <h3>集群拓扑 · 按机房 / 业务线标签</h3>
          <span class="meta">node-label 调度视图</span>
        </div>
        <div class="topo">
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-a">{{ profile.primaryZone }}</span>
              <span class="mono text-dim">{{ Math.round(profile.nodes * 0.42) }}节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 24" :key="i" :class="['node-cell', i <= 18 ? 'a' : '']"></div>
            </div>
          </div>
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-b">{{ profile.secondaryZone }}</span>
              <span class="mono text-dim">{{ Math.round(profile.nodes * 0.36) }}节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 24" :key="i" :class="['node-cell', i <= 20 ? 'b' : '']"></div>
            </div>
          </div>
          <div class="topo-cluster">
            <div class="tname">
              <span class="tag zone-c">弹性云资源池</span>
              <span class="mono text-dim">{{ Math.max(1, profile.nodes - Math.round(profile.nodes * 0.42) - Math.round(profile.nodes * 0.36)) }}节点</span>
            </div>
            <div class="node-grid">
              <div v-for="i in 18" :key="i" :class="['node-cell', i <= 12 ? 'c' : '']"></div>
            </div>
          </div>
        </div>
        <div class="legend">
          <span><span class="node-cell a inline"></span>{{ currentName }}主资源池</span>
          <span><span class="node-cell b inline"></span>{{ currentName }}扩展池</span>
          <span><span class="node-cell c inline"></span>{{ currentName }}弹性池</span>
          <span><span class="node-cell empty inline"></span>空闲</span>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>最近任务</h3>
          <router-link to="/task/log" class="meta">查看全部 →</router-link>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">{{ profile.servicePrefix }} MySQL 主从部署</td>
                <td class="mono">2m14s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text warn">● 执行中</td>
                <td class="strong">{{ profile.servicePrefix }} RKE2 节点加入</td>
                <td class="mono">38s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">{{ profile.servicePrefix }} Redis Cluster 部署</td>
                <td class="mono">3m02s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text err">● 失败</td>
                <td class="strong">{{ profile.servicePrefix }} openresty 网关部署</td>
                <td class="mono">41s</td>
              </tr>
              <tr class="tr-hover">
                <td class="status-text ok">● 成功</td>
                <td class="strong">业务线标签同步</td>
                <td class="mono">5s</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useBusinessLineMockProfile } from '@/utils/businessLineMock'

const { currentName, profile } = useBusinessLineMockProfile()

const refresh = () => {
  // 刷新数据
}
</script>

<style scoped>
/* 公共样式(.page-head, .stat-row, .cols-2, .panel, table, .tag, .status-text)已在 global.css 中定义 */

.topo {
  padding: 18px 16px;
  display: flex;
  gap: 14px;
}

.topo-cluster {
  flex: 1;
  border: 1px dashed var(--line);
  border-radius: var(--radius-lg);
  padding: 12px;
}

.topo-cluster .tname {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.node-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 5px;
}

.node-cell {
  aspect-ratio: 1;
  border-radius: 3px;
  background: var(--node-empty-bg);
  border: 1px solid var(--line-soft);
}

.node-cell.a {
  background: var(--node-a-bg);
  border-color: var(--node-a-border);
}

.node-cell.b {
  background: var(--node-b-bg);
  border-color: var(--node-b-border);
}

.node-cell.c {
  background: var(--node-c-bg);
  border-color: var(--node-c-border);
}

.node-cell.empty {
  background: transparent;
  border-style: dashed;
}

.node-cell.inline {
  display: inline-block;
  width: 9px;
  height: 9px;
  margin-right: 5px;
}

.legend {
  padding: 0 16px 14px;
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: var(--text-dim);
}
</style>
