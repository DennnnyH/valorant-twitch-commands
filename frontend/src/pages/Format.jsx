import { useMemo, useState } from 'react'
import { commandFor, example, formatMessage, lineUrl, presets, regions } from '../format.js'

const bots = [
  ['nightbot', 'Nightbot'],
  ['fossabot', 'Fossabot'],
  ['streamelements', 'StreamElements'],
]

const tokens = Object.keys(example)

export default function FormatPage() {
  const [name, setName] = useState('')
  const [tag, setTag] = useState('')
  const [region, setRegion] = useState('na')
  const [platform, setPlatform] = useState('pc')
  const [origin, setOrigin] = useState(() => window.location.origin)
  const [template, setTemplate] = useState(presets.rank)
  const [bot, setBot] = useState('nightbot')
  const [commandName, setCommandName] = useState('!rank')
  const [copied, setCopied] = useState('')

  const previewData = {
    ...example,
    name: name.trim() || example.name,
    tag: tag.trim() || example.tag,
  }
  const preview = formatMessage(template, previewData)
  const count = [...preview].length
  const url = useMemo(
    () => lineUrl(origin, { name, tag, region, platform }, template),
    [origin, name, tag, region, platform, template],
  )
  const command = commandFor(bot, url, commandName)

  async function copy(value, which) {
    try {
      await navigator.clipboard.writeText(value)
      setCopied(which)
    } catch {
      setCopied('')
    }
  }

  function insertToken(token) {
    setTemplate((current) => `${current}${current.endsWith(' ') || current === '' ? '' : ' '}${token}`)
  }

  return (
    <main className="page format-page">
      <header className="page-head">
        <h1>Format</h1>
        <p>
          Write the sentence, then copy the command into your bot dashboard.
          Replace nothing if the name and tag below are yours. The preview uses sample rank numbers until someone runs the command in chat.
        </p>
      </header>

      <div className="fields">
        <label>
          Name
          <input value={name} onChange={(event) => setName(event.target.value)} autoComplete="off" placeholder="Denny" />
        </label>
        <label>
          Tag
          <input value={tag} onChange={(event) => setTag(event.target.value)} autoComplete="off" placeholder="2005" />
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
        <textarea
          value={template}
          onChange={(event) => setTemplate(event.target.value)}
          rows={3}
          maxLength={240}
        />
      </label>
      <div className="presets">
        <button type="button" onClick={() => setTemplate(presets.rank)}>Rank</button>
        <button type="button" onClick={() => setTemplate(presets.record)}>Record</button>
        <button type="button" onClick={() => setTemplate(presets.both)}>Rank and record</button>
      </div>
      <div className="tokens">
        {tokens.map((key) => (
          <button key={key} type="button" onClick={() => insertToken(`{${key}}`)}>
            {`{${key}}`}
          </button>
        ))}
      </div>
      <label className="origin">
        Public API address
        <input value={origin} onChange={(event) => setOrigin(event.target.value)} spellCheck="false" />
      </label>
      <p className="hint">
        Nightbot cannot call localhost. Put the address where this API is reachable from the internet.
      </p>

      <div className="instrument">
        <section className="panel">
          <div className="panel-head">
            <h2>Live preview</h2>
            <span className={count > 400 ? 'count over' : 'count'}>{count} / 400</span>
          </div>
          <p className="preview-box">
            <span className="who">bot:</span> {preview || 'Write a sentence above.'}
          </p>
        </section>

        <section className="panel">
          <div className="panel-head">
            <h2>Request URL</h2>
          </div>
          <p className="url-box">{url}</p>
          <div className="panel-actions">
            <button type="button" onClick={() => copy(url, 'url')}>
              {copied === 'url' ? 'Copied' : 'Copy URL'}
            </button>
            <a href={url} target="_blank" rel="noreferrer">Open</a>
          </div>
        </section>

        <section className="panel">
          <div className="panel-head">
            <h2>Chatbot command</h2>
          </div>
          <label className="command-name">
            Command
            <input value={commandName} onChange={(event) => setCommandName(event.target.value)} />
          </label>
          <div className="bots" role="tablist">
            {bots.map(([id, label]) => (
              <button
                key={id}
                type="button"
                role="tab"
                aria-selected={bot === id}
                className={bot === id ? 'on' : undefined}
                onClick={() => setBot(id)}
              >
                {label}
              </button>
            ))}
          </div>
          <pre className="command">{command}</pre>
          <div className="panel-actions">
            <button type="button" onClick={() => copy(command, 'command')}>
              {copied === 'command' ? 'Copied' : 'Copy command'}
            </button>
          </div>
          <p className="panel-note">
            Paste into the bot dashboard. Save !rank with the rank sentence, then set the command to !record and save the record sentence as a second command. Nightbot needs the response under 400 characters.
          </p>
        </section>
      </div>
    </main>
  )
}
