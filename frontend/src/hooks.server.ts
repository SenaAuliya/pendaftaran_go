// src/hooks.server.ts (atau .js)
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
  const { url } = event;
  
  // Jika rute adalah '/login', tidak ada layout yang diterapkan
  if (url.pathname === '/auth/login') {
    return resolve(event, {
      transformPageChunk: ({ html }) => html
    });
  }

  // Untuk rute lainnya, gunakan layout default
  return resolve(event);
};
