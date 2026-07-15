<script setup>
import { computed, onMounted, ref } from 'vue'

const TOKEN_KEY = 'authserver_token'
const SSO_LOGIN_PATH = '/auth/api/v1/login/internal-sso'
const WAYEN_LOGIN_PATH = '/auth/api/v1/wayen/login'
const CLOUDDM_LOGIN_PATH = '/auth/api/v1/clouddm/login'

const token = ref(localStorage.getItem(TOKEN_KEY) || '')
const currentUser = ref(null)
const message = ref('')
const loading = ref(false)
const redirectingToSSO = ref(false)
const wayenLoading = ref(false)
const clouddmLoading = ref(false)

const isAuthed = computed(() => Boolean(token.value))
const currentUserLabel = computed(() => {
  if (!currentUser.value) {
    return ''
  }
  return currentUser.value.display_name || currentUser.value.username || ''
})

function clearAuth() {
  token.value = ''
  currentUser.value = null
  localStorage.removeItem(TOKEN_KEY)
}

function relayState(openApp = '') {
  const url = new URL(window.location.href)
  url.searchParams.delete('open_wayen')
  url.searchParams.delete('open_app')
  if (openApp === 'wayen') {
    url.searchParams.set('open_wayen', '1')
  } else if (openApp) {
    url.searchParams.set('open_app', openApp)
  }
  return `${url.pathname}${url.search}${url.hash}` || '/'
}

function ssoLogin({ openApp = '' } = {}) {
  if (redirectingToSSO.value) {
    return
  }
  redirectingToSSO.value = true
  window.location.assign(`${SSO_LOGIN_PATH}?relay_state=${encodeURIComponent(relayState(openApp))}`)
}

async function api(path, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  }
  if (token.value) {
    headers.Authorization = `Bearer ${token.value}`
  }
  const response = await fetch(path, { ...options, headers })
  const data = await response.json().catch(() => ({}))
  if (!response.ok) {
    if (response.status === 401) {
      clearAuth()
      ssoLogin()
    }
    const error = new Error(data.error || `HTTP ${response.status}`)
    error.status = response.status
    throw error
  }
  return data
}

async function run(task, successText = '操作完成') {
  loading.value = true
  message.value = ''
  try {
    const result = await task()
    message.value = successText
    return result
  } catch (error) {
    message.value = error.message
    return null
  } finally {
    loading.value = false
  }
}

function logout() {
  clearAuth()
  ssoLogin()
}

async function loadMe() {
  currentUser.value = await api('/auth/api/v1/users/me')
}

async function openWayen() {
  wayenLoading.value = true
  message.value = ''
  try {
    const response = await fetch(WAYEN_LOGIN_PATH, {
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      credentials: 'include',
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      if (response.status === 401) {
        clearAuth()
        ssoLogin({ openApp: 'wayen' })
      }
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    if (!data.target_url) {
      throw new Error('Wayen 跳转地址为空')
    }
    window.location.assign(data.target_url)
  } catch (error) {
    message.value = error.message
  } finally {
    wayenLoading.value = false
  }
}

async function openCloudDM() {
  clouddmLoading.value = true
  message.value = ''
  try {
    const response = await fetch(CLOUDDM_LOGIN_PATH, {
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      credentials: 'include',
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      if (response.status === 401) {
        clearAuth()
        ssoLogin({ openApp: 'clouddm' })
      }
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    if (!data.target_url) {
      throw new Error('CloudDM 跳转地址为空')
    }
    window.location.assign(data.target_url)
  } catch (error) {
    message.value = error.message
  } finally {
    clouddmLoading.value = false
  }
}

onMounted(async () => {
  const url = new URL(window.location.href)
  const ssoToken = url.searchParams.get('sso_token')
  const openApp = url.searchParams.get('open_wayen') === '1' ? 'wayen' : url.searchParams.get('open_app')
  if (ssoToken) {
    token.value = ssoToken
    localStorage.setItem(TOKEN_KEY, ssoToken)
    url.searchParams.delete('sso_token')
    url.searchParams.delete('open_wayen')
    url.searchParams.delete('open_app')
    window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`)
  }
  if (!token.value) {
    ssoLogin()
    return
  }
  await run(async () => {
    await loadMe()
  }, '已就绪')
  if (openApp === 'wayen') {
    await openWayen()
  } else if (openApp === 'clouddm') {
    await openCloudDM()
  }
})
</script>

<template>
  <main v-if="!isAuthed" class="login-page">
    <section class="login-panel">
      <div>
        <p class="eyebrow">AuthServer</p>
        <h1>应用入口</h1>
      </div>
      <p class="message">正在跳转到 SSO 登录...</p>
    </section>
  </main>

  <main v-else class="app-page">
    <header class="app-header">
      <div>
        <p class="eyebrow">AuthServer</p>
        <h1>应用入口</h1>
      </div>
      <div class="user-block">
        <span>{{ currentUserLabel }}</span>
        <button class="secondary" type="button" @click="logout">退出</button>
      </div>
    </header>

    <section class="app-main">
      <button class="app-button wayen-button" type="button" :disabled="wayenLoading" @click="openWayen">
        <span class="app-icon wayen-icon" aria-hidden="true">W</span>
        <span class="app-title">Wayen</span>
      </button>
      <button class="app-button clouddm-button" type="button" :disabled="clouddmLoading" @click="openCloudDM">
        <span class="app-icon clouddm-icon" aria-hidden="true">C</span>
        <span class="app-title">CloudDM</span>
      </button>
      <p v-if="message" class="message">{{ message }}</p>
    </section>
  </main>
</template>
