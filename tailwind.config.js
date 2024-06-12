/** @type {import('tailwindcss').Config} */
module.exports = {
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        clash: ["Clash Display", "sans-serif"],
      },
      opacity: {
        1: "0.01",
        2.5: "0.025",
        5: "0.05",
        7.5: "0.075",
        15: "0.15",
      },
      width: {
        3.5: "0.875rem",
      },
    },
  },
  content: ["./views/**/*.templ"],
  plugins: [require("tailwindcss-animate"), require("@tailwindcss/forms")],
};
