<template>
  <aside class="sidebar">
    <div class="nav-group">
      <div class="nav-label">概览</div>
      <router-link to="/dashboard" class="nav-item" active-class="active">
        <span class="ic">▣</span>资源大盘
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">机房 &amp; 网络</div>
      <router-link to="/resource/status" class="nav-item" active-class="active">
        <span class="ic">⬡</span>资源状态看板
      </router-link>
      <router-link to="/cluster" class="nav-item" active-class="active">
        <span class="ic">◆</span>集群 / 节点<span class="badge">7</span>
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">资源管理</div>
      <router-link to="/resource/mgmt" class="nav-item" active-class="active">
        <span class="ic">▥</span>资源管理<span class="badge">CMDB</span>
      </router-link>
      <router-link to="/service/mgmt" class="nav-item" active-class="active">
        <span class="ic">◈</span>服务管理<span class="badge">Consul</span>
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">服务交付</div>
      <router-link to="/service/catalog" class="nav-item" active-class="active">
        <span class="ic">▤</span>基础服务目录
      </router-link>
      <router-link to="/tasks" class="nav-item" active-class="active">
        <span class="ic">▧</span>任务中心<span class="badge">2</span>
      </router-link>
      <router-link to="/config" class="nav-item" active-class="active">
        <span class="ic">⚙</span>配置中心管理<span class="badge">Apollo</span>
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">可观测性</div>
      <router-link to="/monitor" class="nav-item" active-class="active">
        <span class="ic">◎</span>监控 / 日志 / 告警
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">统一入口</div>
      <router-link to="/subsystem" class="nav-item" active-class="active">
        <span class="ic">↗</span>子系统导航
      </router-link>
      <router-link to="/subsystem/authz" class="nav-item" active-class="active">
        <span class="ic">▦</span>子系统赋权
      </router-link>
    </div>
    <div class="nav-group">
      <div class="nav-label">审计</div>
      <router-link to="/audit/login" class="nav-item" active-class="active">
        <span class="ic">📋</span>登录审计
      </router-link>
      <router-link to="/audit/ops" class="nav-item" active-class="active">
        <span class="ic">📝</span>运维操作审计
      </router-link>
    </div>
    <div class="sidebar-foot">
      xinfra v3 一期 · 七机房容器化<br>RKE2 · Calico BGP · Ceph · LAS<br>AWS TGW + WireGuard · LDAP 认证
    </div>
  </aside>
</template>

<script setup lang="ts">
</script>

<style scoped>
.sidebar {
  background: var(--bg-panel);
  border-right: 1px solid var(--line);
  padding: 14px 10px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.nav-group {
  margin-bottom: 18px;
}

.nav-label {
  font-size: 13px;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 0 10px;
  margin-bottom: 6px;
  font-weight: 600;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px;
  border-radius: 6px;
  color: var(--text-mid);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 2px;
  border: 1px solid transparent;
  text-decoration: none;
}

.nav-item .ic {
  width: 16px;
  text-align: center;
  font-size: 16px;
  opacity: 0.85;
}

.nav-item:hover {
  background: var(--bg-panel-2);
  color: var(--text-hi);
}

.nav-item.active {
  background: var(--accent-dim);
  color: var(--accent);
  border-color: var(--accent);
}

.nav-item .badge {
  margin-left: auto;
  font-family: var(--mono);
  font-size: 12px;
  background: var(--bg-panel-2);
  color: var(--text-dim);
  padding: 1px 6px;
  border-radius: 10px;
  border: 1px solid var(--line);
}

.nav-item.active .badge {
  background: var(--accent);
  color: #FFFFFF;
  border-color: var(--accent);
}

.sidebar-foot {
  margin-top: auto;
  padding: 10px;
  border-top: 1px solid var(--line-soft);
  font-size: 14px;
  color: var(--text-dim);
  line-height: 1.6;
}
</style>
