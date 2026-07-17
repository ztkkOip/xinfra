import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { redirectToSSO } from '@/utils/sso'

export function useAuth() {
  const authStore = useAuthStore()

  const login = async (username: string, password: string) => {
    const response = await authApi.login({ username, password })
    const { token, user } = response.data
    authStore.setAuth(token, user)
    return response
  }

  const logout = async () => {
    let logoutUrl = ''
    try {
      const response = await authApi.logout()
      logoutUrl = response.data.logout_url || ''
    } finally {
      authStore.clearAuth()
      window.location.assign(logoutUrl || '/login?logged_out=1')
    }
  }

  const checkAuth = () => {
    if (!authStore.isLoggedIn()) {
      redirectToSSO()
      return false
    }
    return true
  }

  return {
    login,
    logout,
    checkAuth,
  }
}
