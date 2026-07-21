import { getToken } from '@/utils/auth'

export interface DeploymentCreatePayload {
  component: string
  business_line_id: number
  params: Record<string, unknown>
}

export interface DeploymentCreateResult {
  deployment_id: string
  status: string
}

export const deploymentApi = {
  async create(payload: DeploymentCreatePayload): Promise<DeploymentCreateResult> {
    const data = await authRequest('/auth/api/v1/deployments', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    return {
      deployment_id: String(data.deployment_id || ''),
      status: String(data.status || ''),
    }
  },

  async cancel(deploymentId: string): Promise<void> {
    await authRequest(`/auth/api/v1/deployments/${encodeURIComponent(deploymentId)}/cancel`, {
      method: 'POST',
    })
  },

  eventsURL(deploymentId: string): string {
    const token = getToken()
    const params = token ? `?access_token=${encodeURIComponent(token)}` : ''
    return `/auth/api/v1/deployments/${encodeURIComponent(deploymentId)}/events${params}`
  },
}

async function authRequest(path: string, init: RequestInit = {}) {
  const token = getToken()
  const response = await fetch(path, {
    ...init,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init.headers,
    },
  })
  const text = await response.text()
  const data = parseResponseBody(text)
  if (!response.ok) {
    const message = data?.error || data?.message || text || `HTTP ${response.status}`
    throw new Error(message)
  }
  return data || {}
}

function parseResponseBody(text: string) {
  if (!text.trim()) {
    return {}
  }
  try {
    return JSON.parse(text)
  } catch {
    return { error: text }
  }
}
