/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte}', './public/index.html'],
  theme: {
    extend: {
      colors: {
        'sand': '#F5E6D3',
        'sand-dark': '#E8D5BC',
        'sand-light': '#FDF8F2',
        'ocean': '#2B8A9E',
        'ocean-dark': '#1E6F80',
        'ocean-light': '#E6F4F7',
        'sunset': '#F4845F',
        'sunset-dark': '#E06840',
        'sunset-light': '#FEF0EC',
        'palm': '#4A7C59',
        'palm-dark': '#3A6247',
        'palm-light': '#EDF5EF',
        'driftwood': '#2C2416',
        'driftwood-mid': '#5C4A32',
        'coconut': '#F9F6F1',
        'coral': '#FF6B6B',
        'mango': '#FFB347',
        'papaya': '#FFDAB9',
        'lime': '#A8E06A',
        'guava': '#FF69B4',
      },
      fontFamily: {
        sans: ['"Nunito"', 'system-ui', '-apple-system', 'sans-serif'],
        display: ['"Pacifico"', 'cursive'],
        surf: ['"Permanent Marker"', 'cursive'],
      },
      backgroundImage: {
        'wave-pattern': "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 100'%3E%3Cpath fill='%23F5E6D3' d='M0,40 C360,80 720,0 1080,40 C1260,60 1380,50 1440,40 L1440,100 L0,100 Z'/%3E%3C/svg%3E\")",
      },
      keyframes: {
        'wave': {
          '0%, 100%': { transform: 'translateX(0)' },
          '50%': { transform: 'translateX(-25px)' },
        },
        'sway': {
          '0%, 100%': { transform: 'rotate(-2deg)' },
          '50%': { transform: 'rotate(2deg)' },
        },
        'float': {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-8px)' },
        },
      },
      animation: {
        'wave': 'wave 6s ease-in-out infinite',
        'sway': 'sway 4s ease-in-out infinite',
        'float': 'float 3s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
