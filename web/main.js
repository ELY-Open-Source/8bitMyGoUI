// Feeds browser input to the Go gallery; blits its frames.
'use strict'

const screen = document.getElementById('screen')
const ctx = screen.getContext('2d')
const ime = document.getElementById('ime')
const blocks = [...document.querySelectorAll('#bar i')]
const params = new URLSearchParams(location.search)
const mac = /Mac|iPhone|iPad/.test(navigator.platform)
const cursors = ['default', 'pointer', 'text', 'move', 'ew-resize', 'ns-resize', 'nwse-resize', 'nesw-resize',
  'not-allowed', 'crosshair', 'grab', 'grabbing', 'n-resize', 'e-resize', 's-resize', 'w-resize', 'col-resize',
  'row-resize', 'vertical-text', 'copy', 'alias', 'context-menu', 'none']
const SHIFT = 1, CTRL = 2, ALT = 4, SUPER = 8
const log = []
const asked = params.get('zoom')
const fixed = /^[1-9]$/.test(asked) ? Number(asked) : 0
let api, zoom, scale, timer, first, active = false, skipText = ''

function note (msg) {
  log.push(msg)
  console.info('8bit:', msg)
}

window.__8bitCrash = (msg) => {
  console.error('8bit: crash', msg)
  clearTimeout(timer)
  document.getElementById('why').textContent = msg
  document.getElementById('over').hidden = false
}
document.getElementById('continue').onclick = () => location.reload()
if (asked !== null && !fixed) note(`ignored zoom=${asked}: use a whole number 1–9`)

// size picks device pixels per dot; returns window in dots.
function size () {
  zoom = fixed || Math.max(1, Math.round(devicePixelRatio * (innerWidth >= 1000 ? 2 : 1)))
  scale = zoom / devicePixelRatio
  return [Math.max(1, Math.floor(innerWidth / scale)), Math.max(1, Math.floor(innerHeight / scale))]
}

function hashPath () { return decodeURIComponent(location.hash.slice(1)) }

// show blits a frame and acts on Go's requests.
function show (r) {
  if (!r) return
  if (screen.width !== r.w || screen.height !== r.h) {
    screen.width = r.w
    screen.height = r.h
  }
  screen.style.width = r.w * scale + 'px'
  screen.style.height = r.h * scale + 'px'
  ctx.putImageData(new ImageData(r.pix, r.w, r.h), 0, 0)
  if (r.path && '#' + r.path !== location.hash && !(r.path === '/' && location.hash === '')) {
    history.pushState(null, '', '#' + r.path)
  }
  screen.style.cursor = cursors[r.cursor]
  active = r.caret !== undefined
  if (active) {
    ime.style.left = r.caret[0] * scale + 'px'
    ime.style.top = r.caret[1] * scale + 'px'
    ime.style.height = r.caret[2] * scale + 'px'
  }
  if (r.copied !== undefined && navigator.clipboard) navigator.clipboard.writeText(r.copied).catch(e => note('clipboard: ' + e))
  for (const url of r.open || []) window.open(url, '_blank', 'noopener')
  clearTimeout(timer)
  if (r.due !== undefined && !document.hidden) timer = setTimeout(() => show(api.frame()), Math.max(0, r.due))
}

function mods (e) {
  const m = (e.shiftKey ? SHIFT : 0) | (e.altKey ? ALT : 0) | (e.ctrlKey ? CTRL : 0) | (e.metaKey ? SUPER : 0)
  // wasm's ui.Cmd is Ctrl; ⌘ stays Super for line edits.
  if (mac && e.metaKey && !/^(Arrow|Backspace|Delete)/.test(e.code)) return (m & ~SUPER) | CTRL
  return m
}

function at (e) {
  const r = screen.getBoundingClientRect()
  return [(e.clientX - r.left) / scale, (e.clientY - r.top) / scale]
}

