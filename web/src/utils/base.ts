/** 安全入口前缀，由服务端注入 index.html */
declare global {
  interface Window {
    __ZPANEL_ENTRY__?: string
  }
}

/** 应用路由 base，如 / 或 /abc123/ */
export function getAppBase(): string {
  const e = (window.__ZPANEL_ENTRY__ || '').replace(/\/$/, '')
  return e ? `${e}/` : '/'
}

/** API base，如 /api/v1 或 /abc123/api/v1 */
export function getAPIBase(): string {
  const e = (window.__ZPANEL_ENTRY__ || '').replace(/\/$/, '')
  return e ? `${e}/api/v1` : '/api/v1'
}
