<template>
  <div class="capture-page">
    <header id="header" class="capture-header">
      <div class="capture-heading">
        <div>
          <p class="eyebrow">RES / SIGNAL DESK</p>
          <div class="title-line">
            <h1>{{ t('index.capture_center') }}</h1>
            <span class="status-pill" :class="{'status-pill--active': isProxy}">
              <span class="status-pill__dot"/>
              {{ isProxy ? t('index.capturing') : t('index.capture_ready') }}
            </span>
          </div>
          <p class="capture-description">{{ t('index.capture_desc') }}</p>
        </div>

        <NButton
            v-if="isProxy"
            size="large"
            type="error"
            secondary
            class="capture-button capture-button--stop"
            @click.stop="close"
        >
          <template #icon><NIcon><StopCircleOutline/></NIcon></template>
          {{ t('index.close_grab') }}
        </NButton>
        <NButton v-else size="large" type="primary" class="capture-button" @click.stop="open">
          <template #icon><NIcon><RadioOutline/></NIcon></template>
          {{ t('index.open_grab') }}
        </NButton>
      </div>

      <div class="metric-strip">
        <div class="metric-card">
          <span class="metric-card__icon metric-card__icon--orange"><NIcon><LayersOutline/></NIcon></span>
          <span class="metric-card__copy">
            <strong>{{ data.length }}</strong>
            <small>{{ t('index.resource_view') }}</small>
          </span>
        </div>
        <div class="metric-card">
          <span class="metric-card__icon metric-card__icon--blue"><NIcon><CodeSlashOutline/></NIcon></span>
          <span class="metric-card__copy">
            <strong>{{ apiData.length }}</strong>
            <small>{{ t('index.api_view') }}</small>
          </span>
        </div>
        <div class="metric-card">
          <span class="metric-card__icon metric-card__icon--green"><NIcon><GlobeOutline/></NIcon></span>
          <span class="metric-card__copy">
            <strong>{{ uniqueDomainCount }}</strong>
            <small>{{ t('index.unique_domains') }}</small>
          </span>
        </div>
        <div class="session-summary">
          <span>{{ t('index.session_records') }}</span>
          <strong>{{ currentDataCount }}</strong>
          <small>{{ t('index.results_visible', {count: visibleDataCount}) }}</small>
        </div>
      </div>
    </header>

    <section class="workspace-panel">
      <div class="workspace-topbar">
        <div class="view-switch" role="tablist">
          <button
              class="view-switch__item"
              :class="{'view-switch__item--active': activeView === 'resources'}"
              type="button"
              @click="activeView = 'resources'"
          >
            <NIcon size="17"><LayersOutline/></NIcon>
            <span>{{ t('index.resource_view') }}</span>
            <b>{{ data.length }}</b>
          </button>
          <button
              class="view-switch__item"
              :class="{'view-switch__item--active': activeView === 'apis'}"
              type="button"
              @click="activeView = 'apis'"
          >
            <NIcon size="17"><CodeSlashOutline/></NIcon>
            <span>{{ t('index.api_view') }}</span>
            <b>{{ apiData.length }}</b>
          </button>
        </div>

        <div class="workspace-actions">
          <NInput v-model:value="searchValue" clearable class="workspace-search" :placeholder="searchPlaceholder">
            <template #prefix><NIcon size="17"><SearchOutline/></NIcon></template>
          </NInput>

          <NButton v-if="activeView === 'resources'" type="primary" secondary @click.stop="batchDown">
            <template #icon><NIcon><DownloadOutline/></NIcon></template>
            {{ t('index.batch_download') }}
          </NButton>

          <NPopover placement="bottom-end" trigger="click" :show-arrow="false">
            <template #trigger>
              <NButton quaternary class="more-button">
                <template #icon><NIcon><EllipsisHorizontal/></NIcon></template>
                {{ t('index.secondary_actions') }}
              </NButton>
            </template>
            <div class="action-menu">
              <button v-if="activeView === 'resources'" type="button" @click="batchCancel">
                <NIcon><CloseOutline/></NIcon>{{ t('index.cancel_down') }}
              </button>
              <button v-if="activeView === 'resources'" type="button" @click="batchExport()">
                <NIcon><ArrowRedoCircleOutline/></NIcon>{{ t('index.batch_export') }}
              </button>
              <button v-if="activeView === 'resources'" type="button" @click="showImport = true">
                <NIcon><ServerOutline/></NIcon>{{ t('index.batch_import') }}
              </button>
              <button v-if="activeView === 'resources'" type="button" @click="batchExport('url')">
                <NIcon><LinkOutline/></NIcon>{{ t('index.export_url') }}
              </button>
              <button v-if="activeView === 'apis'" type="button" @click="batchExportApis('url')">
                <NIcon><LinkOutline/></NIcon>{{ t('index.export_api_url') }}
              </button>
              <button v-if="activeView === 'apis'" type="button" @click="batchExportApis()">
                <NIcon><ServerOutline/></NIcon>{{ t('index.export_api_data') }}
              </button>
              <button v-if="activeView === 'apis'" type="button" @click="batchExportApis('curl')">
                <NIcon><CodeSlashOutline/></NIcon>{{ t('index.export_api_curl') }}
              </button>
            </div>
          </NPopover>
        </div>
      </div>

      <div class="filter-bar">
        <div class="filter-bar__left">
          <span class="filter-label"><NIcon><FilterOutline/></NIcon>{{ t('index.current_target') }}</span>
          <NSelect
              v-if="activeView === 'resources'"
              v-model:value="resourcesType"
              multiple
              clearable
              :max-tag-count="2"
              :options="classify"
              class="type-select"
              :placeholder="t('index.grab_type')"
          />
          <NSelect
              v-else
              v-model:value="apiDomainFilter"
              filterable
              clearable
              :options="apiDomainOptions"
              class="domain-select"
              :placeholder="t('index.filter_domain')"
          />
          <button v-if="searchValue || apiDomainFilter" class="clear-filter" type="button" @click="clearFilters">
            {{ t('index.clear_filter') }}
          </button>
        </div>

        <div class="selection-actions">
          <span v-if="selectedCount" class="selection-count">{{ t('index.selected_count', {count: selectedCount}) }}</span>
          <NButton v-if="rememberChoice" quaternary type="error" @click.stop="clear">
            <template #icon><NIcon><TrashOutline/></NIcon></template>
            {{ t('index.clear_list') }}
          </NButton>
          <NPopconfirm v-else :show-icon="false" @positive-click="()=>{rememberChoice=rememberChoiceTmp;clear()}">
            <template #trigger>
              <NButton quaternary type="error">
                <template #icon><NIcon><TrashOutline/></NIcon></template>
                {{ t('index.clear_list') }}
              </NButton>
            </template>
            <div class="clear-confirm">
              <strong>{{ t('index.clear_list_tip') }}</strong>
              <NCheckbox v-model:checked="rememberChoiceTmp">{{ t('index.remember_clear_choice') }}</NCheckbox>
            </div>
          </NPopconfirm>
        </div>
      </div>

      <div class="table-shell">
        <NDataTable
            v-if="activeView === 'resources'"
            :columns="columns"
            :data="filteredData"
            :bordered="false"
            :max-height="tableHeight"
            :row-key="rowKey"
            :virtual-scroll="true"
            :header-height="46"
            :height-for-row="()=> 52"
            :checked-row-keys="checkedRowKeysValue"
            @update:checked-row-keys="handleCheck"
            @update:filters="updateFilters"
        >
          <template #empty><div class="table-empty"><NIcon size="28"><LayersOutline/></NIcon><strong>{{ emptyTitle }}</strong><span>{{ emptyHint }}</span></div></template>
        </NDataTable>
        <NDataTable
            v-else
            :columns="apiColumns"
            :data="filteredApiData"
            :bordered="false"
            :max-height="tableHeight"
            :row-key="apiRowKey"
            :virtual-scroll="true"
            :header-height="46"
            :height-for-row="()=> 60"
            :checked-row-keys="checkedApiRowKeysValue"
            @update:checked-row-keys="handleApiCheck"
        >
          <template #empty><div class="table-empty"><NIcon size="28"><CodeSlashOutline/></NIcon><strong>{{ emptyTitle }}</strong><span>{{ emptyHint }}</span></div></template>
        </NDataTable>
      </div>
    </section>

    <footer id="bottom" class="capture-footer">
      <span><span class="footer-status-dot" :class="{'footer-status-dot--active': isProxy}"/>{{ isProxy ? t('index.capturing') : t('index.capture_ready') }}</span>
      <div>
        <button type="button" @click="BrowserOpenURL(certUrl)">{{ t('footer.cert_download') }}</button>
        <button type="button" @click="BrowserOpenURL('https://github.com/putyy/res-downloader/issues')">{{ t('footer.help') }}</button>
        <button type="button" @click="BrowserOpenURL('https://github.com/putyy/res-downloader/releases')">{{ t('footer.update_log') }}</button>
      </div>
    </footer>

    <NModal
        v-model:show="showApiBodyModal"
        preset="card"
        class="api-response-modal w-[82vw] max-w-[1180px]"
        :title="apiBodyTitle"
    >
      <div class="flex flex-col gap-3">
        <div class="flex items-center justify-between gap-3 text-xs text-gray-400">
          <span class="ellipsis-2">{{ previewApiRow?.Url }}</span>
          <div class="flex items-center gap-2 shrink-0">
            <span v-if="previewApiRow?.BodyTruncated" class="text-amber-500">{{ t('index.response_body_truncated') }}</span>
            <NButton tertiary size="small" type="primary" @click="copyText(previewApiRow?.Body || '', t('index.response_body_empty'))">
              {{ t('index.copy_response_body') }}
            </NButton>
          </div>
        </div>
        <NInput
            type="textarea"
            readonly
            :value="previewApiRow?.Body || ''"
            :autosize="{minRows: 22, maxRows: 30}"
            class="font-mono text-xs"
        />
      </div>
    </NModal>
    <Preview v-model:showModal="showPreviewRow" :previewRow="previewRow"/>
    <ShowLoading :loadingText="loadingText" :isLoading="loading"/>
    <ImportJson v-model:showModal="showImport" @submit="handleImport"/>
    <Password v-model:showModal="showPassword" @submit="handlePassword"/>
  </div>
