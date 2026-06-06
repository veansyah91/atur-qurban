<template>
  <div class="min-h-screen bg-brand-surface font-sans flex flex-col">
    <AppHeader />
    <main class="flex-grow">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>

<style>
.bg-brand-surface { background-color: var(--color-surface); }
</style>
