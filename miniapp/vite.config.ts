import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
// import fs from 'fs';
// import path from 'path';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // https: {
    //   key: fs.readFileSync(path.resolve(__dirname, 'certs/server.key')),
    //   cert: fs.readFileSync(path.resolve(__dirname, 'certs/server.crt')),
    // },
    port: 3000,
    host: true,
  },
  define: { 
    'process.env.DEBUG': JSON.stringify(process.env.DEBUG || ''),
  },
});
