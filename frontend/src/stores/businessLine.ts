import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface BusinessLine {
  id: string
  name: string
  ou: string
  role: string
  iconText: string
  iconBg: string
  iconColor: string
  authorized: boolean
}

const STORAGE_KEY = 'xinfra-current-bl'

// TODO: 替换为真实 API 调用
const MOCK_BUSINESS_LINES: BusinessLine[] = [
  {
    id: 'las',
    name: 'LAS 业务线',
    ou: 'ou=las',
    role: 'SRE',
    iconText: 'LA',
    iconBg: '#1A1430',
    iconColor: '#C9A6FF',
    authorized: true,
  },
  {
    id: 'kodo',
    name: 'Kodo 业务线',
    ou: 'ou=kodo',
    role: '研发',
    iconText: 'KO',
    iconBg: '#0E1C2C',
    iconColor: '#8EC8FF',
    authorized: true,
  },
  {
    id: 'lingxi',
    name: '灵矽 业务线',
    ou: 'ou=lingxi',
    role: '未授权',
    iconText: 'LX',
    iconBg: '#241B0A',
    iconColor: '#FFC97A',
    authorized: false,
  },
  {
    id: 'ltoken',
    name: 'LTOKEN 业务线',
    ou: 'ou=ltoken',
    role: '研发',
    iconText: 'LT',
    iconBg: '#1A2A1A',
    iconColor: '#A6FFB0',
    authorized: true,
  },
  {
    id: 'maas',
    name: 'MAAS 业务线',
    ou: 'ou=maas',
    role: '未授权',
    iconText: 'MA',
    iconBg: '#2A1A2A',
    iconColor: '#FFB0E0',
    authorized: false,
  },
]

function loadSavedBL(): BusinessLine | null {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) return JSON.parse(saved)
  } catch {
    // ignore
  }
  return null
}

export const useBusinessLineStore = defineStore('businessLine', () => {
  const businessLines = ref<BusinessLine[]>(MOCK_BUSINESS_LINES)
  const currentBL = ref<BusinessLine | null>(loadSavedBL())

  // 默认选中第一个已授权的业务线
  if (!currentBL.value) {
    currentBL.value = businessLines.value.find((bl) => bl.authorized) || null
  }

  const authorizedBLs = computed(() =>
    businessLines.value.filter((bl) => bl.authorized),
  )

  function switchBL(id: string) {
    const bl = businessLines.value.find((b) => b.id === id)
    if (!bl || !bl.authorized) return
    currentBL.value = bl
    localStorage.setItem(STORAGE_KEY, JSON.stringify(bl))
  }

  return {
    businessLines,
    currentBL,
    authorizedBLs,
    switchBL,
  }
})
