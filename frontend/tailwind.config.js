/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        bg: '#0a0e17',
        card: '#111827',
        cardHover: '#1a2332',
        border: '#1e293b',
        accent: '#f59e0b',
        accent2: '#10b981',
        red: '#ef4444',
        green: '#22c55e',
        blue: '#3b82f6',
        purple: '#a78bfa',
        text: '#f8fbff',
        textDim: '#d6e0ee',
        textMuted: '#b4c2d6',
        gold: '#fbbf24',
        orange: '#f97316',
      },
    },
  },
  plugins: [],
}