<template>
  <router-view />
</template>

<script setup>
</script>

<style>
:root {
  --primary-color: #1e3a8a;
  --secondary-color: #2563eb;
  --bg-color: #f8fafc;
  --text-main: #334155;
  --text-muted: #64748b;
  --border-color: #e2e8f0;
  --star-color: #f59e0b;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-main);
  line-height: 1.5;
}
</style>
