<template>
  <div class="row-actions">
    <NTooltip trigger="hover">
      <template #trigger>
        <button class="row-action row-action--primary" type="button" @click="action('down')">
          <NIcon size="16"><DownloadOutline/></NIcon>
        </button>
      </template>
      {{ t('index.direct_download') }}
    </NTooltip>

    <NTooltip trigger="hover">
      <template #trigger>
        <button class="row-action row-action--danger" type="button" @click="action('delete')">
          <NIcon size="16"><TrashOutline/></NIcon>
        </button>
      </template>
      {{ t('index.delete_row') }}
    </NTooltip>

    <NPopover placement="bottom-end" trigger="click" :show-arrow="false">
      <template #trigger>
        <button class="row-action" type="button">
          <NIcon size="17"><EllipsisHorizontal/></NIcon>
        </button>
      </template>
      <div class="row-menu">
        <button v-if="row.Status === 'running' || row.Status === 'pending'" type="button" @click="action('cancel')">
          <NIcon><CloseOutline/></NIcon><span>{{ t('index.cancel_down') }}</span>
        </button>
        <button type="button" @click="action('copy')">
          <NIcon><LinkOutline/></NIcon><span>{{ t('index.copy_link') }}</span>
        </button>
        <button v-if="row.OtherData?.source_url" type="button" @click="action('copySource')">
          <NIcon><CodeSlashOutline/></NIcon><span>{{ t('index.copy_source_api') }}</span>
        </button>
        <button v-if="row.Classify !== 'live' && row.Classify !== 'm3u8'" type="button" @click="action('open')">
          <NIcon><GlobeOutline/></NIcon><span>{{ t('index.open_link') }}</span>
        </button>
        <button v-if="row.DecodeKey" type="button" @click="action('decode')">
          <NIcon><LockOpenSharp/></NIcon><span>{{ t('index.video_decode') }}</span>
        </button>
        <button type="button" @click="action('json')">
          <NIcon><CopyOutline/></NIcon><span>{{ t('index.copy_data') }}</span>
        </button>
      </div>
    </NPopover>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n'
import {NIcon, NPopover, NTooltip} from 'naive-ui'
import {
  DownloadOutline,
  CopyOutline,
  GlobeOutline,
  LockOpenSharp,
  LinkOutline,
  CodeSlashOutline,
  CloseOutline,
  TrashOutline,
  EllipsisHorizontal
} from "@vicons/ionicons5"

const {t} = useI18n()
const props = defineProps<{row: any, index: number}>()
const emits = defineEmits(["action"])

const action = (type: string) => {
  if (type === 'down' && (props.row.Classify === 'live' || props.row.Classify === 'm3u8')) {
    window?.$message?.error(t("index.download_no_tip"))
    return
  }
  emits('action', props.row, props.index, type)
}
</script>

<style scoped>
.row-actions {
  display: flex;
  align-items: center;
  gap: 5px;
  --wails-draggable: no-drag;
}

.row-action {
  display: grid;
  width: 29px;
  height: 29px;
  place-items: center;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  color: var(--app-muted);
  background: var(--app-surface);
  cursor: pointer;
  transition: border-color 0.18s ease, color 0.18s ease, background 0.18s ease, transform 0.18s ease;
}

.row-action:hover {
  border-color: color-mix(in srgb, var(--app-muted) 50%, var(--app-border));
  color: var(--app-text);
  background: var(--app-surface-soft);
  transform: translateY(-1px);
}

.row-action--primary { color: var(--app-success); }
.row-action--danger { color: var(--app-danger); }

.row-menu {
  display: flex;
  min-width: 175px;
  flex-direction: column;
  gap: 3px;
  padding: 3px;
}

.row-menu button {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 9px;
  border: 0;
  border-radius: 7px;
  color: var(--app-text);
  background: transparent;
  cursor: pointer;
  font-size: 12px;
  text-align: left;
}

.row-menu button:hover { background: var(--app-surface-soft); }
.row-menu :deep(.n-icon) { color: var(--app-muted); }
</style>
