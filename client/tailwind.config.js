/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte}', './public/index.html'],
  theme: {
    extend: {
      colors: {
        'je-orange': '#FF8000',
        'je-orange-dark': '#E67300',
        'je-orange-light': '#FFF3E6',
        'je-green': '#00B167',
        'je-green-dark': '#009956',
        'je-dark': '#2E3333',
        'je-grey': '#585C5C',
        'je-grey-light': '#F5F5F5',
        'je-grey-border': '#E0E0E0',
      },
      fontFamily: {
        sans: ['"JustEatBasis"', 'system-ui', '-apple-system', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
