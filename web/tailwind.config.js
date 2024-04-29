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
          lineHeight: '32px',
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
        black: '#9932CC',
        gray: '#909092',
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
        'peatch-bg': '#F8F8F8',
        'peatch-black': '#000100',
        'peatch-gray': '#B6B6B6',
        'peatch-blue': '#3F8AF7',
        'peatch-blue-active': '#E9FDFF',
        'peatch-light-black': '#2F302A',
        'peatch-stroke': '#F6F6F6',
        'peatch-blue-inactive': '#BEDDFC',
        'peatch-green': '#5FA95B',
        'peatch-light-green': '#DBFBD9',
        'peatch-white': '#FFFFFF',
        'peatch-brown-stroke': '#F3F1E6',
        'peatch-black-stroke': '#494A45',
        'peatch-light-gray': '#C9C9C9',
        'peatch-dark-gray': '#949494',
      },
    },
  },
  plugins: [],
};
