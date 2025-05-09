/** @type {import('tailwindcss').Config} */
import animate from 'tailwindcss-animate';

export default {
	darkMode: ['variant', ['.dark &', '[data-kb-theme="dark"] &']],
	content: ['./src/**/*.{ts,tsx}'],
	prefix: '',
	theme: {
		container: {
			center: true,
			padding: '2rem',
			screens: {
				'2xl': '1400px',
			},
		},
		extend: {
			colors: {
				"header_bg_color": "#efeff3",
				"bg_color": "#ffffff",
				"hint_color": "#999999",
				"accent_text_color": "#2481cc",
				"section_bg_color": "#ffffff",
				"bottom_bar_bg_color": "#e4e4e4",
				"section_header_text_color": "#6d6d71",
				"text_color": "#000000",
				"button_text_color": "#ffffff",
				"destructive_text_color": "#ff3b30",
				"section_separator_color": "#eaeaea",
				"subtitle_text_color": "#999999",
				"secondary_bg_color": "#efeff3",
				"link_color": "#2481cc",
				"button_color": "#2481cc",
				border: 'var(--border)',
				input: 'var(--input)',
				ring: 'var(--ring)',
				background: 'var(--background)',
				foreground: 'var(--foreground)',
				primary: {
					DEFAULT: 'var(--primary)',
					foreground: 'var(--primary-foreground)',
				},
				secondary: {
					DEFAULT: 'var(--secondary)',
					foreground: 'var(--secondary-foreground)',
				},
				destructive: {
					DEFAULT: 'var(--destructive)',
					foreground: 'var(--destructive-foreground)',
				},
				info: {
					DEFAULT: 'var(--info)',
					foreground: 'var(--info-foreground)',
				},
				success: {
					DEFAULT: 'var(--success)',
					foreground: 'var(--success-foreground)',
				},
				warning: {
					DEFAULT: 'var(--warning)',
					foreground: 'var(--warning-foreground)',
				},
				error: {
					DEFAULT: 'var(--error)',
					foreground: 'var(--error-foreground)',
				},
				muted: {
					DEFAULT: 'var(--muted)',
					foreground: 'var(--muted-foreground)',
				},
				accent: {
					DEFAULT: 'var(--accent)',
					foreground: 'var(--accent-foreground)',
				},
				popover: {
					DEFAULT: 'var(--popover)',
					foreground: 'var(--popover-foreground)',
				},
				card: {
					DEFAULT: 'var(--card)',
					foreground: 'var(--card-foreground)',
				},
			},
			borderRadius: {
				xl: 'calc(var(--radius) + 4px)',
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)',
			},
			keyframes: {
				'accordion-down': {
					from: { height: 0 },
					to: { height: 'var(--kb-accordion-content-height)' },
				},
				'accordion-up': {
					from: { height: 'var(--kb-accordion-content-height)' },
					to: { height: 0 },
				},
				'content-show': {
					from: { opacity: 0, transform: 'scale(0.96)' },
					to: { opacity: 1, transform: 'scale(1)' },
				},
				'content-hide': {
					from: { opacity: 1, transform: 'scale(1)' },
					to: { opacity: 0, transform: 'scale(0.96)' },
				},
				'caret-blink': {
					'0%,70%,100%': { opacity: '1' },
					'20%,50%': { opacity: '0' },
				},
			},
			animation: {
				'accordion-down': 'accordion-down 0.2s ease-out',
				'accordion-up': 'accordion-up 0.2s ease-out',
				'content-show': 'content-show 0.2s ease-out',
				'content-hide': 'content-hide 0.2s ease-out',
				'caret-blink': 'caret-blink 1.25s ease-out infinite',
			},
		},
	},
	plugins: [animate],
}
