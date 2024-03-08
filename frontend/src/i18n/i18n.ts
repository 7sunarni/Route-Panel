import { derived, writable } from "svelte/store";
import translations from "./translations";

export const locale = writable("en");
export const locales: string[] = Object.keys(translations);

function translate(locale: string, key: string, vars: string) {
  // Let's throw some errors if we're trying to use keys/locales that don't exist.
  // We could improve this by using Typescript and/or fallback values.
  if (!key) throw new Error("no key provided to $18n()");
  if (!locale) throw new Error(`no translation for key "${key}"`);

  // Grab the translation from the translations object.
  let text: string = translations[locale][key];

  if (!text) throw new Error(`no translation found for ${locale}.${key}`);

  // Replace any passed in variables in the translation string.
  Object.keys(vars).map((k) => {
    const regex = new RegExp(`{{${k}}}`, "g");
    text = text.replace(regex, vars[k]);
  });

  return text;
}

export const i18n = derived(locale, ($locale) => (key: string) =>
  translate($locale, key, "")
);
