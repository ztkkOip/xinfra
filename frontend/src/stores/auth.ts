import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi } from '@/api/auth'
import { getToken, setToken, removeToken, getUser, setUser, removeUser } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<any>(getUser())
  const isAdmin = computed(() => user.value?.is_admin === true || decodeAdminClaim(token.value))

  function setAuth(newToken: string, newUser: any) {
    token.value = newToken
    user.value = newUser
    setToken(newToken)
    setUser(newUser)
  }

  function setSessionToken(newToken: string) {
    token.value = newToken
    setToken(newToken)
  }

  async function refreshUser() {
    const currentToken = token.value || getToken()
    if (!currentToken) {
      clearAuth()
      return null
    }
    setSessionToken(currentToken)
    const response = await authApi.getUserInfo()
    setAuth(currentToken, response.data)
    return response.data
  }

  function clearAuth() {
    token.value = null
    user.value = null
    removeToken()
    removeUser()
  }

  function isLoggedIn() {
    return !!token.value
  }

  return {
    token,
    user,
    isAdmin,
    setAuth,
    setSessionToken,
    refreshUser,
    clearAuth,
    isLoggedIn,
  }
})

function decodeAdminClaim(token: string | null): boolean {
	if (!token) return false
	try {
		const payload = JSON.parse(decodeBase64Url(token.split('.')[1] || ''))
		return payload.admin === true || payload.is_admin === true
	} catch {
		return false
	}
}

function decodeBase64Url(value: string): string {
	const normalized = value.replace(/-/g, '+').replace(/_/g, '/')
	const padded = normalized.padEnd(normalized.length + ((4 - (normalized.length % 4)) % 4), '=')
	return atob(padded)
}
