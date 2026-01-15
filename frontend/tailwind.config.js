/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				nascar: {
					red: '#e31837',
					blue: '#1a1a6c',
					yellow: '#ffd100'
				}
			}
		}
	},
	plugins: []
};