// Taps click; held fingers drag; moving fingers scroll.
let touch = null
screen.addEventListener('pointerdown', (e) => {
  e.preventDefault()
  if (!e.isPrimary) return
  ime.focus({ preventScroll: true })
  screen.setPointerCapture(e.pointerId)
  const [x, y] = at(e)
  if (e.pointerType === 'touch') {
    const t = touch = { x, y, drag: false, scroll: false }
    t.hold = setTimeout(() => { t.drag = true; show(api.pointer('down', x, y)) }, 400)
  } else if (e.button === 0) show(api.pointer('down', x, y))
  else if (e.button === 2) show(api.pointer('right', x, y))
})
screen.addEventListener('pointermove', (e) => {
  if (!e.isPrimary) return
  const [x, y] = at(e)
  if (e.pointerType !== 'touch' || touch?.drag) return show(api.pointer('move', x, y))
  if (!touch || (!touch.scroll && Math.hypot(x - touch.x, y - touch.y) * scale <= 8)) return
  clearTimeout(touch.hold)
  touch.scroll = true
  show(api.wheel(x, y, touch.x - x, touch.y - y))
  touch.x = x
  touch.y = y
})
screen.addEventListener('pointerup', (e) => {
  if (!e.isPrimary) return
  const [x, y] = at(e)
  if (e.pointerType !== 'touch') {
    if (e.button === 0) show(api.pointer('up', x, y))
    return
  }
  if (!touch) return
  clearTimeout(touch.hold)
  if (touch.drag) show(api.pointer('up', x, y))
  else if (!touch.scroll) {
    api.pointer('down', x, y)
    show(api.pointer('up', x, y))
  }
  touch = null
})
// A release outside the window cancels the press.
screen.addEventListener('pointercancel', (e) => {
  if (!touch || !e.isPrimary) return
  clearTimeout(touch.hold)
  if (touch.drag) show(api.pointer('up', -1, -1))
  touch = null
})
screen.addEventListener('contextmenu', (e) => e.preventDefault())
screen.addEventListener('wheel', (e) => {
  e.preventDefault()
  const unit = e.deltaMode === 1 ? 16 : e.deltaMode === 2 ? innerHeight : 1
  const [x, y] = at(e)
  show(api.wheel(x, y, e.deltaX * unit / scale, e.deltaY * unit / scale))
}, { passive: false })

const keep = new Set(['Tab', 'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'PageUp', 'PageDown', 'Home', 'End',
  'Backspace', 'Delete', 'Enter', 'NumpadEnter', 'Escape'])
window.addEventListener('keydown', (e) => {
  if (e.isComposing || !api) return
  const m = mods(e)
  const cmd = (m & CTRL) && !(m & ALT)
  if (cmd && e.code === 'KeyV') return
  // Like native hosts: text only to fields focused at keydown.
  const typing = active
  show(api.key(e.code, m))
  if (keep.has(e.code) || (!typing && e.key.length === 1)) e.preventDefault()
})
ime.addEventListener('input', (e) => {
  if (e.isComposing) return
  const text = e.data || ''
  ime.value = ''
  if (!e.inputType.startsWith('insert') || e.inputType === 'insertLineBreak' || e.inputType === 'insertFromPaste') return
  if (text === skipText) {
    skipText = ''
    return
  }
  if (text) show(api.text(text))
})
ime.addEventListener('compositionupdate', (e) => show(api.compose(e.data, [...e.data].length)))
ime.addEventListener('compositionend', (e) => {
  api.compose('', 0)
  skipText = e.data
  setTimeout(() => { skipText = '' }, 0)
  ime.value = ''
  show(api.text(e.data))
})
ime.addEventListener('paste', (e) => {
  e.preventDefault()
  show(api.paste(e.clipboardData.getData('text/plain')))
})
window.addEventListener('focus', () => api && show(api.focus(true)))
window.addEventListener('blur', () => api && show(api.focus(false)))
window.addEventListener('popstate', () => api && show(api.go(hashPath() || '/')))
document.addEventListener('visibilitychange', () => api && !document.hidden && show(api.frame()))
let resizing
window.addEventListener('resize', () => {
  clearTimeout(resizing)
  resizing = setTimeout(() => {
    if (!api) return
    const [w, h] = size()
    show(api.resize(w, h))
  }, 50)
})

// load streams the wasm, filling the bar as bytes arrive.
async function load () {
  const res = await fetch('app.wasm?v=__BUILD__')
  if (!res.ok) throw new Error('app.wasm: HTTP ' + res.status)
  const total = Number(res.headers.get('content-length')) || 0
  const reader = res.body.getReader()
  const parts = []
  let got = 0
  for (;;) {
    const { done, value } = await reader.read()
    if (done) break
    parts.push(value)
    got += value.length
    const lit = total ? Math.floor(got / total * blocks.length) : Math.min(blocks.length, parts.length)
    blocks.forEach((b, i) => b.classList.toggle('on', i < lit))
    document.getElementById('loading').setAttribute('aria-valuenow', lit)
  }
  const bytes = new Uint8Array(got)
  let off = 0
  for (const p of parts) {
    bytes.set(p, off)
    off += p.length
  }
  const go = new Go()
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject)
  const t0 = performance.now()
  go.run(instance).then(() => window.__8bitCrash('The Go program exited.'))
  api = window.__8bitGo
  if (!api) throw new Error('Go did not export __8bitGo')
  note(`wasm ${(got / 1e6).toFixed(1)} MB loaded`)
  const [w, h] = size()
  show(api.start(w, h, hashPath()))
  first = performance.now() - t0
  document.getElementById('loading').remove()
  ime.focus({ preventScroll: true })
  window.__8bit = {
    ready: true,
    stats: () => ({ ...api.stats(), zoom, first }),
    texts: () => api.texts(),
    page: () => api.stats().path,
    find: (label) => api.find(label)?.map(v => v * scale),
    log
  }
  note(`ready ${w}x${h} dots, ${zoom} pixels per dot, first frame ${first.toFixed(0)} ms`)
}

load().catch((err) => window.__8bitCrash(String(err && err.stack || err)))
