import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { businessLineApi, type BusinessLine as ApiBusinessLine } from '@/api/businessLine'

export interface BusinessLine {
  id: number
  name: string
  ou: string
  role: string
  iconText: string
  iconBg: string
  iconColor: string
  authorized: boolean
  permission?: 0 | 1
  created_at?: string
  updated_at?: string
}

const STORAGE_KEY = 'xinfra-current-bl'
const COLORS = [
  ['#1A1430', '#C9A6FF'],
  ['#0E1C2C', '#8EC8FF'],
  ['#241B0A', '#FFC97A'],
  ['#1A2A1A', '#A6FFB0'],
  ['#2A1A2A', '#FFB0E0'],
] as const

function loadSavedBL(): BusinessLine | null {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) return JSON.parse(saved)
  } catch {
    // ignore
  }
  return null
}

function saveBL(bl: BusinessLine | null) {
  if (!bl) {
    localStorage.removeItem(STORAGE_KEY)
    return
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(bl))
}

function initials(name: string) {
  const compact = name.replace(/\s+/g, '')
  return compact.slice(0, 2).toUpperCase() || 'BL'
}

function mapBusinessLine(item: ApiBusinessLine, index: number): BusinessLine {
  const [iconBg, iconColor] = COLORS[index % COLORS.length]
  return {
    id: item.id,
    name: item.name,
    ou: `ou=${item.name}`,
    role: item.permission === 0 ? '管理员' : '普通用户',
    iconText: initials(item.name),
    iconBg,
    iconColor,
    authorized: true,
    permission: item.permission,
    created_at: item.created_at,
    updated_at: item.updated_at,
  }
}

export const useBusinessLineStore = defineStore('businessLine', () => {
  const businessLines = ref<BusinessLine[]>([])
  const currentBL = ref<BusinessLine | null>(loadSavedBL())
  const loading = ref(false)

  const items = computed(() => businessLines.value)
  const current = computed(() => currentBL.value)
  const isCurrentAdmin = computed(() => currentBL.value?.permission === 0)

  const authorizedBLs = computed(() =>
    businessLines.value.filter((bl) => bl.authorized),
  )

  async function loadMine() {
    loading.value = true
    try {
      const rows = await businessLineApi.listMine()
      const mapped = rows.map(mapBusinessLine)
      businessLines.value = mapped

      const savedID = currentBL.value?.id
      const next = mapped.find((bl) => bl.id === savedID) || mapped[0] || null
      currentBL.value = next
      saveBL(next)
    } finally {
      loading.value = false
    }
  }

  function switchBL(id: number) {
    const bl = businessLines.value.find((b) => b.id === id)
    if (!bl || !bl.authorized) return
    currentBL.value = bl
    saveBL(bl)
  }

  function clear() {
    businessLines.value = []
    currentBL.value = null
    saveBL(null)
  }

  return {
    businessLines,
    items,
    currentBL,
    current,
    loading,
    isCurrentAdmin,
    authorizedBLs,
    loadMine,
    switchBL,
    clear,
  }
})
