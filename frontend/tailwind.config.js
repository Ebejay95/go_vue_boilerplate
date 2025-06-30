/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
	  "./public/**/*.html",
	  "./src/**/*.{vue,js,ts,jsx,tsx}"
	],
	theme: {
	  extend: {
		fontFamily: {
		  sans: ['Avenir', 'Helvetica', 'Arial', 'sans-serif']
		}
	  },
	},
	plugins: [],
  }
