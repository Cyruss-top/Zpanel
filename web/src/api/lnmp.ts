import { request } from './client'

export interface ComponentStatus {
  installed: boolean
  version: string
  running: boolean
}

export interface LNMPStatus {
  installed: boolean
  platform: string
  components: Record<string, ComponentStatus>
}

export function fetchLNMPStatus() {
  return request<LNMPStatus>('/lnmp/status')
}

export function installLNMP() {
  return request<{ ok: boolean; nginx: string; php: string; mysql: string }>('/lnmp/install', {
    method: 'POST',
  })
}
