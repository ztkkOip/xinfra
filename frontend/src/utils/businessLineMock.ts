import { computed } from 'vue'
import { useBusinessLineStore } from '@/stores/businessLine'

export const BUSINESS_LINE_NAMES = ['kodo', 'linxi', 'xinfra', 'las'] as const
export type BusinessLineName = (typeof BUSINESS_LINE_NAMES)[number]

interface BusinessLineMockProfile {
  name: BusinessLineName
  clusters: number
  nodes: number
  physicalMachines: number
  virtualMachines: number
  cpuAllocated: number
  components: number
  mysql: number
  redis: number
  alertsP0: number
  alertsP1: number
  tasksRunning: number
  primaryZone: string
  secondaryZone: string
  servicePrefix: string
}

const profiles: Record<BusinessLineName, BusinessLineMockProfile> = {
  kodo: {
    name: 'kodo',
    clusters: 3,
    nodes: 128,
    physicalMachines: 186,
    virtualMachines: 512,
    cpuAllocated: 61,
    components: 214,
    mysql: 38,
    redis: 92,
    alertsP0: 1,
    alertsP1: 5,
    tasksRunning: 2,
    primaryZone: 'IDC-华北机房',
    secondaryZone: 'IDC-华东机房',
    servicePrefix: 'kodo',
  },
  linxi: {
    name: 'linxi',
    clusters: 2,
    nodes: 74,
    physicalMachines: 96,
    virtualMachines: 238,
    cpuAllocated: 47,
    components: 128,
    mysql: 21,
    redis: 46,
    alertsP0: 0,
    alertsP1: 2,
    tasksRunning: 1,
    primaryZone: 'IDC-华东机房',
    secondaryZone: '阿里云-华南',
    servicePrefix: 'linxi',
  },
  xinfra: {
    name: 'xinfra',
    clusters: 2,
    nodes: 52,
    physicalMachines: 68,
    virtualMachines: 156,
    cpuAllocated: 39,
    components: 84,
    mysql: 12,
    redis: 31,
    alertsP0: 0,
    alertsP1: 1,
    tasksRunning: 1,
    primaryZone: 'IDC-华北机房',
    secondaryZone: '香港 IDC',
    servicePrefix: 'xinfra',
  },
  las: {
    name: 'las',
    clusters: 1,
    nodes: 34,
    physicalMachines: 42,
    virtualMachines: 118,
    cpuAllocated: 31,
    components: 66,
    mysql: 8,
    redis: 18,
    alertsP0: 0,
    alertsP1: 3,
    tasksRunning: 0,
    primaryZone: '七牛-新加坡',
    secondaryZone: 'AWS-美国',
    servicePrefix: 'las',
  },
}

export function normalizeBusinessLineName(name?: string | null): BusinessLineName {
  const normalized = String(name || '').trim().toLowerCase()
  if (BUSINESS_LINE_NAMES.includes(normalized as BusinessLineName)) {
    return normalized as BusinessLineName
  }
  return 'kodo'
}

export function useBusinessLineMockProfile() {
  const businessLineStore = useBusinessLineStore()
  const currentName = computed(() => normalizeBusinessLineName(businessLineStore.current?.name))
  const profile = computed(() => profiles[currentName.value])

  return {
    currentName,
    profile,
  }
}
