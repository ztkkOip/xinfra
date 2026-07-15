import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { redirectToSSO } from '@/utils/sso'

const SSO_SIGNOUT_URL = 'https://bo-staging-sso-internal.jfcs-k8s-qa1.qiniu.io/signout'

export function useAuth() {
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
    window.location.assign(SSO_SIGNOUT_URL)
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
