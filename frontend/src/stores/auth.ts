import { defineStore } from 'pinia'
import { ref } from 'vue'
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
    clearAuth,
    isLoggedIn,
  }
})
