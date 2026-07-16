<template>
  <NConfigProvider class="h-full" :theme="theme" :theme-overrides="themeOverrides" :locale="uiLocale">
    <NaiveProvider>
      <RouterView/>
    </NaiveProvider>
    <NGlobalStyle/>
    <NModalProvider/>
  </NConfigProvider>
</template>

<script setup lang="ts">
import NaiveProvider from '@/components/NaiveProvider.vue'
import {darkTheme, lightTheme, zhCN, enUS} from 'naive-ui'
import {useIndexStore} from "@/stores"
import {computed, onMounted} from "vue"
import {useEventStore} from "@/stores/event"
import type {appType} from "@/types/app"
import type {GlobalThemeOverrides} from "naive-ui"
import {useI18n} from 'vue-i18n'

const store = useIndexStore()
const eventStore = useEventStore()
const {locale} = useI18n()

const theme = computed(() => {
  if (store.globalConfig.Theme === "darkTheme") {
    document.documentElement.classList.add('dark');
    return darkTheme
  }
  document.documentElement.classList.remove('dark');
  return lightTheme
})

const uiLocale = computed(() => {
  locale.value = store.globalConfig.Locale
  if (store.globalConfig.Locale === "zh") {
    return zhCN
  }
  return enUS
})

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const isDark = store.globalConfig.Theme === "darkTheme"
  return {
    common: {
      primaryColor: "#dc6b2f",
      primaryColorHover: "#ef7c3d",
      primaryColorPressed: "#bd5420",
      primaryColorSuppl: "#ef7c3d",
      successColor: "#23835f",
      warningColor: "#d18a27",
      errorColor: "#d3544d",
      infoColor: "#477d9f",
      borderRadius: "10px",
      borderRadiusSmall: "7px",
      fontFamily: "'Avenir Next', 'PingFang SC', 'Microsoft YaHei', sans-serif",
      fontFamilyMono: "'SFMono-Regular', 'Cascadia Code', 'Roboto Mono', monospace",
      bodyColor: isDark ? "#101614" : "#f2f3ee",
      cardColor: isDark ? "#161e1b" : "#fffefa",
      modalColor: isDark ? "#161e1b" : "#fffefa",
      popoverColor: isDark ? "#1a2420" : "#fffefa",
      tableColor: isDark ? "#161e1b" : "#fffefa",
      tableHeaderColor: isDark ? "#1b2521" : "#f5f4ee",
      hoverColor: isDark ? "rgba(255,255,255,.055)" : "rgba(22,41,33,.045)",
      textColorBase: isDark ? "#e8eee9" : "#1b2923",
      textColor1: isDark ? "#e8eee9" : "#1b2923",
      textColor2: isDark ? "#afbbb4" : "#52625a",
      textColor3: isDark ? "#7d8d84" : "#7a8881",
      borderColor: isDark ? "#2a3731" : "#dedfd7",
      dividerColor: isDark ? "#27332e" : "#e4e4dc"
    },
    Button: {
      fontWeight: "600",
      heightMedium: "38px",
      paddingMedium: "0 16px"
    },
    DataTable: {
      thFontWeight: "650",
      thColor: isDark ? "#1b2521" : "#f5f4ee",
      thColorHover: isDark ? "#202c27" : "#efeee7",
      tdColor: isDark ? "#161e1b" : "#fffefa",
      tdColorHover: isDark ? "#1a2420" : "#faf9f4",
      borderColor: isDark ? "#26332d" : "#e8e7df"
    },
    Input: {borderRadius: "9px"},
    Select: {
      peers: {InternalSelection: {borderRadius: "9px"}}
    },
    Tabs: {tabFontWeightActive: "650"}
  }
})

onMounted(async () => {
  await store.init()

  eventStore.init()
  eventStore.addHandle({
    type: "message",
    event: (res: appType.Message) => {
      switch (res?.code) {
        case 0:
          window.$message?.error(res.message)
          break
        case 1:
          window.$message?.success(res.message)
          break
      }
    }
  })
})
</script>
