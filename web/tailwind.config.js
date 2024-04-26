/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
  ],
  theme: {
    fontSize: {
      '3xl': ['27px', {
        lineHeight: '32px',
        fontWeight: '500',
      }],
      '2xl': ['20px', {
        lineHeight: '25px',
        fontWeight: '500',
      }],
      'xl': ['17px', {
        lineHeight: '32px',
        fontWeight: '500',
      }],
      'lg': ['15px', {
        lineHeight: '22px',
        fontWeight: '500',
      }],
      'base': ['15px', {
        lineHeight: '22px',
        fontWeight: '400',
      }],
      'sm': ['13px', {
        lineHeight: '18px',
        fontWeight: '400',
      }],
    },
    extend: {
      textColor: {
        'black': '#000100',
        'gray': '#909092',
      },
      colors: {
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