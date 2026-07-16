<template>
  <div class="settings-page app-scrollbar" :key="renderKey">
    <header class="settings-header">
      <div>
        <p class="settings-eyebrow">RES / PREFERENCES</p>
        <h1>{{ t('setting.setting_title') }}</h1>
        <p>{{ t('setting.setting_desc') }}</p>
      </div>
      <span class="autosave-badge"><NIcon><CheckmarkCircleOutline/></NIcon>{{ t('setting.saved_automatically') }}</span>
    </header>

    <NTabs type="line" animated class="settings-tabs">
      <NTabPane name="basic" :tab="t('setting.basic_setting')">
        <NForm :model="formValue" size="medium" label-placement="top" class="settings-form">
          <div class="settings-grid">
            <section class="settings-card">
              <div class="card-heading">
                <span class="card-heading__icon card-heading__icon--orange"><NIcon><FolderOpenOutline/></NIcon></span>
                <div><h2>{{ t('setting.storage_title') }}</h2><p>{{ t('setting.storage_desc') }}</p></div>
              </div>

              <NFormItem :label="t('setting.save_dir')" path="SaveDirectory">
                <div class="inline-field">
                  <NInput :value="formValue.SaveDirectory" :placeholder="t('setting.save_dir')" readonly/>
                  <NButton secondary type="primary" @click="selectDir">{{ t('common.select') }}</NButton>
                </div>
              </NFormItem>

              <div class="two-column-fields">
                <NFormItem :label="t('setting.filename_rules')" path="FilenameLen">
                  <div class="inline-field inline-field--compact">
                    <NInputNumber v-model:value="formValue.FilenameLen" :min="0" :max="9999" placeholder="0"/>
                    <NSwitch v-model:value="formValue.FilenameTime"/>
                    <NTooltip trigger="hover"><template #trigger><NIcon class="help-icon"><HelpCircleOutline/></NIcon></template>{{ t('setting.filename_rules_tip') }}</NTooltip>
                  </div>
                </NFormItem>
                <NFormItem :label="t('setting.quality')" path="Quality">
                  <div class="inline-field">
                    <NSelect v-model:value="formValue.Quality" :options="options"/>
                    <NTooltip trigger="hover"><template #trigger><NIcon class="help-icon"><HelpCircleOutline/></NIcon></template>{{ t('setting.quality_tip') }}</NTooltip>
                  </div>
                </NFormItem>
              </div>
            </section>

            <section class="settings-card">
              <div class="card-heading">
                <span class="card-heading__icon card-heading__icon--green"><NIcon><PulseOutline/></NIcon></span>
                <div><h2>{{ t('setting.behavior_title') }}</h2><p>{{ t('setting.behavior_desc') }}</p></div>
              </div>

              <div class="toggle-list">
                <label class="toggle-row">
                  <span><strong>{{ t('setting.auto_proxy') }}</strong><small>{{ t('setting.auto_proxy_tip') }}</small></span>
                  <NSwitch v-model:value="formValue.AutoProxy"/>
                </label>
                <label class="toggle-row">
                  <span><strong>{{ t('setting.full_intercept') }}</strong><small>{{ t('setting.full_intercept_tip') }}</small></span>
                  <NSwitch v-model:value="formValue.WxAction"/>
                </label>
                <label class="toggle-row">
                  <span><strong>{{ t('setting.insert_tail') }}</strong><small>{{ t('setting.insert_tail_tip') }}</small></span>
                  <NSwitch v-model:value="formValue.InsertTail"/>
                </label>
              </div>
            </section>

            <section class="settings-card settings-card--danger settings-card--wide">
              <div class="danger-copy">
                <span class="card-heading__icon card-heading__icon--red"><NIcon><WarningOutline/></NIcon></span>
                <div><h2>{{ t('setting.danger_title') }}</h2><p>{{ t('setting.danger_desc') }}</p></div>
              </div>
              <NPopconfirm @positive-click="resetHandle">
                <template #trigger><NButton secondary type="error">{{ t('index.start_err_positiveText') }}</NButton></template>
                {{ t('index.reset_app_tip') }}
              </NPopconfirm>
            </section>
          </div>
        </NForm>
      </NTabPane>

      <NTabPane name="advanced" :tab="t('setting.advanced_setting')">
        <NForm :model="formValue" size="medium" label-placement="top" class="settings-form">
          <div class="settings-grid">
            <section class="settings-card">
              <div class="card-heading">
                <span class="card-heading__icon card-heading__icon--blue"><NIcon><OptionsOutline/></NIcon></span>
                <div><h2>{{ t('setting.network_title') }}</h2><p>{{ t('setting.network_desc') }}</p></div>
              </div>

              <div class="two-column-fields">
                <NFormItem label="Host" path="Host" :validation-status="hostValidationFeedback === '' ? undefined : 'error'" :feedback="hostValidationFeedback">
                  <NInput v-model:value="formValue.Host" placeholder="127.0.0.1"/>
                </NFormItem>
                <NFormItem label="Port" path="Port" :validation-status="portValidationFeedback === '' ? undefined : 'error'" :feedback="portValidationFeedback">
                  <NInput v-model:value="formValue.Port" placeholder="8899"/>
                </NFormItem>
              </div>
              <NFormItem :label="t('setting.upstream_proxy')" path="UpstreamProxy">
                <div class="inline-field">
                  <NInput v-model:value="formValue.UpstreamProxy" placeholder="http://127.0.0.1:7890"/>
                  <NSwitch v-model:value="formValue.OpenProxy"/>
                </div>
                <template #feedback>{{ t('setting.upstream_proxy_tip') }}</template>
              </NFormItem>
              <label class="toggle-row toggle-row--standalone">
                <span><strong>{{ t('setting.download_proxy') }}</strong><small>{{ t('setting.download_proxy_tip') }}</small></span>
                <NSwitch v-model:value="formValue.DownloadProxy"/>
              </label>
            </section>

            <section class="settings-card">
              <div class="card-heading">
                <span class="card-heading__icon card-heading__icon--green"><NIcon><SpeedometerOutline/></NIcon></span>
                <div><h2>{{ t('setting.performance_title') }}</h2><p>{{ t('setting.performance_desc') }}</p></div>
              </div>
              <div class="two-column-fields">
                <NFormItem :label="t('setting.connections')" path="TaskNumber"><NInputNumber v-model:value="formValue.TaskNumber" :min="2" :max="64"/></NFormItem>
                <NFormItem :label="t('setting.down_number')" path="DownNumber"><NInputNumber v-model:value="formValue.DownNumber" :min="1" :max="10"/></NFormItem>
              </div>
              <NFormItem label="User-Agent" path="UserAgent"><NInput v-model:value="formValue.UserAgent" placeholder="User-Agent"/></NFormItem>
              <NFormItem label="Headers" path="Headers"><NInput v-model:value="formValue.UseHeaders" placeholder="User-Agent, Referer, Authorization, Cookie"/></NFormItem>
            </section>

            <section class="settings-card settings-card--wide">
              <div class="card-heading">
                <span class="card-heading__icon card-heading__icon--orange"><NIcon><ShieldCheckmarkOutline/></NIcon></span>
                <div><h2>{{ t('setting.rules_title') }}</h2><p>{{ t('setting.rules_desc') }}</p></div>
              </div>
              <div class="rule-grid">
                <NFormItem :label="t('setting.domain_rule')" path="DomainRule">
                  <NInput v-model:value="formValue.Rule" type="textarea" :rows="10" :placeholder="t('setting.domain_rule_tip')" class="code-input"/>
                </NFormItem>
                <NFormItem :label="t('setting.mime_map')" path="MimeMap" :validation-status="mimeValidationFeedback ? 'error' : undefined" :feedback="mimeValidationFeedback">
                  <NInput v-model:value="MimeMap" type="textarea" :rows="10" placeholder='{"video/mp4": { "Type": "video", "Suffix": ".mp4"}}' class="code-input"/>
                </NFormItem>
              </div>
            </section>
          </div>
        </NForm>
      </NTabPane>
    </NTabs>
  </div>
