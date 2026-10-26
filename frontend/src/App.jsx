import { useEffect, useState } from 'react'
import Home from './pages/Home.jsx'
import RoutesPage from './pages/Routes.jsx'
import FormatPage from './pages/Format.jsx'
import TryPage from './pages/Try.jsx'

const pages = {
  '/': Home,
  '/routes': RoutesPage,
  '/format': FormatPage,
  '/try': TryPage,
}

function pathFromLocation() {
  const path = window.location.pathname
  return pages[path] ? path : '/'
}

export default function App() {
  const [path, setPath] = useState(pathFromLocation)
  const [theme, setTheme] = useState(() => document.documentElement.dataset.theme || 'light')

  function toggleTheme() {
    const next = theme === 'dark' ? 'light' : 'dark'
    document.documentElement.dataset.theme = next
    localStorage.setItem('veyl-theme', next)
    setTheme(next)
  }

  useEffect(() => {
    const onPop = () => setPath(pathFromLocation())
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  function go(event, next) {
    event.preventDefault()
    if (next === path) return
    window.history.pushState(null, '', next)
    setPath(next)
    window.scrollTo(0, 0)
  }

  const Page = pages[path] ?? Home

  return (
    <>
      <header className="bar">
        <a className="wordmark" href="/" onClick={(event) => go(event, '/')}>
          veyl
        </a>
        <div className="bar-links">
          <nav>
            <a href="/routes" aria-current={path === '/routes' ? 'page' : undefined} onClick={(event) => go(event, '/routes')}>
              Routes
            </a>
            <a href="/format" aria-current={path === '/format' ? 'page' : undefined} onClick={(event) => go(event, '/format')}>
              Format
            </a>
            <a href="/try" aria-current={path === '/try' ? 'page' : undefined} onClick={(event) => go(event, '/try')}>
              Try
            </a>
          </nav>
          <button
            type="button"
            className="theme-toggle"
            aria-pressed={theme === 'dark'}
            onClick={toggleTheme}
          >
            {theme === 'dark' ? 'Light' : 'Dark'}
          </button>
        </div>
      </header>
      <Page go={go} />
    </>
  )
}
