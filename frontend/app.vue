<script setup>
const store = useMainStore()
</script>

<template>
  <div class="min-h-screen bg-brand-surface">
    <UContainer>
      <header class="flex items-center justify-between py-6 border-b">
        <div class="flex items-center gap-4">
          <img src="~/assets/logo-horizontal.png" alt="Logo" class="h-10 object-contain" />
          <h1 class="text-2xl font-bold text-brand-primary">{{ store.appName }}</h1>
        </div>
        <UNavbar>
          <UButton color="brand-primary" variant="ghost">Home</UButton>
          <UButton color="brand-primary" variant="ghost">About</UButton>
        </UNavbar>
      </header>

      <main class="py-10">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
          <div>
            <h2 class="text-4xl font-bold text-brand-primary-dark mb-4">
              Welcome to Our Project
            </h2>
            <p class="text-lg text-gray-700 mb-6">
              This project is built with Nuxt.js, Nuxt UI, and Pinia.
              The colors and assets are configured according to the planning document.
            </p>
            
            <div class="flex gap-4 items-center">
              <UButton 
                size="xl" 
                class="bg-brand-primary hover:bg-brand-primary-dark text-white"
                @click="store.increment"
              >
                Counter: {{ store.counter }}
              </UButton>
              <UButton 
                size="xl" 
                variant="outline"
                class="border-brand-accent text-brand-accent-dark hover:bg-brand-accent hover:text-white"
              >
                Explore More
              </UButton>
            </div>
          </div>
          
          <div class="flex justify-center">
            <img src="~/assets/keuangan-umum.jpeg" alt="Illustration" class="rounded-xl shadow-2xl max-w-full h-auto" />
          </div>
        </div>

        <section class="mt-16 grid grid-cols-1 sm:grid-cols-3 gap-6">
          <UCard class="border-t-4 border-brand-primary">
            <template #header>
              <div class="flex items-center gap-2 text-brand-primary font-bold">
                <UIcon name="i-heroicons-shield-check" />
                Primary Style
              </div>
            </template>
            Menggunakan warna primary #094927 untuk elemen utama.
          </UCard>
          
          <UCard class="border-t-4 border-brand-accent">
            <template #header>
              <div class="flex items-center gap-2 text-brand-accent-dark font-bold">
                <UIcon name="i-heroicons-star" />
                Accent Style
              </div>
            </template>
            Warna accent #ce981d digunakan untuk penekanan atau highlight.
          </UCard>

          <UCard class="border-t-4 border-brand-secondary">
            <template #header>
              <div class="flex items-center gap-2 text-brand-secondary font-bold">
                <UIcon name="i-heroicons-cog" />
                Secondary
              </div>
            </template>
            Secondary color terintegrasi dengan nuansa hijau yang konsisten.
          </UCard>
        </section>
      </main>

      <footer class="mt-20 py-8 border-t text-center text-gray-500">
        <div class="flex justify-center mb-4">
          <img src="~/assets/logo-vertical.png" alt="Vertical Logo" class="h-12 grayscale opacity-50" />
        </div>
        <p>&copy; 2024 Nuxt Brand Project. Built with Nuxt UI.</p>
      </footer>
    </UContainer>
  </div>
</template>

<style>
/* Custom utility mappings if needed */
.bg-brand-primary { background-color: var(--color-primary); }
.text-brand-primary { color: var(--color-primary); }
.bg-brand-surface { background-color: var(--color-surface); }
.text-brand-primary-dark { color: var(--color-primary-dark); }
.text-brand-accent-dark { color: var(--color-accent-dark); }
.border-brand-primary { border-color: var(--color-primary); }
.border-brand-accent { border-color: var(--color-accent); }
.border-brand-secondary { border-color: var(--color-secondary); }
.text-brand-secondary { color: var(--color-secondary); }
</style>
