// Load runtime config. It first tries to fetch /env-config.json (can be
// generated at deploy time from .env). If that fails, it falls back to
// `window.API_DOMAIN` (injected by server/container) and finally to
// `window.location.origin`.

export async function loadConfig() {
  try {
    const resp = await fetch('/env-config.json', { cache: 'no-store' });
    if (resp.ok) {
      const j = await resp.json();
      return { API_DOMAIN: j.API_DOMAIN || window.API_DOMAIN || window.location.origin };
    }
  } catch (e) {
    // ignore and fallback
  }

  return { API_DOMAIN: window.API_DOMAIN || window.location.origin };
}

export default loadConfig;
