<template>
  <div id="content" class="app-shell">
    <Sider/>
    <main class="app-main">
      <RouterView v-slot="{ Component, route }">
        <keep-alive v-if="route.meta.keepAlive">
          <component :is="Component"/>
        </keep-alive>
        <component v-else :is="Component"/>
      </RouterView>
    </main>
  </div>
</template>
<script lang="ts" setup>
import Sider from "./Sider.vue"
</script>

<style scoped>
.app-shell {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--app-bg);
}

.app-main {
  min-width: 0;
  flex: 1;
  height: 100%;
  overflow: hidden;
  background:
    radial-gradient(circle at 90% -12%, rgba(220, 107, 47, 0.09), transparent 28rem),
    var(--app-bg);
}
</style>
