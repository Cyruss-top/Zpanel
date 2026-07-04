import { request } from './client'

export interface MonitorOverview {
  cpu: { usage_percent: number; cores: number }
  memory: {
    total: number
    used: number
    available: number
    used_percent: number
  }
  disk: {
    total: number
    used: number
    available: number
    used_percent: number
    path: string
  }
  load: { load1: number; load5: number; load15: number }
  uptime_seconds: number
  os: string
  arch: string
  platform: string
}

export function fetchOverview() {
  return request<MonitorOverview>('/monitor/overview')
}
