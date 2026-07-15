<template>
  <div>
    <div class="page-head">
      <div>
        <h1>任务中心</h1>
        <p>所有 ansible-playbook / Wayne 发布任务的执行记录与实时日志</p>
      </div>
    </div>

    <div class="cols-2-equal">
      <div class="panel">
        <div class="panel-head">
          <h3>任务列表</h3>
        </div>
        <div class="panel-body">
          <table>
            <tbody>
              <tr v-for="task in tasks" :key="task.id" class="tr-hover" :class="{ active: task.id === selectedTask }" @click="selectedTask = task.id">
                <td :class="['status-text', task.statusClass]">● {{ task.status }}</td>
                <td class="strong">{{ task.name }}</td>
                <td class="mono text-xs text-dim">{{ task.playbook }}</td>
              </tr>
            </tbody>
          </table>
          <div class="pagination">
            <span>共 86 条 · 每页 6 条</span>
            <div class="pg-btns">
              <span class="pg-btn disabled">‹</span>
              <span class="pg-btn active">1</span>
              <span class="pg-btn">2</span>
              <span class="pg-btn">3</span>
              <span class="pg-sep">…</span>
              <span class="pg-btn">15</span>
              <span class="pg-btn">›</span>
            </div>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h3>实时日志 · RKE2 节点加入</h3>
          <span class="meta">WebSocket 流式输出</span>
        </div>
        <div class="panel-body log-stream">
          <div v-for="(log, index) in logs" :key="index" :class="['task-log-line', log.class]">
            <span class="t">{{ log.time }}</span>{{ log.message }}
          </div>
          <div class="task-log-line">
            <span class="t">10:42:33</span>
            <span class="blink">▌</span> 等待节点 Ready...
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const selectedTask = ref(1)

const tasks = ref([
  { id: 1, name: 'RKE2 节点加入 · bj-node-061', playbook: 'roles/rke2-node-join', status: '执行中', statusClass: 'warn' },
  { id: 2, name: 'MySQL 主从部署 · kodo', playbook: 'roles/mysql-deploy', status: '成功', statusClass: 'ok' },
  { id: 3, name: 'Redis Cluster 部署 · las', playbook: 'roles/redis-deploy', status: '成功', statusClass: 'ok' },
  { id: 4, name: 'openresty 网关部署', playbook: 'roles/openresty-deploy', status: '失败', statusClass: 'err' },
  { id: 5, name: '业务线标签同步 · LDAP', playbook: 'internal/label-sync', status: '成功', statusClass: 'ok' },
  { id: 6, name: 'MySQL 主从部署 · ltoken', playbook: 'roles/mysql-deploy', status: '成功', statusClass: 'ok' },
])

const logs = ref([
  { time: '10:42:01', message: 'PLAY [rke2-node-join] **********************', class: '' },
  { time: '10:42:02', message: 'TASK [初始化内核参数] ...', class: '' },
  { time: '10:42:04', message: 'ok: [bj-node-061]', class: 'ok' },
  { time: '10:42:05', message: 'TASK [安装 containerd] ...', class: '' },
  { time: '10:42:18', message: 'ok: [bj-node-061]', class: 'ok' },
  { time: '10:42:19', message: 'TASK [写入 RKE2 node-labels: business-line=kodo] ...', class: '' },
  { time: '10:42:20', message: 'changed: [bj-node-061] => labels applied', class: 'tag-ok' },
  { time: '10:42:21', message: 'TASK [加入集群 rke2-bj-prod-01] ...', class: '' },
])
</script>

<style scoped>
/* 公共样式已在 global.css 中定义 */

.log-stream {
  max-height: 340px;
  overflow-y: auto;
}

.blink {
  animation: blink 1s step-start infinite;
  color: var(--accent);
}

@keyframes blink {
  50% { opacity: 0; }
}
</style>
