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
        // for testin purposes
        bg_color: '#ffffff',
        text_color: '#000000',
        hint_color: '#707579',
        link_color: '#3390ec',
        button_color: '#3390ec',
        button_text_color: '#ffffff',
        secondary_bg_color: '#f4f4f5',
        header_bg_color: '#ffffff',
        accent_text_color: '#3390ec',
        section_bg_color: '#ffffff',
        section_header_text_color: '#707579',
        subtitle_text_color: '#707579',
        destructive_text_color: '#e53935',
        // for testin purposes
        bg_color_dark: '#212121',
        text_color_dark: '#ffffff',
        hint_color_dark: '#aaaaaa',
        link_color_dark: '#8774e1',
        button_color_dark: '#8774e1',
        button_text_color_dark: '#ffffff',
        secondary_bg_color_dark: '#0f0f0f',
        header_bg_color_dark: '#212121',
        accent_text_color_dark: '#8774e1',
        section_bg_color_dark: '#212121',
        section_header_text_color_dark: '#aaaaaa',
        subtitle_text_color_dark: '#aaaaaa',
        destructive_text_color_dark: '#e53935',

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
        hint: 'var(--telegram-hint-color)',
        accent: 'var(--telegram-accent-text-color)',
        main: 'var(--telegram-secondary-bg-color)',
        secondary: 'var(--telegram-section-bg-color)',
        button: 'var(--telegram-button-color)',
        'header-bg': 'var(--telegram-header-bg-color)',
        'section-bg': 'var(--telegram-section-bg-color)',
      },
    },
  },
  plugins: [],
};