</template>

<script lang="ts" setup>
import {
  HelpCircleOutline,
  FolderOpenOutline,
  PulseOutline,
  OptionsOutline,
  SpeedometerOutline,
  ShieldCheckmarkOutline,
  WarningOutline,
  CheckmarkCircleOutline
} from "@vicons/ionicons5"
import {ref, watch} from "vue"
import {useIndexStore} from "@/stores"
import type {appType} from "@/types/app"
import appApi from "@/api/app"
import {computed} from "vue"
import {useI18n} from 'vue-i18n'
import {isValidHost, isValidPort} from '@/func'
import {NButton, NIcon} from "naive-ui"
import * as bind from "../../wailsjs/go/core/Bind"

const {t} = useI18n()
const store = useIndexStore()

const options = computed(() =>
    t("setting.quality_value")
        .split(",")
        .map((value, index) => ({ value: index, label: value }))
)

const formValue = ref<appType.Config>(Object.assign({}, store.globalConfig))

const MimeMap = ref(formValue.value.MimeMap ? JSON.stringify(formValue.value.MimeMap, null, 2) : "")
const renderKey = ref(999)

const hostValidationFeedback = ref("")
const portValidationFeedback = ref("")
const mimeValidationFeedback = ref("")

watch(formValue.value, () => {
  formValue.value.Port = formValue.value.Port.trim()
  formValue.value.Host = formValue.value.Host.trim()

  if (!isValidHost(formValue.value.Host)) {
    hostValidationFeedback.value = t("setting.host_format_error")
    return
  } else {
    hostValidationFeedback.value = ''
  }

  if (!isValidPort(parseInt(formValue.value.Port))) {
    portValidationFeedback.value = t("setting.port_format_error")
    return
  } else {
    portValidationFeedback.value = ''
  }
  store.setConfig(formValue.value)
}, {deep: true})

