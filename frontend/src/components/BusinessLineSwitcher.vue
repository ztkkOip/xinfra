<template>
  <div class="bl-switch" ref="switchRef">
    <div class="bl-trigger" @click="toggleMenu">
      <div
        class="bl-ic"
        :style="{ background: currentBL?.iconBg, color: currentBL?.iconColor }"
      >
        {{ currentBL?.iconText }}
      </div>
      <span class="bl-name-trigger">{{ currentBL?.name || '未选择' }}</span>
      <span class="bl-role" v-if="currentBL">{{ currentBL.role }} · 授权 {{ authCount }} 子系统</span>
      <span class="caret">▼</span>
    </div>
    <Transition name="bl-menu">
      <div class="bl-menu" v-show="menuVisible">
        <div class="bl-menu-label">切换业务线（仅展示被授权范围）</div>
        <div
          v-for="bl in businessLines"
          :key="bl.id"
          class="bl-opt"
          :class="{ sel: bl.id === currentBL?.id, locked: !bl.authorized }"
          @click="handleSelect(bl)"
        >
          <div
            class="bl-ic"
            :style="{ background: bl.iconBg, color: bl.iconColor }"
          >
            {{ bl.iconText }}
          </div>
          <div>
            <div class="bl-name">{{ bl.name }}</div>
            <div class="bl-sub">{{ bl.ou }} · {{ bl.role }}</div>
          </div>
          <span v-if="bl.authorized" class="bl-check">✓</span>
          <span v-else class="bl-lock">🔒 无权限</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useBusinessLineStore, type BusinessLine } from '@/stores/businessLine'

const blStore = useBusinessLineStore()

const businessLines = computed(() => blStore.businessLines)
const currentBL = computed(() => blStore.currentBL)
const authCount = computed(() => blStore.authorizedBLs.length)

const menuVisible = ref(false)
const switchRef = ref<HTMLElement | null>(null)

function toggleMenu() {
  menuVisible.value = !menuVisible.value
}

function handleSelect(bl: BusinessLine) {
  if (!bl.authorized) return
  blStore.switchBL(bl.id)
  menuVisible.value = false
  ElMessage.success(`已切换业务线：${bl.name}（资源与告警按授权范围过滤）`)
}

function onClickOutside(e: MouseEvent) {
  if (switchRef.value && !switchRef.value.contains(e.target as Node)) {
    menuVisible.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', onClickOutside)
})
</script>

<style scoped>
.bl-switch {
  position: relative;
}

.bl-trigger {
  display: flex;
  align-items: center;
  gap: 9px;
  background: var(--bg-panel-2);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 6px 12px;
  color: var(--text-hi);
  font-size: 12.5px;
  font-weight: 500;
  cursor: pointer;
  transition: border-color 0.2s;
}

.bl-trigger:hover {
  border-color: var(--accent);
}

.bl-trigger .bl-ic {
  width: 20px;
  height: 20px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 10px;
  font-weight: 600;
}

.bl-trigger .caret {
  color: var(--text-dim);
  font-size: 10px;
  margin-left: 2px;
}

.bl-trigger .bl-role {
  font-family: var(--mono);
  font-size: 10px;
  color: var(--text-dim);
}

.bl-name-trigger {
  font-size: 12.5px;
}

.bl-menu {
  position: absolute;
  top: 44px;
  left: 0;
  width: 280px;
  z-index: 60;
  background: var(--bg-panel);
  border: 1px solid var(--line);
  border-radius: 10px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.18);
  padding: 6px;
}

.bl-menu-label {
  font-size: 10px;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 8px 10px 4px;
}

.bl-opt {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 7px;
  cursor: pointer;
  transition: background 0.15s;
}

.bl-opt:hover {
  background: var(--bg-panel-2);
}

.bl-opt.sel {
  background: var(--accent-dim);
}

.bl-opt .bl-ic {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 600;
}

.bl-opt .bl-name {
  font-size: 13px;
  color: var(--text-hi);
  font-weight: 500;
}

.bl-opt .bl-sub {
  font-size: 10.5px;
  color: var(--text-dim);
  font-family: var(--mono);
}

.bl-opt .bl-check {
  margin-left: auto;
  color: var(--accent);
  font-size: 13px;
  opacity: 0;
}

.bl-opt.sel .bl-check {
  opacity: 1;
}

.bl-opt.locked {
  opacity: 0.45;
  cursor: not-allowed;
}

.bl-opt.locked .bl-lock {
  margin-left: auto;
  color: var(--text-dim);
  font-size: 11px;
}

/* 下拉动画 */
.bl-menu-enter-active,
.bl-menu-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.bl-menu-enter-from,
.bl-menu-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
