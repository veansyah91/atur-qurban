<template>
  <footer class="bg-brand-surface border-t border-gray-200 py-12">
    <UContainer>
      <div class="flex flex-col items-center text-center gap-6">
        <p class="text-gray-500 font-medium">Support by: keuanganumum.com</p>
        <img src="~/assets/keuangan-umum.jpeg" alt="Keuangan Umum Logo" class="h-12 md:h-16 object-contain rounded" />
        <p class="text-sm text-gray-400 mt-4">
          &copy; {{ new Date().getFullYear() }} Atur Qurban. All rights reserved.
        </p>
      </div>
    </UContainer>
  </footer>
</template>

<style scoped>
.bg-brand-surface { background-color: var(--color-surface); }
</style>
