/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    fontSize: {
      '3xl': [
        '27px',
        {
          lineHeight: '32px',
          fontWeight: '800',
        },
      ],
      '2xl': [
        '20px',
        {
          lineHeight: '25px',
          fontWeight: '500',
        },
      ],
      xl: [
        '17px',
        {
          lineHeight: '23px',
          fontWeight: '500',
        },
      ],
      lg: [
        '15px',
        {
          lineHeight: '22px',
          fontWeight: '500',
        },
      ],
      base: [
        '15px',
        {
          lineHeight: '22px',
          fontWeight: '400',
        },
      ],
      sm: [
        '13px',
        {
          lineHeight: '18px',
          fontWeight: '400',
        },
      ],
      xs: [
        '11px',
        {
          lineHeight: '18px',
          fontWeight: '400',
        },
      ],
    },
    extend: {
      textColor: {
        main: 'var(--telegram-text-color)',
        secondary: 'var(--telegram-subtitle-text-color)',
        hint: 'var(--telegram-hint-color)',
        link: 'var(--telegram-link-color)',
        button: 'var(--telegram-button-text-color)',
        accent: 'var(--telegram-accent-text-color)',
        destructive: 'var(--telegram-destructive-text-color)',
        'section-header': 'var(--telegram-section-header-text-color)',
        pink: '#EF5DA8',
        red: '#FE5F55',
        green: '#408F1B',
        orange: '#FF8C42',
        blue: '#3478F6',
      },
      colors: {
        pink: '#EF5DA8',
        red: '#FE5F55',
        orange: '#FF8C42',
        blue: '#3478F6',
        'peatch-hint': 'var(--telegram-hint-color)',
        'peatch-accent': 'var(--telegram-accent-color)',
        'peatch-main': 'var(--telegram-background-color)',
        'peatch-button': 'var(--telegram-button-color)',
        'peatch-secondary': 'var(--telegram-secondary-bg-color)',
        'peatch-header-bg': 'var(--telegram-header-bg-color)',
        'peatch-section-bg': 'var(--telegram-section-bg-color)',
      },
    },
  },
  plugins: [],
};
