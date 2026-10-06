import type { Content, Work } from './types';

async function req<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, { credentials: 'same-origin', ...init });
  if (!res.ok) {
    let msg = res.statusText;
    try { msg = (await res.json()).error ?? msg; } catch {}
    throw new Error(msg);
  }
  return res.status === 204 ? (undefined as T) : res.json();
}

const json = (method: string, body: unknown): RequestInit => ({
  method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
});

export const api = {
  content: () => req<Content>('/api/content'),
  session: () => req<{ authed: boolean }>('/api/admin/session'),
  login: (password: string) => req<void>('/api/admin/login', json('POST', { password })),
  logout: () => req<void>('/api/admin/logout', { method: 'POST' }),
  saveContent: (c: Content) => req<Content>('/api/admin/content', json('PUT', c)),
  createWork: (w: Omit<Work, 'id'>) => req<Work>('/api/admin/works', json('POST', w)),
  updateWork: (w: Work) => req<Work>(`/api/admin/works/${w.id}`, json('PUT', w)),
  deleteWork: (id: string) => req<void>(`/api/admin/works/${id}`, { method: 'DELETE' }),
  orderWorks: (ids: string[]) => req<void>('/api/admin/works/order', json('PUT', { ids })),
  upload: (file: File) => {
    const fd = new FormData();
    fd.append('file', file);
    return req<{ url: string; size: string }>('/api/admin/upload', { method: 'POST', body: fd });
  },
};
