<template>
  <aside class="sidebar" :class="{'sidebar--darwin': envInfo.platform === 'darwin'}">
    <Screen v-if="envInfo.platform !== 'darwin'" class="window-controls"/>

    <button class="brand" type="button" @click="handleUtility('github')" :title="store.appInfo.AppName || 'res-downloader'">
      <span class="brand-mark">
        <img src="@/assets/image/logo.png" alt="res-downloader"/>
      </span>
      <span v-if="showUpdate" class="update-dot">New</span>
    </button>

    <nav class="primary-nav" :aria-label="t('menu.navigation')">
      <NTooltip v-for="item in navigation" :key="item.key" placement="right" trigger="hover">
        <template #trigger>
          <button
              class="nav-item"
              :class="{'nav-item--active': menuValue === item.key}"
              type="button"
              @click="handleNavigate(item.key)"
          >
            <NIcon size="22"><component :is="item.icon"/></NIcon>
            <span class="nav-label">{{ item.label }}</span>
          </button>
        </template>
        {{ item.label }}
      </NTooltip>
    </nav>

    <div class="sidebar-spacer"/>

    <div class="capture-indicator" :class="{'capture-indicator--active': store.isProxy}">
      <span class="capture-indicator__dot"/>
      <span class="capture-indicator__text">{{ store.isProxy ? t('index.capture_active_short') : t('index.capture_idle_short') }}</span>
    </div>

    <div class="utility-nav">
      <NTooltip v-for="item in utilities" :key="item.key" placement="right" trigger="hover">
        <template #trigger>
          <button class="utility-button" type="button" @click="handleUtility(item.key)">
            <NIcon size="19"><component :is="item.icon"/></NIcon>
          </button>
        </template>
        {{ item.label }}
      </NTooltip>
    </div>
  </aside>

  <Footer v-model:showModal="showAppInfo"/>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref, watch} from "vue"
import {useRoute, useRouter} from "vue-router"
import {
  CloudOutline,
  SettingsOutline,
  HelpCircleOutline,
  MoonOutline,
  SunnyOutline,
  LanguageSharp,
  LogoGithub
} from "@vicons/ionicons5"
import {NIcon, NTooltip} from "naive-ui"
import {useIndexStore} from "@/stores"
import Footer from "@/components/Footer.vue"
import Screen from "@/components/Screen.vue"
import {BrowserOpenURL} from "../../../wailsjs/runtime"
import {useI18n} from "vue-i18n"
import request from "@/api/request"
import {compareVersions} from "@/func"

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const store = useIndexStore()
const envInfo = store.envInfo
const showAppInfo = ref(false)
const showUpdate = ref(false)
const menuValue = ref(route.fullPath.substring(1))

const navigation = computed(() => [
  {label: t("menu.index"), key: "index", icon: CloudOutline},
  {label: t("menu.setting"), key: "setting", icon: SettingsOutline}
])

const utilities = computed(() => [
  {label: "GitHub", key: "github", icon: LogoGithub},
  {label: t("menu.locale"), key: "locale", icon: LanguageSharp},
  {
    label: t("menu.theme"),
    key: "theme",
    icon: store.globalConfig.Theme === "darkTheme" ? SunnyOutline : MoonOutline
  },
  {label: t("menu.about"), key: "about", icon: HelpCircleOutline}
])

watch(() => route.path, () => {
  menuValue.value = route.fullPath.substring(1)
})

onMounted(() => {
  request({
    url: "https://res.putyy.com/version.json?v=" + Date.now(),
    method: "get"
  }).then((res) => {
    showUpdate.value = compareVersions(res.version, store.appInfo.Version) === 1
  }).catch(() => {
    showUpdate.value = false
  })
})

const handleNavigate = (key: string) => {
  menuValue.value = key
  router.push({path: "/" + key})
}

