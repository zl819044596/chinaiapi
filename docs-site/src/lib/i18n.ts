import { defineI18n } from 'fumadocs-core/i18n';

const _i18n = defineI18n({
  defaultLanguage: 'en',
  languages: ['en', 'zh', 'ja'],
  parser: 'dir',
});

// Plain config object (no functions) safe for serialization to client components
export const i18n = {
  defaultLanguage: _i18n.defaultLanguage,
  languages: _i18n.languages,
  parser: _i18n.parser,
} as const;

// Factory for the new fumadocs-ui i18nProvider API — server-side only
export function createTranslations() {
  return _i18n.translations();
}

export function getLocalePath(lang: string, path = ''): string {
  const cleanPath = path.startsWith('/') ? path.slice(1) : path;
  return cleanPath ? `/${lang}/${cleanPath}` : `/${lang}`;
}
