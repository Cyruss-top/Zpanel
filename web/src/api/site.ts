import { request } from './client'

export type SiteType = 'html' | 'php' | 'go'

export interface Site {
  id: string
  name: string
  type: SiteType
  status: string
  domains: string[]
  root: string
  php_version?: string
  go_port?: number
  go_binary?: string
  nginx_config_path: string
  created_at: string
  updated_at: string
}

export interface CreateSitePayload {
  name: string
  type: SiteType
  domains: string[]
  php_version?: string
  go_port?: number
  go_binary?: string
}

export function fetchSites() {
  return request<Site[]>('/sites')
}

export function createSite(payload: CreateSitePayload) {
  return request<Site>('/sites', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function deleteSite(id: string) {
  return request<{ deleted: boolean }>(`/sites/${id}`, { method: 'DELETE' })
}