const handleUtility = (key: string) => {
  if (key === "about") {
    showAppInfo.value = true
  } else if (key === "github") {
    BrowserOpenURL("https://github.com/putyy/res-downloader")
  } else if (key === "theme") {
    store.setConfig({Theme: store.globalConfig.Theme === "darkTheme" ? "lightTheme" : "darkTheme"})
  } else if (key === "locale") {
    store.setConfig({Locale: store.globalConfig.Locale === "zh" ? "en" : "zh"})
  }
}
</script>

<style scoped>
.sidebar {
  position: relative;
  z-index: 50;
  display: flex;
  width: 92px;
  min-width: 92px;
  height: 100%;
  flex-direction: column;
  align-items: center;
  padding: 14px 12px 16px;
  color: #dce7e0;
  background:
    linear-gradient(155deg, rgba(255, 255, 255, 0.035), transparent 42%),
    var(--app-sidebar);
  box-shadow: 10px 0 28px rgba(12, 21, 17, 0.09);
  --wails-draggable: drag;
}

.sidebar--darwin {
  padding-top: 30px;
}

.window-controls {
  width: 100%;
  margin-bottom: 12px;
  --wails-draggable: no-drag;
}

.brand {
  position: relative;
  display: grid;
  width: 56px;
  height: 56px;
  place-items: center;
  border: 0;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.075);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08);
  cursor: pointer;
  --wails-draggable: no-drag;
}

.brand-mark img {
  display: block;
  width: 42px;
  height: 42px;
  object-fit: contain;
  transition: transform 0.25s ease;
}

.brand:hover img {
  transform: rotate(-5deg) scale(1.05);
}

.update-dot {
  position: absolute;
  top: -5px;
  right: -10px;
  padding: 2px 5px;
  border: 2px solid var(--app-sidebar);
  border-radius: 999px;
  color: #fff;
  background: var(--app-danger);
  font-size: 8px;
  font-weight: 800;
  letter-spacing: 0.04em;
}

.primary-nav {
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: 9px;
  margin-top: 34px;
}

.nav-item {
  position: relative;
  display: flex;
  width: 100%;
  height: 58px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  border: 0;
  border-radius: 16px;
  color: #8fa198;
  background: transparent;
  cursor: pointer;
  transition: color 0.2s ease, background 0.2s ease, transform 0.2s ease;
  --wails-draggable: no-drag;
}

.nav-item:hover {
  color: #f0f5f2;
  background: rgba(255, 255, 255, 0.065);
  transform: translateY(-1px);
}

.nav-item--active {
  color: #fff6ef;
  background: rgba(220, 107, 47, 0.2);
  box-shadow: inset 0 0 0 1px rgba(235, 132, 78, 0.18);
}

.nav-item--active::before {
  position: absolute;
  left: -12px;
  width: 3px;
  height: 26px;
  border-radius: 0 4px 4px 0;
  background: #ed8148;
  content: "";
}

.nav-label {
  font-size: 10px;
  font-weight: 650;
  letter-spacing: 0.04em;
}

.sidebar-spacer {
  flex: 1;
}

.capture-indicator {
  display: flex;
  width: 100%;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  margin-bottom: 16px;
  color: #71837a;
}

.capture-indicator__dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #607169;
  box-shadow: 0 0 0 4px rgba(96, 113, 105, 0.12);
}

.capture-indicator--active {
  color: #6fceaa;
}

.capture-indicator--active .capture-indicator__dot {
  background: #49bd90;
  box-shadow: 0 0 0 4px rgba(73, 189, 144, 0.13), 0 0 12px rgba(73, 189, 144, 0.58);
  animation: capture-pulse 1.8s ease-in-out infinite;
}

.capture-indicator__text {
  font-size: 8px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.utility-nav {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 7px;
}

.utility-button {
  display: grid;
  width: 30px;
  height: 30px;
  place-items: center;
  border: 0;
  border-radius: 9px;
  color: #7f9188;
  background: transparent;
  cursor: pointer;
  transition: color 0.2s ease, background 0.2s ease;
  --wails-draggable: no-drag;
}

.utility-button:hover {
  color: #e9f0ec;
  background: rgba(255, 255, 255, 0.08);
}

@keyframes capture-pulse {
  50% { opacity: 0.62; transform: scale(0.88); }
}
</style>
