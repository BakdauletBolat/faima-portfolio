export interface Stat { label: string; value: string }
export interface Profile {
  name: string; tagline: string; available: boolean;
  stats: Stat[]; heroImage: string; resumePdf: string;
}
export interface Experience { period: string; role: string; company: string; desc: string }
export interface Tool { name: string; level: number }
export type Cat = 'home' | 'house' | 'commercial';
export interface Work {
  id: string; cat: Cat; type: string; title: string; area: string; year: string;
  scope: string; cover: string; pdf: string; size: string;
}
export interface Contacts { email: string; telegram: string; phone: string }
export interface Content {
  profile: Profile; experience: Experience[]; tools: Tool[]; skills: string[];
  works: Work[]; contacts: Contacts;
}

export const CATS: { key: 'all' | Cat; label: string }[] = [
  { key: 'all', label: 'Все' },
  { key: 'home', label: 'Квартиры' },
  { key: 'house', label: 'Дома' },
  { key: 'commercial', label: 'Коммерческие' },
];

export const LEVELS = ['', 'Начальный', 'Базовый', 'Уверенный', 'Продвинутый', 'Эксперт'];
