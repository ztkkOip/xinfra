import { ref, onMounted } from 'vue'

export type Theme = 'light' | 'dark'

const THEME_KEY = 'xinfra-theme'

// 默认浅色模式
const theme = ref<Theme>('light')

export function useTheme() {
  // 初始化主题
  const initTheme = () => {
    const savedTheme = localStorage.getItem(THEME_KEY) as Theme | null
    if (savedTheme && ['light', 'dark'].includes(savedTheme)) {
      theme.value = savedTheme
    } else {
      // 默认浅色模式
      theme.value = 'light'
    }
    applyTheme(theme.value)
  }

  // 应用主题
  const applyTheme = (newTheme: Theme) => {
    document.documentElement.setAttribute('data-theme', newTheme)
    document.documentElement.classList.toggle('dark', newTheme === 'dark')
  }

  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem(THEME_KEY, theme.value)
    applyTheme(theme.value)
  }

  // 设置主题
  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme
    localStorage.setItem(THEME_KEY, theme.value)
    applyTheme(theme.value)
  }

  // 监听系统主题变化（可选）
  const watchSystemTheme = () => {
    if (window.matchMedia) {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      mediaQuery.addEventListener('change', (e) => {
        // 只有在用户没有手动设置过主题时才跟随系统
        if (!localStorage.getItem(THEME_KEY)) {
          theme.value = e.matches ? 'dark' : 'light'
          applyTheme(theme.value)
        }
      })
    }
  }

  onMounted(() => {
    initTheme()
    watchSystemTheme()
  })

  return {
    theme,
    toggleTheme,
    setTheme,
    initTheme
  }
}