watch(MimeMap, () => {
  try {
    const mimeMap = JSON.parse(MimeMap.value)
    mimeValidationFeedback.value = ""
    store.setConfig({MimeMap: mimeMap})
  } catch (e) {
    mimeValidationFeedback.value = t("setting.mime_format_error")
  }
})

watch(() => {
  return store.globalConfig.Theme
}, () => {
  formValue.value.Theme = store.globalConfig.Theme
})

watch(() => store.globalConfig.Locale, () => {
  formValue.value.Locale = store.globalConfig.Locale
  renderKey.value++
})

const selectDir = () => {
  appApi.openDirectoryDialog().then((res: any) => {
    if (res.code === 1) {
      formValue.value.SaveDirectory = res.data.folder
    }
  }).catch((err: any) => {
    window?.$message?.error(err)
  })
}

const resetHandle = ()=>{
  localStorage.clear()
  bind.ResetApp()
}
</script>
<style scoped>
.settings-page {
  height: 100%;
  overflow-y: auto;
  padding: 25px 28px 36px;
  color: var(--app-text);
}

.settings-header {
  display: flex;
  max-width: 1120px;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin: 0 auto 22px;
  animation: settings-rise 0.4s ease-out both;
}

.settings-eyebrow {
  margin-bottom: 5px;
  color: var(--app-accent) !important;
  font-family: 'SFMono-Regular', 'Cascadia Code', monospace;
  font-size: 9px !important;
  font-weight: 750;
  letter-spacing: 0.18em;
}

.settings-header h1 {
  color: var(--app-text);
  font-size: clamp(26px, 2.6vw, 34px);
  font-weight: 750;
  letter-spacing: -0.04em;
  line-height: 1.15;
}

.settings-header p {
  margin-top: 6px;
  color: var(--app-muted);
  font-size: 12px;
}

.autosave-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 3px;
  padding: 6px 10px;
  border: 1px solid rgba(35, 131, 95, 0.18);
  border-radius: 999px;
  color: var(--app-success);
  background: rgba(35, 131, 95, 0.07);
  font-size: 10px;
  font-weight: 650;
  white-space: nowrap;
}

.settings-tabs {
  max-width: 1120px;
  margin: 0 auto;
  animation: settings-rise 0.48s 0.05s ease-out both;
}