</template>

<script lang="ts" setup>
import {NButton, NIcon, NImage, NInput, NSpace, NTooltip, NPopover, NGradientText, NModal} from "naive-ui"
import {computed, h, onMounted, ref, watch} from "vue"
import type {appType} from "@/types/app"
import type {DataTableRowKey, ImageRenderToolbarProps, DataTableFilterState, DataTableBaseColumn} from "naive-ui"
import Preview from "@/components/Preview.vue"
import ShowLoading from "@/components/ShowLoading.vue"
// @ts-ignore
import {getDecryptionArray} from '@/assets/js/decrypt.js'
import {useIndexStore} from "@/stores"
import appApi from "@/api/app"
import Action from "@/components/Action.vue"
import ActionDesc from "@/components/ActionDesc.vue"
import ImportJson from "@/components/ImportJson.vue"
import {useEventStore} from "@/stores/event"
import {BrowserOpenURL, ClipboardSetText} from "../../wailsjs/runtime"
import Password from "@/components/Password.vue"
import ShowOrEdit from "@/components/ShowOrEdit.vue"
import {useI18n} from 'vue-i18n'
import {
  DownloadOutline,
  ArrowRedoCircleOutline,
  ServerOutline,
  SearchOutline,
  TrashOutline,
  CloseOutline,
  CopyOutline,
  LinkOutline,
  CodeSlashOutline,
  EyeOutline,
  RadioOutline,
  StopCircleOutline,
  LayersOutline,
  GlobeOutline,
  EllipsisHorizontal,
  FilterOutline
} from "@vicons/ionicons5"
import {useDialog} from 'naive-ui'
import * as bind from "../../wailsjs/go/core/Bind"
import {Quit} from "../../wailsjs/runtime"
import {DialogOptions} from "naive-ui/es/dialog/src/DialogProvider"
import {formatSize} from "@/func"

const {t} = useI18n()
const eventStore = useEventStore()
const dialog = useDialog()
const store = useIndexStore()
const isProxy = computed(() => {
  return store.isProxy
})
const certUrl = computed(() => {
  return store.baseUrl + "/api/cert"
})
const activeView = ref<"resources" | "apis">("resources")
const data = ref<any[]>([])
const apiData = ref<appType.ApiInfo[]>([])
const resourceSearchValue = ref("")
const apiSearchValue = ref("")
const apiDomainFilter = ref<string | null>(null)
const filterClassify = ref<string[]>([])
const currentDataCount = computed(() => {
  return activeView.value === "apis" ? apiData.value.length : data.value.length
})
const filteredData = computed(() => {
  let result = data.value

  if (filterClassify.value.length > 0) {
    result = result.filter(item => filterClassify.value.includes(item.Classify))
  }

  if (descriptionSearchValue.value) {
    result = result.filter(item => item.Description?.toLowerCase().includes(descriptionSearchValue.value.toLowerCase()))
  }

  if (urlSearchValue.value) {
    result = result.filter(item => item.Url?.toLowerCase().includes(urlSearchValue.value.toLowerCase()))
  }

  if (resourceSearchValue.value) {
    const keyword = resourceSearchValue.value.toLowerCase()
    result = result.filter(item =>
        item.Url?.toLowerCase().includes(keyword) ||
        item.Domain?.toLowerCase().includes(keyword) ||
        item.Description?.toLowerCase().includes(keyword) ||
        item.Classify?.toLowerCase().includes(keyword)
    )
  }

  return result
})
const filteredApiData = computed(() => {
  let result = apiData.value
  if (apiDomainFilter.value) {
    result = result.filter(item => item.Domain === apiDomainFilter.value)
  }
  if (apiSearchValue.value) {
    const keyword = apiSearchValue.value.toLowerCase()
    result = result.filter(item =>
        item.Url?.toLowerCase().includes(keyword) ||
        item.Body?.toLowerCase().includes(keyword) ||
        item.Domain?.toLowerCase().includes(keyword)
    )
  }
  return result
})
const apiDomainOptions = computed(() => {
  const domains = Array.from(new Set(apiData.value.map(item => item.Domain).filter(Boolean))).sort()
  return domains.map(domain => ({
    label: `${domain} (${apiData.value.filter(item => item.Domain === domain).length})`,
    value: domain
  }))
})

