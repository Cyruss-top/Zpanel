import { getAPIBase } from '@/utils/base'

export interface APIResponse<T = unknown> {
  ok: boolean
  message: string
  code?: string
  data?: T
}

export class APIError extends Error {
  code?: string
  constructor(message: string, code?: string) {
    super(message)
    this.code = code
  }
}

const BASE = getAPIBase()

type RouterLike = { push: (path: string) => void }

let routerRef: RouterLike | null = null

export function setAPIRouter(router: RouterLike) {
  routerRef = router
}

export async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const token = localStorage.getItem('token')
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options?.headers,
    },
  })

  if (res.status === 401) {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    routerRef?.push('/login')
    throw new APIError('Unauthorized', 'UNAUTHORIZED')
  }

  const data: APIResponse<T> = await res.json()
  if (!data.ok) {
    throw new APIError(data.message || 'request failed', data.code)
  }
  return data.data as T
}

export async function login(username: string, password: string) {
  return request<{ token: string; username: string; expires: string }>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}
