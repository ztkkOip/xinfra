import { setToken } from '@/utils/auth'

const SSO_LOGIN_PATH = '/auth/api/v1/login/internal-sso'

export function relayState(openApp = ''): string {
  const url = new URL(window.location.href)
  url.searchParams.delete('sso_token')
  url.searchParams.delete('open_app')
  if (openApp) {
    url.searchParams.set('open_app', openApp)
  }
  return `${url.pathname}${url.search}${url.hash}` || '/'
}

export function redirectToSSO(openApp = ''): void {
  window.location.assign(`${SSO_LOGIN_PATH}?relay_state=${encodeURIComponent(relayState(openApp))}`)
}

export function consumeSSOToken(): string {
  const url = new URL(window.location.href)
  const token = url.searchParams.get('sso_token') || ''
  if (!token) {
    return ''
  }

  setToken(token)
  url.searchParams.delete('sso_token')
  window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`)
  return token
}

export function consumeOpenApp(): string {
  const url = new URL(window.location.href)
  const openApp = url.searchParams.get('open_app') || ''
  if (!openApp) {
    return ''
  }

  url.searchParams.delete('open_app')
  window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`)
  return openApp
}
