export default function Home({ go }) {
  return (
    <main>
      <section className="hero">
        <p className="kicker">Veyl API</p>
        <h1>A chat command for your Valorant rank.</h1>
        <p className="lede">
          Viewers type <span className="cmd">!rank</span> or <span className="cmd">!record</span>.
          Your bot answers with the tier, RR, and the wins and losses since the stream started.
          You paste one line into Nightbot, Fossabot, or StreamElements.
        </p>
        <div className="hero-actions">
          <a className="button solid" href="/format" onClick={(event) => go(event, '/format')}>
            Make the command
          </a>
          <a className="button" href="/routes" onClick={(event) => go(event, '/routes')}>
            See the routes
          </a>
        </div>
      </section>

      <figure className="chat">
        <figcaption>What chat sees</figcaption>
        <p className="chat-line">
          <span>viewer</span> !rank
        </p>
        <p className="chat-line bot">
          <span>bot</span> Denny#2005 → Diamond 2 88 RR
        </p>
      </figure>

      <section className="steps">
        <h2>Three steps, then it stays in chat.</h2>
        <ol>
          <li>
            <strong>Write the sentence.</strong>
            <p>Pick the words chat should read. Drop in rank, RR, peak, or the stream record.</p>
          </li>
          <li>
            <strong>Copy the command.</strong>
            <p>The formatter fills your Riot ID and prints the exact line for your bot dashboard.</p>
          </li>
          <li>
            <strong>Paste it once.</strong>
            <p>Nightbot, Fossabot, and StreamElements fetch the line and post it under 400 characters.</p>
          </li>
        </ol>
      </section>

      <section className="routes-tease">
        <h2>The endpoints behind the commands.</h2>
        <ul>
          <li>
            <code>GET /api/valorant/line</code>
            <p>Plain text. This is the URL your bot should call.</p>
          </li>
          <li>
            <code>GET /api/valorant/rank</code>
            <p>Current tier, RR, last change, elo, and peak, as JSON.</p>
          </li>
          <li>
            <code>GET /api/valorant/record</code>
            <p>Wins, losses, and RR change since you started the stream.</p>
          </li>
        </ul>
        <a href="/try" onClick={(event) => go(event, '/try')}>Try a Riot ID</a>
      </section>
    </main>
  )
}
