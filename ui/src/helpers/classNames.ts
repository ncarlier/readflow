export function classNames(...names: (string | undefined | null)[]) {
  return names.filter((name) => name).join(' ')
}