const searchValue = computed({
  get: () => activeView.value === "apis" ? apiSearchValue.value : resourceSearchValue.value,
  set: (value: string) => {
    if (activeView.value === "apis") {
      apiSearchValue.value = value
    } else {
      resourceSearchValue.value = value
    }
  }
})
const searchPlaceholder = computed(() => activeView.value === "apis" ? t("index.search_api") : t("index.search_resources"))
const uniqueDomainCount = computed(() => {
  const items = activeView.value === "apis" ? apiData.value : data.value
  return new Set(items.map(item => item.Domain).filter(Boolean)).size
})
const visibleDataCount = computed(() => activeView.value === "apis" ? filteredApiData.value.length : filteredData.value.length)
const selectedCount = computed(() => activeView.value === "apis" ? checkedApiRowKeysValue.value.length : checkedRowKeysValue.value.length)
const emptyTitle = computed(() => searchValue.value || apiDomainFilter.value ? t("index.no_results") : t("index.waiting_results"))
const emptyHint = computed(() => searchValue.value || apiDomainFilter.value ? t("index.no_results_hint") : t("index.waiting_results_hint"))

const clearFilters = () => {
  resourceSearchValue.value = ""
  apiSearchValue.value = ""
  apiDomainFilter.value = null
}

const tableHeight = ref(800)
const resourcesType = ref<string[]>(["all"])

const classifyAlias: { [key: string]: any } = {
  image: computed(() => t("index.image")),
  audio: computed(() => t("index.audio")),
  video: computed(() => t("index.video")),
  m3u8: computed(() => t("index.m3u8")),
  live: computed(() => t("index.live")),
  xls: computed(() => t("index.xls")),
  doc: computed(() => t("index.doc")),
  pdf: computed(() => t("index.pdf")),
  stream: computed(() => t("index.stream")),
  font: computed(() => t("index.font"))
}

const dwStatus = computed<any>(() => {
  return {
    ready: t("index.ready"),
    pending: t("index.pending"),
    running: t("index.running"),
    error: t("index.error"),
    done: t("index.done"),
    handle: t("index.handle")
  }
})

const maxConcurrentDownloads = computed(() => {
  return store.globalConfig.DownNumber
})

const classify = ref([
  {
    value: "all",
    label: computed(() => t("index.all")),
  },
])

const descriptionSearchValue = ref("")
const urlSearchValue = ref("")
const rememberChoice = ref(false)
const rememberChoiceTmp = ref(false)

function copyText(value: string, emptyMessage = t("common.copy_fail")) {
  if (!value) {
    window?.$message?.error(emptyMessage)
    return
  }
  ClipboardSetText(value).then((is: boolean) => {
    if (is) {
      window?.$message?.success(t("common.copy_success"))
    } else {
      window?.$message?.error(t("common.copy_fail"))
    }
  })
}