.settings-tabs :deep(.n-tabs-nav) {
  position: sticky;
  z-index: 10;
  top: -25px;
  padding-top: 3px;
  background: color-mix(in srgb, var(--app-bg) 92%, transparent);
  backdrop-filter: blur(14px);
}

.settings-tabs :deep(.n-tabs-tab) {
  padding: 10px 4px 12px;
  font-size: 12px;
}

.settings-tabs :deep(.n-tab-pane) { padding-top: 18px; }
.settings-form { --wails-draggable: no-drag; }

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.settings-card {
  min-width: 0;
  padding: 19px;
  border: 1px solid var(--app-border);
  border-radius: 17px;
  background: var(--app-surface);
  box-shadow: 0 10px 28px rgba(28, 42, 34, 0.045);
}

.settings-card--wide { grid-column: 1 / -1; }

.card-heading {
  display: flex;
  align-items: flex-start;
  gap: 11px;
  margin-bottom: 18px;
}

.card-heading__icon {
  display: grid;
  width: 34px;
  height: 34px;
  flex: none;
  place-items: center;
  border-radius: 10px;
  font-size: 17px;
}

.card-heading__icon--orange { color: #c85b24; background: rgba(220, 107, 47, 0.11); }
.card-heading__icon--green { color: #23835f; background: rgba(35, 131, 95, 0.11); }
.card-heading__icon--blue { color: #477d9f; background: rgba(71, 125, 159, 0.11); }
.card-heading__icon--red { color: #c64c46; background: rgba(211, 84, 77, 0.11); }

.card-heading h2,
.danger-copy h2 {
  color: var(--app-text);
  font-size: 13px;
  font-weight: 700;
  line-height: 1.3;
}

.card-heading p,
.danger-copy p {
  margin-top: 3px;
  color: var(--app-muted);
  font-size: 10px;
  line-height: 1.5;
}

.settings-card :deep(.n-form-item-label) {
  color: var(--app-muted);
  font-size: 10px;
  font-weight: 650;
  letter-spacing: 0.025em;
}

.settings-card :deep(.n-form-item:last-child) { margin-bottom: 0; }

.inline-field {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 8px;
}

.inline-field > :deep(.n-input),
.inline-field > :deep(.n-select) { flex: 1; }
.inline-field--compact > :deep(.n-input-number) { flex: 1; }

.two-column-fields {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.help-icon {
  flex: none;
  color: var(--app-muted);
  cursor: help;
  font-size: 16px;
}

.toggle-list {
  display: flex;
  flex-direction: column;
}

.toggle-row {
  display: flex;
  min-height: 57px;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 10px 2px;
  border-bottom: 1px solid var(--app-border);
  cursor: pointer;
}

.toggle-row:last-child { border-bottom: 0; }
.toggle-row--standalone { margin-top: 1px; padding-inline: 0; border-bottom: 0; }

.toggle-row > span {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.toggle-row strong {
  color: var(--app-text);
  font-size: 12px;
  font-weight: 650;
}

.toggle-row small {
  margin-top: 3px;
  overflow: hidden;
  color: var(--app-muted);
  font-size: 9px;
  line-height: 1.4;
  text-overflow: ellipsis;
}

.settings-card--danger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  border-color: rgba(211, 84, 77, 0.18);
  background:
    linear-gradient(110deg, rgba(211, 84, 77, 0.055), transparent 45%),
    var(--app-surface);
}

.danger-copy { display: flex; align-items: center; gap: 11px; }
.rule-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
.code-input :deep(textarea) { font-family: 'SFMono-Regular', 'Cascadia Code', monospace; font-size: 11px; line-height: 1.65; }

@keyframes settings-rise {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 1050px) {
  .settings-page { padding-inline: 20px; }
  .settings-grid { grid-template-columns: 1fr; }
  .settings-card--wide { grid-column: auto; }
}

@media (max-width: 760px) {
  .two-column-fields,
  .rule-grid { grid-template-columns: 1fr; }
  .settings-card--danger { align-items: flex-start; flex-direction: column; }
}

@media (prefers-reduced-motion: reduce) {
  .settings-header,
  .settings-tabs { animation: none; }
}
</style>
