const units = {
  second: 60,
  minute: 60,
  hour: 24,
  day: 7,
  week: 4.35,
  month: 12,
  year: 10000
}

export function ago(date: Date) {
  let delta = Math.floor((Date.now() - date.getTime()) / 1000)
  if (delta == 0) {
    return 'now'
  }
  for (const key in units) {
    const unit = units[key as keyof typeof units]
    const value = delta % unit
    if (!(delta = 0 | delta / unit)) {
        if (value < 0) {
          return `in ${Math.abs(value)} ${key}${value < -1 ? 's' : ''}`
        } else {
          return `${value} ${key}${value > 1 ? 's' : ''} ago`
        }
    }
  }
  return 'undefined'
}