function shellQuote(value: string) {
  return "'" + String(value).replace(/'/g, "'\\''") + "'"
}

function buildCurl(row: appType.ApiInfo) {
  const lines = [`curl --location ${shellQuote(row.Url)}`]
  if (row.Method && row.Method !== "GET") {
    lines.push(`  -X ${row.Method}`)
  }

  try {
    const headers = JSON.parse(row.RequestHeaders || "{}")
    Object.entries(headers).forEach(([key, values]) => {
      const headerValue = Array.isArray(values) ? values[0] : values
      const lowerKey = key.toLowerCase()
      if (!headerValue || ["accept-encoding", "content-length", "host"].includes(lowerKey)) {
        return
      }
      lines.push(`  -H ${shellQuote(`${key}: ${headerValue}`)}`)
    })
  } catch (e) {
    console.log(e)
  }

  if (row.RequestBody) {
    lines.push(`  --data-raw ${shellQuote(row.RequestBody)}`)
  }

  return lines.join(" \\\n")
}

function showApiBody(row: appType.ApiInfo) {
  previewApiRow.value = row
  showApiBodyModal.value = true
}

const columns = ref<any[]>([
  {
    type: "selection",
  },
  {
    title: () => {
      if (checkedRowKeysValue.value.length > 0) {
        return h(NGradientText, {type: "success"}, t("index.choice") + `(${checkedRowKeysValue.value.length})`)
      }
      return t('index.domain')
    },
    key: "Domain",
    width: 90,
    render: (row: appType.MediaInfo) => {
      return h(NTooltip, {
        trigger: 'hover',
        placement: 'top'
      }, {
        trigger: () => h('span', {
          class: 'cursor-default'
        }, row.Domain),
        default: () => row.Url
      })
    }
  },
  {
    title: computed(() => t("index.type")),
    key: "Classify",
    width: 80,
    filterOptions: computed(() => Array.from(classify.value).slice(1)),
    filterMultiple: true,
    filter: (value: string, row: appType.MediaInfo) => {
      return !!~row.Classify.indexOf(String(value))
    },
    render: (row: appType.MediaInfo) => {
      const item = classify.value.find(item => item.value === row.Classify)
      return item ? item.label : row.Classify
    }
  },
  {
    title: computed(() => t("index.preview")),
    key: "Url",
    width: 80,
    render: (row: appType.MediaInfo) => {
      if (row.Classify === "image") {
        return h("div", {
          style: "width: 100%;max-height:80px;overflow:hidden;"
        }, h(NImage, {
          objectFit: "contain",
          lazy: true,
          "render-toolbar": renderToolbar,
          src: row.Url
        }))
      }
      return [
        h(
            NButton,
            {
              strong: true,
              tertiary: true,
              type: "info",
              size: "small",
              style: {
                margin: "2px"
              },
              onClick: () => {
                if (row.Classify === "audio" || row.Classify === "video" || row.Classify === "m3u8" || row.Classify === "live") {
                  previewRow.value = row
                  showPreviewRow.value = true
                }
              }
            },
            {
              default: () => {
                if (row.Classify === "audio" || row.Classify === "video" || row.Classify === "m3u8" || row.Classify === "live") {
                  return t("index.preview")
                }
                return t("index.preview_tip")
              }
            }
        ),
      ]
    }
  },
  {
    title: computed(() => t("index.status")),
    key: "Status",
    width: 80,
    render: (row: appType.MediaInfo, index: number) => {
      let status = "info"
      if (row.Status === "done" || row.Status === "running") {
        status = "success"
      } else if (row.Status === "pending") {
        status = "warning"
      }

      return h(
          NButton,
          {
            tertiary: true,
            type: status as any,
            size: "small",
            style: {
              margin: "2px"
            },
            onClick: () => {
              if (row.SavePath && row.Status === "done") {
                appApi.openFolder({filePath: row.SavePath})
              } else if (row.Status === "ready") {
                download(row, index)
              }
            }
          },
          {
            default: () => {
              return row.Status === "running" ? row.SavePath : dwStatus.value[row.Status as keyof typeof dwStatus]
            }
          }
      )
    }
  },
  {
    title: computed(() => t('index.description')),
    key: "Description",
    width: 150,
    render: (row: appType.MediaInfo, index: number) => {
      return h(ShowOrEdit, {
        value: row.Description,
        onUpdateValue(v: string) {
          data.value[index].Description = v
          cacheData()
        }
      })
    }
  },
  {
    title: computed(() => t("index.resource_size")),
    key: "Size",
    width: 120,
    sorter: (row1: appType.MediaInfo, row2: appType.MediaInfo) => row1.Size - row2.Size,
    render(row: appType.MediaInfo, index: number) {
      return formatSize(row.Size)
    }
  },
  {
    title: computed(() => t("index.save_path")),
    key: "SavePath",
    render(row: appType.MediaInfo, index: number) {
      return h("a",
          {
            href: "javascript:;",
            class: "ellipsis-2",
            style: {
              color: "#5a95d0"
            },
            onClick: () => {
              if (row.SavePath && row.Status === "done") {
                appApi.openFolder({filePath: row.SavePath})
              }
            }
          },
          row.Status === "running" ? "" : row.SavePath
      )
    }
  },
  {
    key: "actions",
    width: 130,
    render(row: appType.MediaInfo, index: number) {
      return h(Action, {key: index, row: row, index: index, onAction: dataAction})
    },
    title() {
      return h(ActionDesc)
    }
  }
])

const apiColumns = ref<any[]>([
  {
    type: "selection",
  },
  {
    title: computed(() => t("index.time")),
    key: "CreatedAt",
    width: 82,
  },
  {
    title: "Method",
    key: "Method",
    width: 86,
    render(row: appType.ApiInfo) {
      return h(NButton, {
        tertiary: true,
        size: "small",
        type: row.Method === "GET" ? "info" : "warning",
      }, row.Method)
    }
  },
  {
    title: computed(() => t("index.status")),
    key: "StatusCode",
    width: 88,
    render(row: appType.ApiInfo) {
      const type = row.StatusCode >= 400 ? "error" : row.StatusCode >= 300 ? "warning" : "success"
      return h(NButton, {
        tertiary: true,
        size: "small",
        type: type as any,
      }, String(row.StatusCode))
    }
  },
  {
    title: computed(() => t('index.api_url')),
    key: "Url",
    minWidth: 360,
    render(row: appType.ApiInfo) {
      return h(NTooltip, {
        trigger: 'hover',
        placement: 'top'
      }, {
        trigger: () => h("span", {
          class: "ellipsis-2 cursor-default",
          style: {color: "#5a95d0"}
        }, row.Url),
        default: () => row.Url
      })
    }
  },
  {
    title: computed(() => t("index.resource_size")),
    key: "Size",
    width: 110,
    sorter: (row1: appType.ApiInfo, row2: appType.ApiInfo) => row1.Size - row2.Size,
    render(row: appType.ApiInfo) {
      return formatSize(row.Size)
    }
  },
  {
    title: computed(() => t("index.hit_resources")),
    key: "MediaCount",
    width: 110,
    sorter: (row1: appType.ApiInfo, row2: appType.ApiInfo) => row1.MediaCount - row2.MediaCount,
    render(row: appType.ApiInfo) {
      if (!row.MediaCount) {
        return h("span", {class: "text-gray-400"}, "0")
      }
      return h(NButton, {
        tertiary: true,
        size: "small",
        type: "primary",
        onClick: () => copyText(row.MediaUrls.join("\n"))
      }, String(row.MediaCount))
    }
  },
  {
    title: computed(() => t("index.response_body")),
    key: "Body",
    minWidth: 260,
    render(row: appType.ApiInfo) {
      const body = row.Body || ""
      const preview = body.replace(/\s+/g, " ").slice(0, 180)
      return h(NTooltip, {
        trigger: 'hover',
        placement: 'top',
        style: {
          maxWidth: "640px",
          whiteSpace: "pre-wrap",
          wordBreak: "break-all"
        }
      }, {
        trigger: () => h("span", {
          class: "ellipsis-2 font-mono text-xs cursor-pointer hover:text-blue-400",
          onClick: () => showApiBody(row)
        }, preview + (row.BodyTruncated ? " ..." : "")),
        default: () => body.slice(0, 1200) + (row.BodyTruncated ? "\n..." : "")
      })
    }
  },
  {
    key: "actions",
    width: 190,
    title: computed(() => t("index.operation")),
    render(row: appType.ApiInfo) {
      return h(NSpace, {size: 4, wrapItem: false}, {
        default: () => [
          h(NTooltip, {trigger: 'hover'}, {
            trigger: () => h(NButton, {
              tertiary: true,
              circle: true,
              size: "small",
              type: "primary",
              onClick: () => copyText(row.Url)
            }, {
              icon: () => h(NIcon, null, {default: () => h(LinkOutline)})
            }),
            default: () => t("index.copy_link")
          }),
          h(NTooltip, {trigger: 'hover'}, {
            trigger: () => h(NButton, {
              tertiary: true,
              circle: true,
              size: "small",
              type: "warning",
              onClick: () => copyText(buildCurl(row))
            }, {
              icon: () => h(NIcon, null, {default: () => h(CodeSlashOutline)})
            }),
            default: () => t("index.copy_curl")
          }),
          h(NTooltip, {trigger: 'hover'}, {
            trigger: () => h(NButton, {
              tertiary: true,
              circle: true,
              size: "small",
              type: "success",
              onClick: () => showApiBody(row)
            }, {
              icon: () => h(NIcon, null, {default: () => h(EyeOutline)})
            }),
            default: () => t("index.preview_response_body")
          }),
          h(NTooltip, {trigger: 'hover'}, {
            trigger: () => h(NButton, {
              tertiary: true,
              circle: true,
              size: "small",
              type: "info",
              onClick: () => copyText(row.Body, t("index.response_body_empty"))
            }, {
              icon: () => h(NIcon, null, {default: () => h(CopyOutline)})
            }),
            default: () => t("index.copy_response_body")
          })
        ]
      })
    }
  }
])

const checkedRowKeysValue = ref<DataTableRowKey[]>([])
const checkedApiRowKeysValue = ref<DataTableRowKey[]>([])
const showPreviewRow = ref(false)
const previewRow = ref<appType.MediaInfo>()
const showApiBodyModal = ref(false)
const previewApiRow = ref<appType.ApiInfo | null>(null)
const apiBodyTitle = computed(() => {
  if (!previewApiRow.value) {
    return t("index.response_body")
  }
  return `${previewApiRow.value.Method} ${previewApiRow.value.StatusCode} ${previewApiRow.value.Domain}`
})
const loading = ref(false)
const loadingText = ref("")
const showImport = ref(false)
const showPassword = ref(false)
const downloadQueue = ref<appType.MediaInfo[]>([])
let activeDownloads = 0
let isOpenProxy = false
let isInstall = false

onMounted(() => {
  try {
    window.addEventListener("resize", () => {
      resetTableHeight()
    })
    loading.value = true
    handleInstall().then((is: boolean) => {
      isInstall = true
      loading.value = false
    })

    checkLoading()
    watch(showPassword, () => {
      if (!showPassword.value) {
        checkLoading()
      }
    })
  } catch (e) {
    window.$message?.error(JSON.stringify(e), {duration: 5000})
  }

  buildClassify()

  const temp = localStorage.getItem("resources-type")
  if (temp) {
    resourcesType.value = JSON.parse(temp).res
  } else {
    appApi.setType(resourcesType.value)
  }

  const cache = localStorage.getItem("resources-data")
  if (cache) {
    data.value = JSON.parse(cache)
  }

  const choiceCache = localStorage.getItem("remember-clear-choice")
  if (choiceCache === "1") {
    rememberChoice.value = true
  }

  watch(rememberChoice, (n, o) => {
    if (rememberChoice.value) {
      localStorage.setItem("remember-clear-choice", "1")
    } else {
      localStorage.removeItem("remember-clear-choice")
    }
  })

  resetTableHeight()

  eventStore.addHandle({
    type: "newResources",
    event: (res: appType.MediaInfo) => {
      if (store.globalConfig.InsertTail) {
        data.value.push(res)
      } else {
        data.value.unshift(res)
      }
      cacheData()
    }
  })

  eventStore.addHandle({
    type: "newApis",
    event: (res: appType.ApiInfo) => {
      const index = apiData.value.findIndex(item => item.UrlSign === res.UrlSign)
      if (index !== -1) {
        apiData.value[index] = res
        return
      }

      if (store.globalConfig.InsertTail) {
        apiData.value.push(res)
      } else {
        apiData.value.unshift(res)
      }

      if (apiData.value.length > 500) {
        if (store.globalConfig.InsertTail) {
          apiData.value.splice(0, apiData.value.length - 500)
        } else {
          apiData.value.splice(500)
        }
      }
    }
  })

  eventStore.addHandle({
    type: "downloadProgress",
    event: (res: { Id: string, SavePath: string, Status: string, Message: string }) => {
      switch (res.Status) {
        case "running":
          updateItem(res.Id, item => {
            item.SavePath = res.Message
            item.Status = 'running'
          })
          break
        case "done":
          updateItem(res.Id, item => {
            item.SavePath = res.SavePath
            item.Status = 'done'
          })
          if (activeDownloads > 0) {
            activeDownloads--
          }
          cacheData()
          checkQueue()
          break
        case "error":
          updateItem(res.Id, item => {
            item.SavePath = res.Message
            item.Status = 'error'
          })
          if (activeDownloads > 0) {
            activeDownloads--
          }
          cacheData()
          checkQueue()
          break
      }
    }
  })
})

watch(() => {
  return store.globalConfig.MimeMap
}, () => {
  buildClassify()
})

watch(resourcesType, (n, o) => {
  localStorage.setItem("resources-type", JSON.stringify({res: resourcesType.value}))
  appApi.setType(resourcesType.value)
})

const updateItem = (id: string, updater: (item: any) => void) => {
  const item = data.value.find(i => i.Id === id)
  if (item) updater(item)
}

function cacheData() {
  localStorage.setItem("resources-data", JSON.stringify(data.value))
}

const resetTableHeight = () => {
  try {
    const headerHeight = document.getElementById("header")?.offsetHeight || 0
    const bottomHeight = document.getElementById("bottom")?.offsetHeight || 0
    const workspaceTopbarHeight = document.querySelector<HTMLElement>(".workspace-topbar")?.offsetHeight || 0
    const filterBarHeight = document.querySelector<HTMLElement>(".filter-bar")?.offsetHeight || 0
    const height = document.documentElement.clientHeight || window.innerHeight
    tableHeight.value = Math.max(220, height - headerHeight - bottomHeight - workspaceTopbarHeight - filterBarHeight - 72)
  } catch (e) {
    console.log(e)
  }
}

const buildClassify = () => {
  const mimeMap = store.globalConfig.MimeMap ?? {}
  const seen = new Set()
  classify.value = [
    {value: "all", label: computed(() => t("index.all"))},
    ...Object.values(mimeMap)
        .filter(({Type}) => {
          if (seen.has(Type)) return false
          seen.add(Type)
          return true
        })
        .map(({Type}) => ({
          value: Type,
          label: classifyAlias[Type] ?? Type,
        })),
  ]
}

const dataAction = (row: appType.MediaInfo, index: number, type: string) => {
  switch (type) {
    case "down":
      download(row, index)
      break
    case "cancel":
      if (row.Status === "pending") {
        const queueIndex = downloadQueue.value.findIndex(item => item.Id === row.Id)
        if (queueIndex !== -1) {
          downloadQueue.value.splice(queueIndex, 1)
        }
        updateItem(row.Id, item => {
          item.Status = 'ready'
          item.SavePath = ''
        })
        cacheData()
      } else if (row.Status === "running") {
        appApi.cancel({id: row.Id}).then((res) => {
          updateItem(row.Id, item => {
            item.Status = 'ready'
            item.SavePath = ''
          })
          if (activeDownloads > 0) {
            activeDownloads--
          }
          cacheData()
          checkQueue()
          if (res.code === 0) {
            window?.$message?.error(res.message)
            return
          }
        })
      }
      break
    case "copy":
      ClipboardSetText(row.Url).then((is: boolean) => {
        if (is) {
          window?.$message?.success(t("common.copy_success"))
        } else {
          window?.$message?.error(t("common.copy_fail"))
        }
      })
      break
    case "copySource":
      if (!row.OtherData?.source_url) {
        window?.$message?.error(t("index.source_api_empty"))
        return
      }
      ClipboardSetText(row.OtherData.source_url).then((is: boolean) => {
        if (is) {
          window?.$message?.success(t("common.copy_success"))
        } else {
          window?.$message?.error(t("common.copy_fail"))
        }
      })
      break
    case "json":
      ClipboardSetText(encodeURIComponent(JSON.stringify(row))).then((is: boolean) => {
        if (is) {
          window?.$message?.success(t("common.copy_success"))
        } else {
          window?.$message?.error(t("common.copy_fail"))
        }
      })
      break
    case "open":
      BrowserOpenURL(row.Url)
      break
    case "decode":
      decodeWxFile(row, index)
      break
    case "delete":
      if (row.Status === "pending" || row.Status === "running") {
        window?.$message?.error(t("index.delete_tip"))
        return
      }
      appApi.delete({sign: [row.UrlSign]}).then(() => {
        data.value.splice(index, 1)
        cacheData()
      })
      break
  }
}

const renderToolbar = ({nodes}: ImageRenderToolbarProps) => {
  return [
    nodes.rotateCounterclockwise,
    nodes.rotateClockwise,
    nodes.resizeToOriginalSize,
    nodes.zoomOut,
    nodes.zoomIn,
    nodes.close
  ]
}

const rowKey = (row: appType.MediaInfo) => {
  return row.Id
}

const apiRowKey = (row: appType.ApiInfo) => {
  return row.Id
}

const handleCheck = (rowKeys: DataTableRowKey[]) => {
  checkedRowKeysValue.value = rowKeys
}

const handleApiCheck = (rowKeys: DataTableRowKey[]) => {
  checkedApiRowKeysValue.value = rowKeys
}

const updateFilters = (filters: DataTableFilterState, initiatorColumn: DataTableBaseColumn) => {
  filterClassify.value = filters.Classify as string[]
}

const batchDown = async () => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }

  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  data.value.forEach((item, index) => {
    if (checkedRowKeysValue.value.includes(item.Id) && item.Classify !== 'live' && item.Classify !== 'm3u8') {
      download(item, index)
    }
  })

  checkedRowKeysValue.value = []
}

