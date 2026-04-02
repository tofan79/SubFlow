import flowbitePlugin from 'flowbite/plugin';

/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'
	],
	darkMode: 'class',
	theme: {
		extend: {
			fontFamily: {
				sans: ['Inter', 'sans-serif'],
				headline: ['Sora', 'sans-serif'],
				label: ['Space Grotesk', 'monospace']
			},
			colors: {
				primary: {
					50: '#fff0f5',
					100: '#ffe0eb',
					200: '#ffc2d7',
					300: '#ff8fb4',
					400: '#ff5c91',
					500: '#ff2d78',
					600: '#e01a5f',
					700: '#bc1250',
					800: '#9a1345',
					900: '#82143d',
					DEFAULT: '#ff2d78',
					container: '#b3004e',
					fixed: '#ffe0ec',
					'fixed-dim': '#ff80aa'
				},
				secondary: {
					50: '#e6fffc',
					100: '#ccfff9',
					200: '#99fff3',
					300: '#66ffed',
					400: '#33ffe6',
					500: '#00ffcc',
					600: '#00cca3',
					700: '#00997a',
					800: '#006652',
					900: '#003329',
					DEFAULT: '#00ffcc',
					container: '#004d3d',
					fixed: '#c0fff4',
					'fixed-dim': '#00e6b8'
				},
				tertiary: {
					50: '#fffceb',
					100: '#fff8d6',
					200: '#fff1ad',
					300: '#ffeb85',
					400: '#ffe55c',
					500: '#ffe04a',
					600: '#ccb33b',
					700: '#99862c',
					800: '#66591e',
					900: '#332c0f',
					DEFAULT: '#ffe04a',
					container: '#665200',
					fixed: '#fff0c0',
					'fixed-dim': '#ffe04a'
				},
				surface: {
					background: '#0a0a12',
					DEFAULT: '#0f0f1a',
					container: '#141422',
					elevated: '#1a1a2e',
					border: '#2a2a3e',
					dim: '#0f0f1a',
					bright: '#1a1a2e',
					'container-lowest': '#0a0a12',
					'container-low': '#111118',
					'container-high': '#1e1e30',
					'container-highest': '#28283e',
					variant: '#1e1e30',
					tint: '#ff2d78'
				},
				background: '#0a0a12',
				'on-background': '#e8e0f0',
				error: '#ff4444',
				'error-container': '#3d0f0f',
				'on-error': '#1a0000',
				'on-error-container': '#ffa0a0',
				outline: '#5a5068',
				'outline-variant': '#302840',
				'on-surface': '#e8e0f0',
				'on-surface-variant': '#a098b0',
				'inverse-surface': '#e8e0f0',
				'inverse-on-surface': '#0a0a12',
				'inverse-primary': '#8c0038',
				'on-primary': '#1a0010',
				'on-primary-container': '#ffe0ec',
				'on-primary-fixed': '#3d0020',
				'on-primary-fixed-variant': '#8c0038',
				'on-secondary': '#001a1a',
				'on-secondary-container': '#c0fff4',
				'on-secondary-fixed': '#001a1a',
				'on-secondary-fixed-variant': '#004d4d',
				'on-tertiary': '#1a1000',
				'on-tertiary-container': '#fff0c0',
				'on-tertiary-fixed': '#1a1000',
				'on-tertiary-fixed-variant': '#665200'
			},
			boxShadow: {
				'neon-primary': '0 0 12px rgba(255, 45, 120, 0.4), 0 0 24px rgba(255, 45, 120, 0.2)',
				'neon-secondary': '0 0 12px rgba(0, 255, 204, 0.4), 0 0 24px rgba(0, 255, 204, 0.2)',
				'neon-tertiary': '0 0 12px rgba(255, 224, 74, 0.4), 0 0 24px rgba(255, 224, 74, 0.2)'
			},
			animation: {
				'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
				'glow': 'glow 2s ease-in-out infinite alternate'
			},
			keyframes: {
				glow: {
					'0%': { boxShadow: '0 0 8px rgba(255, 45, 120, 0.4)' },
					'100%': { boxShadow: '0 0 16px rgba(255, 45, 120, 0.6), 0 0 24px rgba(255, 45, 120, 0.3)' }
				}
			}
		}
	},
	plugins: [flowbitePlugin]
};
