<script setup>
const isOpen = ref(false)

const navLinks = [
  { label: 'Masuk', to: '#' },
  { label: 'Daftar', to: '#' }
]
</script>

<template>
  <header class="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200">
    <UContainer>
      <div class="flex items-center justify-between h-24 md:h-32">
        <!-- Logo Horizontal - Zoomed in 10x (scale-10 is huge, typically we use zoom via scale or explicit size) -->
        <!-- Note: 10x is very large for a header, applying scale and adjusting container -->
        <div class="flex items-center overflow-visible">
          <img 
            src="~/assets/logo-horizontal.png" 
            alt="Logo Qurban" 
            class="h-8 md:h-8 object-contain transform scale-[3] origin-left transition-transform" 
          />
        </div>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center gap-4">
          <AppButton v-for="link in navLinks" :key="link.label" :to="link.to" variant="ghost">
            {{ link.label }}
          </AppButton>
        </div>

        <!-- Mobile Hamburger -->
        <div class="md:hidden">
          <UButton
            icon="i-heroicons-bars-3"
            color="gray"
            variant="ghost"
            aria-label="Menu"
            @click="isOpen = true"
          />
        </div>
      </div>
    </UContainer>

    <!-- Sidebar Mobile -->
    <USlideover v-model="isOpen" side="right">
      <div class="p-4 flex-1">
        <div class="flex items-center justify-between mb-8">
          <h3 class="text-base font-semibold text-gray-900">Menu</h3>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="isOpen = false" />
        </div>
        <div class="flex flex-col gap-4">
          <AppButton
            v-for="link in navLinks"
            :key="link.label"
            :to="link.to"
            variant="primary"
            class="w-full"
            @click="isOpen = false"
          >
            {{ link.label }}
          </AppButton>
        </div>
      </div>
    </USlideover>
  </header>
</template>