const batchCancel = async () => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }
  loading.value = true
  const cancelTasks: Promise<any>[] = []
  data.value.forEach((item, index) => {
    if (!checkedRowKeysValue.value.includes(item.Id)) {
      return
    }

    if (item.Status === "pending") {
      const queueIndex = downloadQueue.value.findIndex(qItem => qItem.Id === item.Id)
      if (queueIndex !== -1) {
        downloadQueue.value.splice(queueIndex, 1)
      }
      item.Status = 'ready'
      item.SavePath = ''
      return
    }

    if (item.Status === "running") {
      if (activeDownloads > 0) {
        activeDownloads--
      }
      cancelTasks.push(appApi.cancel({id: item.Id}).then(() => {
        item.Status = 'ready'
        item.SavePath = ''
        checkQueue()
      }))
    }
  })
  await Promise.allSettled(cancelTasks)
  loading.value = false
  checkedRowKeysValue.value = []
  cacheData()
}

const batchExport = (type?: string) => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }

  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  loadingText.value = t("common.loading")
  loading.value = true

  let jsonData = data.value.filter(item => checkedRowKeysValue.value.includes(item.Id))

  if (type === "url") {
    jsonData = jsonData.map(item => item.Url)
  } else {
    jsonData = jsonData.map(item => encodeURIComponent(JSON.stringify(item)))
  }

  appApi.batchExport({content: jsonData.join("\n")}).then((res: appType.Res) => {
    loading.value = false
    if (res.code === 0) {
      window?.$message?.error(res.message)
      return
    }
    window?.$message?.success(t("index.import_success"))
    window?.$message?.info(t("index.save_path") + "：" + res.data?.file_name, {
      duration: 5000
    })
  })
}

