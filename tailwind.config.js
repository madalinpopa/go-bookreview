/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["ui/**/*.{tmpl,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
