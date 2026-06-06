import { defineTailwindConfig } from 'tailwindcss'

export default {
  theme: {
    extend: {
      colors: {
        brand: {
          primary: '#094927',
          'primary-dark': '#063318',
          'primary-light': '#12693a',
          accent: '#ce981d',
          'accent-dark': '#a67b15',
          secondary: '#12693a',
          warning: '#f59e0b',
          error: '#dc2626',
          surface: '#f8f9fa',
        }
      }
    }
  }
}
