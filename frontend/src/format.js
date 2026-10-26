const example = {
  name: 'John Doe',
  tag: 'BOBA',
  rank: 'Diamond 2',
  rr: 88,
  lastRrChange: '+17',
  elo: 1688,
  peak: 'Ascendant 1',
  wins: 3,
  losses: 1,
  games: 4,
  rrChange: '+42',
}

const regions = ['na', 'eu', 'ap', 'kr', 'latam', 'br']

const presets = {
  rank: '{name}#{tag} → {rank} {rr} RR',
  record: '{name}#{tag} is {wins}-{losses} this stream ({rrChange} RR).',
  both: '{name}#{tag} is {rank} {rr} RR. Stream record {wins}-{losses}.',
}

function signed(value) {
  if (typeof value !== 'number') return value
  return value > 0 ? `+${value}` : String(value)
}

export function formatMessage(template, data) {
  return template.replace(/\{([a-zA-Z]+)(?::(lower|upper|num))?\}/g, (_, key, mod) => {
    if (!(key in data)) return `{${key}}`
    const value = data[key]
    if (mod === 'num') return String(value).replace(/[^\d-]/g, '')
    if (mod === 'lower') return String(value).toLowerCase()
    if (mod === 'upper') return String(value).toUpperCase()
    if (key === 'wins') return `${value}W`
    if (key === 'losses') return `${value}L`
    return String(value)
  })
}

function lineUrl(origin, identity, template) {
  const params = new URLSearchParams({
    name: identity.name.trim() || 'NAME',
    tag: identity.tag.trim() || 'TAG',
    region: identity.region,
    platform: identity.platform,
  })
  if (template.trim()) params.set('template', template.trim())
  const base = (origin || '').replace(/\/$/, '')
  return `${base}/api/valorant/line?${params.toString()}`
}

function commandFor(bot, url, commandName) {
  const name = commandName.trim().startsWith('!') ? commandName.trim() : `!${commandName.trim() || 'rank'}`
  if (bot === 'fossabot') return `!addcom ${name} $(customapi ${url})`
  if (bot === 'streamelements') return `!addcom ${name} $(customapi ${url})`
  return `!addcom ${name} $(urlfetch ${url})`
}

export { example, regions, presets, signed, lineUrl, commandFor }
