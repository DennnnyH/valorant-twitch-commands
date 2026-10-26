const routes = [
  {
    method: 'GET',
    path: '/api/valorant/line',
    summary: 'The chat line. Nightbot, Fossabot, and StreamElements should call this.',
    params: [
      ['name', 'Riot name, without the tag'],
      ['tag', 'Riot tag, without the #'],
      ['region', 'na, eu, ap, kr, latam, or br'],
      ['platform', 'pc or console. Defaults to pc'],
      ['template', 'Optional sentence. Defaults to {name}#{tag} → {rank} {rr} RR'],
    ],
    returns: 'text/plain, one sentence, meant to stay under 400 characters.',
    example: 'Denny#2005 → Diamond 2 88 RR',
  },
  {
    method: 'GET',
    path: '/api/valorant/rank',
    summary: 'Current rank as JSON, when you want the fields yourself.',
    params: [
      ['name', 'Riot name'],
      ['tag', 'Riot tag'],
      ['region', 'na, eu, ap, kr, latam, or br'],
      ['platform', 'pc or console. Defaults to pc'],
    ],
    returns: 'JSON with riotId, current rank and RR, peak, and the season record.',
    example: `{
  "riotId": { "name": "Denny", "tag": "2005" },
  "current": { "rankName": "Diamond 2", "rr": 88, "lastRrChange": 17 },
  "peak": { "rankName": "Ascendant 1" }
}`,
  },
  {
    method: 'GET',
    path: '/api/valorant/record',
    summary: 'Wins and losses since the stream record was started.',
    params: [
      ['name', 'Riot name'],
      ['tag', 'Riot tag'],
      ['region', 'na, eu, ap, kr, latam, or br'],
      ['platform', 'pc or console. Defaults to pc'],
    ],
    returns: 'JSON with wins, losses, games, rrChange, and the matches in the window.',
    example: `{
  "wins": 3,
  "losses": 1,
  "games": 4,
  "rrChange": 42
}`,
  },
  {
    method: 'POST',
    path: '/api/valorant/record/start',
    summary: 'Starts the stream record over. Call this when you go live.',
    params: [],
    returns: 'JSON with startedAt, an RFC3339 timestamp.',
    example: `{ "startedAt": "2026-10-24T05:00:00Z" }`,
  },
]

export default function RoutesPage() {
  return (
    <main className="page">
      <header className="page-head">
        <h1>Routes</h1>
        <p>
          Every lookup takes a Riot ID split into name and tag, plus a region.
          Keys stay on the server. Chat only sees the sentence you wrote.
        </p>
      </header>
      <div className="route-list">
        {routes.map((route) => (
          <article key={route.path} className="route">
            <h2>
              <span className={route.method === 'POST' ? 'method post' : 'method'}>{route.method}</span>
              <code>{route.path}</code>
            </h2>
            <p>{route.summary}</p>
            {route.params.length > 0 && (
              <table>
                <thead>
                  <tr>
                    <th>Query</th>
                    <th>What to send</th>
                  </tr>
                </thead>
                <tbody>
                  {route.params.map(([name, detail]) => (
                    <tr key={name}>
                      <td><code>{name}</code></td>
                      <td>{detail}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
            <p className="returns">{route.returns}</p>
            <pre>{route.example}</pre>
          </article>
        ))}
      </div>
      <section className="token-note">
        <h2>Tokens in a template</h2>
        <p>
          <code>{'{name}'}</code> <code>{'{tag}'}</code> <code>{'{rank}'}</code> <code>{'{rr}'}</code>{' '}
          <code>{'{lastRrChange}'}</code> <code>{'{elo}'}</code> <code>{'{peak}'}</code>{' '}
          <code>{'{wins}'}</code> <code>{'{losses}'}</code> <code>{'{games}'}</code> <code>{'{rrChange}'}</code>
        </p>
        <p>
          <code>{'{wins}'}</code> prints 3W. <code>{'{wins:num}'}</code> prints 3.
          Add <code>:lower</code> or <code>:upper</code> to change case, as in <code>{'{rank:lower}'}</code>.
        </p>
      </section>
    </main>
  )
}
