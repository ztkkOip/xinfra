import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/auth'
import { getToken, setToken, removeToken, getUser, setUser, removeUser } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<any>(getUser())

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
    setAuth,
    setSessionToken,
    refreshUser,
    clearAuth,
    isLoggedIn,
  }
})
