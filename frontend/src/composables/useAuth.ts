import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

export function useAuth() {
  const router = useRouter()
  const authStore = useAuthStore()

  const login = async (username: string, password: string) => {
    const response = await authApi.login({ username, password })
    const { token, user } = response.data
    authStore.setAuth(token, user)
    return response
  }

  const logout = async () => {
    await authApi.logout()
    authStore.clearAuth()
    router.push('/login')
  }

  const checkAuth = () => {
    if (!authStore.isLoggedIn()) {
      router.push('/login')
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
