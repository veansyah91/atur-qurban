import { defineStore } from 'pinia'

export const useMainStore = defineStore('main', {
  state: () => ({
    counter: 0,
    appName: 'Nuxt Brand App'
  }),
  actions: {
    increment() {
      this.counter++
    }
  }
})
