<script setup>
const features = [
  {
    title: 'Gratis',
    description: 'Aplikasi ini dapat digunakan sepenuhnya tanpa biaya sepeserpun.',
    icon: 'i-heroicons-gift'
  },
  {
    title: 'Buat Kupon Otomatis',
    description: 'Generate kupon qurban dengan cepat dan mudah untuk para mustahik.',
    icon: 'i-heroicons-ticket'
  },
  {
    title: 'Tabungan Qurban',
    description: 'Kelola rencana tabungan qurban jamaah dengan sistematis.',
    icon: 'i-heroicons-banknotes'
  }
]
</script>

<template>
  <div>
    <!-- Hero Section: Desktop Grid 2 Column -->
    <section class="py-16 md:py-24">
      <UContainer>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
          <!-- Column 1: Logo Vertical -->
          <div class="flex justify-center md:justify-center border">
            <img src="~/assets/logo-vertical.png" alt="Qurban Logo" class="w-48 md:w-64" />
          </div>

          <!-- Column 2: Tagline & Button -->
          <div class="text-center md:text-left">
            <h1 class="text-3xl md:text-5xl font-extrabold text-brand-primary leading-tight">
              Solusi Gratis Manajemen Qurban
            </h1>
            <p class="mt-6 text-lg md:text-xl text-gray-600 leading-relaxed">
              Atur distribusi daging, buat kupon otomatis, dan kelola kegiatan qurban dengan mudah.
            </p>
            <div class="mt-10">
              <AppButton size="xl" variant="primary" class="shadow-lg">
                Daftar Sekarang
              </AppButton>
            </div>
          </div>
        </div>
      </UContainer>
    </section>

    <!-- Features Section -->
    <section class="py-16 bg-white">
      <UContainer>
        <div class="text-center mb-12">
          <h2 class="text-2xl md:text-3xl font-bold text-gray-900">Fitur Unggulan</h2>
          <div class="h-1 w-20 bg-brand-accent mx-auto mt-4 rounded-full"></div>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <UCard v-for="feature in features" :key="feature.title" class="text-center hover:shadow-md transition-shadow">
            <div class="flex flex-col items-center p-4">
              <div class="w-16 h-16 bg-brand-surface rounded-2xl flex items-center justify-center mb-6 text-brand-primary">
                <UIcon :name="feature.icon" class="w-8 h-8" />
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-3">{{ feature.title }}</h3>
              <p class="text-gray-600 leading-relaxed">
                {{ feature.description }}
              </p>
            </div>
          </UCard>
        </div>
      </UContainer>
    </section>
  </div>
</template>

<style scoped>
.text-brand-primary { color: var(--color-primary); }
.bg-brand-surface { background-color: var(--color-surface); }
.bg-brand-accent { background-color: var(--color-accent); }
</style>
