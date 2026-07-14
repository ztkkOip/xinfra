<template>
  <div>
    <div class="page-head">
      <div>
        <h1>基础服务目录</h1>
        <p>点击卡片，通过封装好的 ansible-playbook 一键部署标准化基础服务</p>
      </div>
    </div>

    <div class="svc-grid">
      <div v-for="service in services" :key="service.name" class="svc-card" @click="handleDeploy(service)">
        <div class="svc-top">
          <div class="svc-ic" :style="{ background: service.bgColor, color: service.iconColor }">
            {{ service.icon }}
          </div>
          <h4>{{ service.name }}</h4>
        </div>
        <div class="svc-desc">{{ service.description }}</div>
        <div class="svc-foot">
          <span class="playbook">{{ service.playbook }}</span>
          <span>{{ service.version }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const services = ref([
  { name: 'MySQL', icon: 'My', bgColor: '#1C2E4A', iconColor: '#7FB8FF', description: '一主多从，自动注册进 open-cdm 数据源，支持 LDAP 账号统一审核 SQL 上线。', playbook: 'roles/mysql-deploy', version: 'v8.0.36' },
  { name: 'Redis', icon: 'Rd', bgColor: '#3A1C28', iconColor: '#FF8E9C', description: '支持 Standalone / Sentinel / Cluster，部署完成后自动注册进 CacheCloud。', playbook: 'roles/redis-deploy', version: 'v7.2' },
  { name: 'openresty', icon: 'Or', bgColor: '#1C3A2E', iconColor: '#7FFFC2', description: '容器化部署在 K8s 内，作为业务边缘网关，支持灰度路由配置。', playbook: 'roles/openresty-deploy', version: 'v1.25' },
  { name: 'dpvs', icon: 'Dp', bgColor: '#3A2E1C', iconColor: '#FFC97A', description: '四层负载均衡，直接管理宿主机网络命名空间，不纳入容器编排。', playbook: 'roles/dpvs-deploy', version: 'v1.9' },
  { name: 'PgSQL', icon: 'Pg', bgColor: '#2E1C3A', iconColor: '#C9A6FF', description: '流复制主从架构，规划中 · 第二期接入 open-cdm 统一审核体系。', playbook: 'roles/pgsql-deploy', version: '规划中' },
])

const handleDeploy = (service: any) => {
  ElMessage.success(`${service.name} 部署任务已提交`)
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

.svc-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
}

.svc-card {
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.svc-card:hover {
  border-color: #3A4356;
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
</style>
