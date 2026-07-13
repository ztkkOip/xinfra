<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

const TOKEN_KEY = 'authserver_token'
const SSO_LOGIN_PATH = '/auth/api/v1/login/internal-sso'
const WAYEN_LOGIN_PATH = '/auth/api/v1/wayen/login'
const WAYEN_CREDENTIAL_PATH = '/auth/api/v1/wayen/credential'

const token = ref(localStorage.getItem(TOKEN_KEY) || '')
const currentUser = ref(null)
const message = ref('')
const loading = ref(false)
const ssoLoading = ref(false)
const wayenLoading = ref(false)
const wayenCredentialLoading = ref(false)
const wayenConfigured = ref(false)

const wayenForm = reactive({
  email: '',
  password: '',
})

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

function ssoLogin() {
  const relayState = `${window.location.pathname}${window.location.search}${window.location.hash}` || '/'
  ssoLoading.value = true
  window.location.assign(`${SSO_LOGIN_PATH}?relay_state=${encodeURIComponent(relayState)}`)
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
  message.value = '已退出'
}

async function loadMe() {
  currentUser.value = await api('/auth/api/v1/users/me')
}

async function loadWayenCredential() {
  const result = await api(WAYEN_CREDENTIAL_PATH)
  wayenForm.email = result.email || (currentUser.value && currentUser.value.email) || ''
  wayenConfigured.value = Boolean(result.configured)
}

async function saveWayenCredential() {
  wayenCredentialLoading.value = true
  message.value = ''
  try {
    const result = await api(WAYEN_CREDENTIAL_PATH, {
      method: 'PUT',
      body: JSON.stringify(wayenForm),
    })
    wayenForm.email = result.email || wayenForm.email
    wayenForm.password = ''
    wayenConfigured.value = Boolean(result.configured)
    if (currentUser.value) {
      currentUser.value.email = wayenForm.email
    }
    message.value = 'Wayen 凭据已保存'
  } catch (error) {
    message.value = error.message
  } finally {
    wayenCredentialLoading.value = false
  }
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

onMounted(async () => {
  const url = new URL(window.location.href)
  const ssoToken = url.searchParams.get('sso_token')
  if (ssoToken) {
    token.value = ssoToken
    localStorage.setItem(TOKEN_KEY, ssoToken)
    url.searchParams.delete('sso_token')
    window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`)
  }
  if (!token.value) {
    return
  }
  await run(async () => {
    await loadMe()
    await loadWayenCredential()
  }, '已就绪')
})
</script>

<template>
  <main v-if="!isAuthed" class="login-page">
    <section class="login-panel">
      <div>
        <p class="eyebrow">AuthServer</p>
        <h1>Wayen 登录入口</h1>
      </div>
      <p class="message">当前仅保留 SSO 登录和 Wayen 单点跳转。</p>
      <button class="primary" type="button" :disabled="ssoLoading" @click="ssoLogin">SSO 登录</button>
      <p v-if="message" class="message">{{ message }}</p>
    </section>
  </main>

  <main v-else class="wayen-page">
    <header class="wayen-header">
      <div>
        <p class="eyebrow">AuthServer</p>
        <h1>Wayen</h1>
      </div>
      <div class="user-block">
        <span>{{ currentUserLabel }}</span>
        <button class="secondary" type="button" @click="logout">退出</button>
      </div>
    </header>

    <section class="wayen-main">
      <button class="wayen-button" type="button" :disabled="wayenLoading" @click="openWayen">
        <span class="wayen-icon" aria-hidden="true">W</span>
        <span class="wayen-title">Wayen</span>
      </button>
      <form class="wayen-config form stack" @submit.prevent="saveWayenCredential">
        <p class="config-title">Wayen 凭据</p>
        <label>
          <span>Wayen 账号</span>
          <input v-model="wayenForm.email" autocomplete="email" />
        </label>
        <label>
          <span>Wayen 密码</span>
          <input v-model="wayenForm.password" type="password" autocomplete="current-password" />
        </label>
        <button class="secondary" type="submit" :disabled="wayenCredentialLoading">
          {{ wayenConfigured ? '更新凭据' : '保存凭据' }}
        </button>
      </form>
      <p v-if="message" class="message">{{ message }}</p>
    </section>
  </main>
</template>