const batchExportApis = (type?: string) => {
  const exportData = checkedApiRowKeysValue.value.length > 0
      ? apiData.value.filter(item => checkedApiRowKeysValue.value.includes(item.Id))
      : filteredApiData.value

  if (exportData.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }

  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  loadingText.value = t("common.loading")
  loading.value = true

  let exportContent = ""
  if (type === "url") {
    exportContent = exportData.map(item => item.Url).join("\n")
  } else if (type === "curl") {
    exportContent = exportData.map(item => buildCurl(item)).join("\n\n")
  } else {
    exportContent = exportData.map(item => encodeURIComponent(JSON.stringify(item))).join("\n")
  }

  appApi.batchExport({content: exportContent}).then((res: appType.Res) => {
    loading.value = false
    if (res.code === 0) {
      window?.$message?.error(res.message)
      return
    }
    window?.$message?.success(t("index.import_success"))
    window?.$message?.info(t("index.save_path") + "：" + res.data?.file_name, {
      duration: 5000
    })
  })
}

const uint8ArrayToBase64 = (bytes: any) => {
  return window.btoa(Array.from(bytes, (byte: any) => String.fromCharCode(byte)).join(''))
}

const download = (row: appType.MediaInfo, index: number) => {
  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  if (data.value.some(item => item.Id === row.Id && item.Status === "running")) {
    return
  }

  if (downloadQueue.value.some(item => item.Id === row.Id || item.Url === row.Url)) {
    return
  }

  if (activeDownloads >= maxConcurrentDownloads.value) {
    row.Status = "pending"
    downloadQueue.value.push(row)
    window?.$message?.info(t("index.download_queued", {count: downloadQueue.value.length}))
    return
  }

  startDownload(row, index)
}

const startDownload = (row: appType.MediaInfo, index: number) => {
  activeDownloads++

  const decodeStr = row.DecodeKey
      ? uint8ArrayToBase64(getDecryptionArray(row.DecodeKey))
      : ""

  appApi.download({...row, decodeStr}).then((res: appType.Res) => {
    if (res.code === 0) {
      window?.$message?.error(res.message)
    }
  })
}

const checkQueue = () => {
  if (downloadQueue.value.length > 0 && activeDownloads < maxConcurrentDownloads.value) {
    const nextItem = downloadQueue.value.shift()
    if (nextItem) {
      const index = data.value.findIndex(item => item.Id === nextItem.Id)
      if (index !== -1) {
        startDownload(nextItem, index)
      }
    }
  }
}

const open = () => {
  isOpenProxy = true
  store.openProxy().then((res: appType.Res) => {
    if (res.code === 1) {
      return
    }

    if (["darwin", "linux"].includes(store.envInfo.platform)) {
      showPassword.value = true
    } else {
      window.$message?.error(res.message)
    }
  })
}

