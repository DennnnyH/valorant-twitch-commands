import { useState } from 'react'
import { example, formatMessage, presets, regions, signed } from '../format.js'

export default function TryPage() {
  const [name, setName] = useState('')
  const [tag, setTag] = useState('')
  const [region, setRegion] = useState('na')
  const [platform, setPlatform] = useState('pc')
  const [template, setTemplate] = useState(presets.both)
  const [note, setNote] = useState('Load a Riot ID to see the live response.')
  const [busy, setBusy] = useState(false)
  const [data, setData] = useState(null)
  const [raw, setRaw] = useState(null)

  function params() {
    return new URLSearchParams({
      name: name.trim(),
      tag: tag.trim(),
      region,
      platform,
    })
  }

  async function load() {
    if (!name.trim() || !tag.trim()) {
      setNote('Enter a name and a tag.')
      return
    }
    setBusy(true)
    setNote('Loading…')
    const query = params()
    query.set('template', template)
    try {
      const [rankRes, recordRes, lineRes] = await Promise.all([
        fetch(`/api/valorant/rank?${params()}`),
        fetch(`/api/valorant/record?${params()}`),
        fetch(`/api/valorant/line?${query}`),
      ])
      const rank = await rankRes.json()
      const record = recordRes.ok ? await recordRes.json() : null
      const line = await lineRes.text()
      setRaw({ rank, record, line })
      if (!rankRes.ok) {
        setData(null)
        setNote(rank.error || 'Rank lookup failed.')
        return
      }
      setData({
        name: rank.riotId.name,
        tag: rank.riotId.tag,
        rank: rank.current.rankName,
        rr: rank.current.rr,
        lastRrChange: signed(rank.current.lastRrChange),
        elo: rank.current.elo,
        peak: rank.peak.rankName,
        wins: record?.wins ?? 0,
        losses: record?.losses ?? 0,
        games: record?.games ?? 0,
        rrChange: signed(record?.rrChange ?? 0),
      })
      setNote(lineRes.ok ? 'Loaded from the API.' : line)
    } catch {
      setData(null)
      setRaw(null)
      setNote('The API is not running. Start it with go run ./cmd/server.')
    } finally {
      setBusy(false)
    }
  }

  async function resetRecord() {
    setNote('Resetting the stream record…')
    try {
      const res = await fetch('/api/valorant/record/start', { method: 'POST' })
      const body = await res.json()
      if (!res.ok) {
        setNote('Could not reset the stream record.')
        return
      }
      setNote(`Stream record reset at ${body.startedAt}.`)
    } catch {
      setNote('The API is not running. Start it with go run ./cmd/server.')
    }
  }

  const shown = data ?? { ...example, name: name.trim() || example.name, tag: tag.trim() || example.tag }

  return (
    <main className="page">
      <header className="page-head">
        <h1>Try</h1>
        <p>Look up a Riot ID and read the same sentence your chat command will post.</p>
      </header>

      <form
        onSubmit={(event) => {
          event.preventDefault()
          load()
        }}
      >
        <div className="fields">
          <label>
            Name
            <input value={name} onChange={(event) => setName(event.target.value)} autoComplete="off" />
          </label>
          <label>
            Tag
            <input value={tag} onChange={(event) => setTag(event.target.value)} autoComplete="off" />
          </label>
          <label>
            Region
            <select value={region} onChange={(event) => setRegion(event.target.value)}>
              {regions.map((value) => (
                <option key={value} value={value}>{value}</option>
              ))}
            </select>
          </label>
          <label>
            Platform
            <select value={platform} onChange={(event) => setPlatform(event.target.value)}>
              <option value="pc">pc</option>
              <option value="console">console</option>
            </select>
          </label>
        </div>
        <label className="template">
          Sentence
          <textarea value={template} onChange={(event) => setTemplate(event.target.value)} rows={3} />
        </label>
        <div className="actions">
          <button type="submit" disabled={busy}>Load</button>
          <button type="button" onClick={resetRecord}>Reset stream record</button>
        </div>
      </form>
      <p className="note">{note}</p>

      <figure className="chat try-chat">
        <figcaption>{data ? 'From the API' : 'Example, until you load a Riot ID'}</figcaption>
        <p className="chat-line bot">
          <span>bot</span> {formatMessage(template, shown)}
        </p>
      </figure>

      {raw && (
        <div className="raw">
          <h2>Rank</h2>
          <pre>{JSON.stringify(raw.rank, null, 2)}</pre>
          <h2>Record</h2>
          <pre>{JSON.stringify(raw.record, null, 2)}</pre>
          <h2>Line</h2>
          <pre>{raw.line}</pre>
        </div>
      )}
    </main>
  )
}
