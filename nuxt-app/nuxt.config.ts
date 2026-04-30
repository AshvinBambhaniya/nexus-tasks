import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },

  modules: ["@pinia/nuxt", "pinia-plugin-persistedstate/nuxt", "@nuxt/eslint"],

  css: ["~/assets/css/main.css"],

  vite: {
    plugins: [tailwindcss()],
  },

  runtimeConfig: {
    apiUrl: process.env.NUXT_API_URL || "http://backend:8000",
    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || "http://localhost:8000",
    },
  },

  app: {
    head: {
      htmlAttrs: {
        lang: "en",
      },
      title: "Nexus Tasks",
      meta: [
        { name: "description", content: "Unified Task Management System" },
      ],
    },
  },
});