const close = () => {
  store.unsetProxy()
}

const clear = async () => {
  if (activeView.value === "apis") {
    if (checkedApiRowKeysValue.value.length > 0) {
      apiData.value = apiData.value.filter(item => !checkedApiRowKeysValue.value.includes(item.Id))
      checkedApiRowKeysValue.value = []
    } else {
      apiData.value = []
    }
    return
  }

  const newData = [] as any[]
  const signs: string[] = []
  if (checkedRowKeysValue.value.length > 0) {
    data.value.forEach((item, index) => {
      if (checkedRowKeysValue.value.includes(item.Id) && item.Status !== "pending" && item.Status !== "running") {
        signs.push(item.UrlSign)
      } else {
        newData.push(item)
      }
    })
    checkedRowKeysValue.value = []
  } else {
    data.value.forEach((item, index) => {
      if (item.Status === "pending" || item.Status === "running") {
        newData.push(item)
      } else {
        signs.push(item.UrlSign)
      }
    })
  }
  await appApi.delete({sign: signs})
  data.value = newData
  cacheData()
}

const decodeWxFile = (row: appType.MediaInfo, index: number) => {
  if (!row.DecodeKey) {
    window?.$message?.error(t("index.video_decode_no"))
    return
  }
  appApi.openFileDialog().then((res: appType.Res) => {
    if (res.code === 0) {
      window?.$message?.error(res.message)
      return
    }
    if (res.data.file) {
      loadingText.value = t("index.video_decode_loading")
      loading.value = true
      appApi.wxFileDecode({
        ...row,
        filename: res.data.file,
        decodeStr: uint8ArrayToBase64(getDecryptionArray(row.DecodeKey))
      }).then((res: appType.Res) => {
        loading.value = false
        if (res.code === 0) {
          window?.$message?.error(res.message)
          return
        }
        data.value[index].SavePath = res.data.save_path
        data.value[index].Status = "done"
        cacheData()
        window?.$message?.success(t("index.video_decode_success"))
      })
    }
  })
}

const handleImport = (content: string) => {
  if (!content) {
    window?.$message?.error(t("view.import_empty"))
    return
  }
  let newItems = [] as any[]
  content.split("\n").forEach((line, index) => {
    try {
      let res = JSON.parse(decodeURIComponent(line))
      if (res && res?.Id) {
        res.Id = res.Id + Math.floor(Math.random() * 100000)
        res.SavePath = ""
        res.Status = "ready"
        newItems.push(res)
      }
    } catch (e) {
      console.log(e)
    }
  })
  if (newItems.length > 0) {
    data.value = [...newItems, ...data.value]
    cacheData()
  }
  showImport.value = false
}

const handlePassword = async (password: string, isCache: boolean) => {
  const res = await appApi.setSystemPassword({password, isCache})
  if (res.code === 0) {
    window.$message?.error(res.message)
    return
  }

  if (isOpenProxy) {
    showPassword.value = false
    store.openProxy()
    return
  }

  handleInstall().then((is: boolean) => {
    if (is) {
      showPassword.value = false
    }
  })
}

const handleInstall = async () => {
  isOpenProxy = false
  const res = await appApi.install()
  if (res.code === 1) {
    store.globalConfig.AutoProxy && store.openProxy()
    return true
  }

  window.$message?.error(res.message, {duration: 5000})

  if (store.envInfo.platform === "windows" && res.message.includes("Access is denied")) {
    window.$message?.error(t("index.win_install_tip"))
  } else if (["darwin", "linux"].includes(store.envInfo.platform)) {
    showPassword.value = true
  }
  return false
}

const checkLoading = () => {
  setTimeout(() => {
    if (loading.value && !isInstall && !showPassword.value) {
      dialog.warning({
        title: t("index.start_err_tip"),
        content: t("index.start_err_content"),
        positiveText: t("index.start_err_positiveText"),
        negativeText: t("index.start_err_negativeText"),
        draggable: false,
        closeOnEsc: false,
        closable: false,
        maskClosable: false,
        onPositiveClick: () => {
          bind.ResetApp()
        },
        onNegativeClick: () => {
          Quit()
        }
      } as DialogOptions)
    }
  }, 6000)
}
</script>

<style scoped>
.capture-page {
  display: flex;
  height: 100%;
  min-width: 0;
  flex-direction: column;
  gap: 14px;
  padding: 24px 26px 0;
  overflow: hidden;
  color: var(--app-text);
}

.capture-header {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  gap: 18px;
  animation: rise-in 0.42s ease-out both;
}

.capture-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
}

.eyebrow {
  margin-bottom: 5px;
  color: var(--app-accent);
  font-family: 'SFMono-Regular', 'Cascadia Code', monospace;
  font-size: 9px;
  font-weight: 750;
  letter-spacing: 0.18em;
}

.title-line {
  display: flex;
  align-items: center;
  gap: 13px;
}

.title-line h1 {
  color: var(--app-text);
  font-size: clamp(24px, 2.3vw, 32px);
  font-weight: 750;
  letter-spacing: -0.035em;
  line-height: 1.15;
}

.capture-description {
  margin-top: 6px;
  color: var(--app-muted);
  font-size: 12px;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 9px;
  border: 1px solid var(--app-border);
  border-radius: 999px;
  color: var(--app-muted);
  background: color-mix(in srgb, var(--app-surface) 82%, transparent);
  font-size: 10px;
  font-weight: 650;
}

.status-pill__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #94a199;
}

.status-pill--active {
  border-color: rgba(35, 131, 95, 0.25);
  color: var(--app-success);
  background: rgba(35, 131, 95, 0.08);
}

.status-pill--active .status-pill__dot {
  background: #2da274;
  box-shadow: 0 0 0 4px rgba(45, 162, 116, 0.1);
  animation: signal-pulse 1.6s ease-in-out infinite;
}

.capture-button {
  min-width: 142px;
  height: 44px;
  border-radius: 12px;
  box-shadow: 0 10px 20px rgba(220, 107, 47, 0.16);
  --wails-draggable: no-drag;
}

.capture-button--stop {
  box-shadow: none;
}

.metric-strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(128px, 176px)) minmax(175px, 1fr);
  gap: 10px;
}

.metric-card,
.session-summary {
  min-height: 62px;
  border: 1px solid var(--app-border);
  border-radius: 14px;
  background: color-mix(in srgb, var(--app-surface) 88%, transparent);
  box-shadow: 0 6px 20px rgba(25, 37, 31, 0.025);
}

.metric-card {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 11px 13px;
}

.metric-card__icon {
  display: grid;
  width: 35px;
  height: 35px;
  flex: none;
  place-items: center;
  border-radius: 11px;
  font-size: 17px;
}

