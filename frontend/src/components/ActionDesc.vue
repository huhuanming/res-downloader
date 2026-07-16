<template>
  <span class="action-heading">{{ t('index.operation') }}</span>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n"
const {t} = useI18n()
</script>

<style scoped>
.action-heading {
  color: var(--app-muted);
  font-size: 10px;
  font-weight: 650;
  letter-spacing: 0.035em;
}
</style>