.metric-card__icon--orange { color: #c85b24; background: rgba(220, 107, 47, 0.11); }
.metric-card__icon--blue { color: #477d9f; background: rgba(71, 125, 159, 0.11); }
.metric-card__icon--green { color: #23835f; background: rgba(35, 131, 95, 0.11); }

.metric-card__copy {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.metric-card__copy strong,
.session-summary strong {
  color: var(--app-text);
  font-family: 'DIN Alternate', 'Avenir Next Condensed', sans-serif;
  font-size: 21px;
  font-weight: 700;
  line-height: 1;
}

.metric-card__copy small,
.session-summary span,
.session-summary small {
  overflow: hidden;
  color: var(--app-muted);
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 0.035em;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.metric-card__copy small { margin-top: 5px; }

.session-summary {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  padding: 11px 15px;
  background:
    linear-gradient(110deg, rgba(220, 107, 47, 0.08), transparent 58%),
    color-mix(in srgb, var(--app-surface) 88%, transparent);
}

.session-summary strong {
  grid-row: span 2;
  font-size: 27px;
}

.workspace-panel {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--app-border);
  border-radius: 18px;
  background: var(--app-surface);
  box-shadow: var(--app-shadow);
  animation: rise-in 0.48s 0.05s ease-out both;
}

.workspace-topbar,
.filter-bar {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.workspace-topbar {
  min-height: 61px;
  padding: 11px 14px 10px;
  border-bottom: 1px solid var(--app-border);
}

.view-switch {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 12px;
  background: var(--app-surface-soft);
}

.view-switch__item {
  display: flex;
  height: 34px;
  align-items: center;
  gap: 7px;
  padding: 0 11px;
  border: 0;
  border-radius: 9px;
  color: var(--app-muted);
  background: transparent;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  transition: color 0.18s ease, background 0.18s ease, box-shadow 0.18s ease;
}

.view-switch__item b {
  min-width: 19px;
  padding: 2px 5px;
  border-radius: 999px;
  color: inherit;
  background: rgba(111, 126, 118, 0.11);
  font-family: 'DIN Alternate', sans-serif;
  font-size: 10px;
  font-weight: 700;
}

.view-switch__item--active {
  color: var(--app-text);
  background: var(--app-surface);
  box-shadow: 0 3px 10px rgba(31, 45, 38, 0.08), inset 0 0 0 1px var(--app-border);
}

.view-switch__item--active :deep(.n-icon) { color: var(--app-accent); }

.workspace-actions,
.filter-bar__left,
.selection-actions {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 8px;
}

.workspace-search { width: clamp(190px, 24vw, 310px); }
.more-button { color: var(--app-muted); }

.action-menu {
  display: flex;
  min-width: 190px;
  flex-direction: column;
  gap: 3px;
  padding: 3px;
}

.action-menu button {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border: 0;
  border-radius: 8px;
  color: var(--app-text);
  background: transparent;
  cursor: pointer;
  font-size: 12px;
  text-align: left;
}

.action-menu button:hover { background: var(--app-surface-soft); }
.action-menu :deep(.n-icon) { color: var(--app-muted); }

.filter-bar {
  min-height: 49px;
  padding: 7px 14px;
  border-bottom: 1px solid var(--app-border);
  background: color-mix(in srgb, var(--app-surface-soft) 72%, var(--app-surface));
}

.filter-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--app-muted);
  font-size: 10px;
  font-weight: 650;
  white-space: nowrap;
}

.type-select { width: 235px; }
.domain-select { width: min(320px, 28vw); }

.clear-filter {
  border: 0;
  color: var(--app-accent);
  background: transparent;
  cursor: pointer;
  font-size: 11px;
}

.selection-count {
  padding: 5px 8px;
  border-radius: 7px;
  color: var(--app-success);
  background: rgba(35, 131, 95, 0.09);
  font-size: 10px;
  font-weight: 650;
  white-space: nowrap;
}

.clear-confirm {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 3px;
}

.clear-confirm strong { font-size: 13px; font-weight: 650; }

.table-shell {
  min-height: 0;
  flex: 1;
  overflow: hidden;
}

.table-shell :deep(.n-data-table) {
  --wails-draggable: no-drag;
}

.table-shell :deep(.n-data-table-th) {
  color: var(--app-muted);
  font-size: 10px;
  letter-spacing: 0.035em;
}

.table-shell :deep(.n-data-table-td) {
  font-size: 12px;
}

.table-shell :deep(.n-data-table-tr:last-child .n-data-table-td) {
  border-bottom: 0;
}

.table-empty {
  display: flex;
  min-height: 210px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--app-muted);
}

.table-empty :deep(.n-icon) {
  margin-bottom: 10px;
  color: color-mix(in srgb, var(--app-muted) 62%, transparent);
}

.table-empty strong {
  color: var(--app-text);
  font-size: 13px;
  font-weight: 650;
}

.table-empty span {
  margin-top: 4px;
  font-size: 11px;
}

.capture-footer {
  display: flex;
  min-height: 29px;
  flex: none;
  align-items: flex-start;
  justify-content: space-between;
  color: var(--app-muted);
  font-size: 9px;
}

.capture-footer > span {
  display: flex;
  align-items: center;
  gap: 6px;
}

.capture-footer div { display: flex; gap: 15px; }

.capture-footer button {
  padding: 0;
  border: 0;
  color: var(--app-muted);
  background: transparent;
  cursor: pointer;
}

.capture-footer button:hover { color: var(--app-accent); }

.footer-status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #9ba69f;
}

.footer-status-dot--active {
  background: var(--app-success);
  box-shadow: 0 0 0 3px rgba(35, 131, 95, 0.1);
}

@keyframes rise-in {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes signal-pulse {
  50% { opacity: 0.55; transform: scale(0.86); }
}

@media (max-width: 1080px) {
  .capture-page { padding-inline: 18px; }
  .metric-strip { grid-template-columns: repeat(3, minmax(112px, 1fr)); }
  .session-summary { display: none; }
  .view-switch__item { padding-inline: 9px; }
  .workspace-actions :deep(.n-button__content) { font-size: 0; }
  .workspace-actions :deep(.n-button__icon) { margin: 0; }
}

@media (max-height: 700px) {
  .capture-page { gap: 9px; padding-top: 15px; }
  .capture-header { gap: 10px; }
  .capture-description { display: none; }
  .metric-card, .session-summary { min-height: 50px; }
  .metric-card { padding-block: 7px; }
  .metric-card__icon { width: 31px; height: 31px; }
  .workspace-topbar { min-height: 53px; padding-block: 7px; }
  .filter-bar { min-height: 43px; padding-block: 5px; }
  .capture-footer { min-height: 22px; }
}

@media (prefers-reduced-motion: reduce) {
  .capture-header,
  .workspace-panel,
  .status-pill--active .status-pill__dot { animation: none; }
}
</style>
