import { brotliCompressSync, constants } from 'node:zlib'
import { test, expect } from '@playwright/test'

// pageErrors keeps each page's console and page errors.
const pageErrors = new WeakMap()

test.beforeEach(async ({ page }) => {
  const seen = []
  pageErrors.set(page, seen)
  page.on('console', m => { if (m.type() === 'error') seen.push(m.text()) })
  page.on('pageerror', e => seen.push(e.message))
})

// Every page must end free of console errors.
test.afterEach(async ({ page }) => {
  if (test.info().annotations.some(a => a.type === 'console errors')) return
  expect(pageErrors.get(page)).toEqual([])
})

// open loads the page and returns errors seen so far.
async function open (page, hash = '') {
  await page.goto('/' + hash)
  await page.waitForFunction(() => window.__8bit?.ready)
  return pageErrors.get(page)
}

const texts = page => page.evaluate(() => window.__8bit.texts())
const where = page => page.evaluate(() => window.__8bit.page())

// searched reports whether the page shows results for q.
const searched = async (page, q) => (await texts(page)).some(t => t.endsWith(`results for "${q}"`))

// pixels hashes the canvas, to see that it changed.
const pixels = page => page.evaluate(() => {
  const c = document.getElementById('screen')
  const d = c.getContext('2d').getImageData(0, 0, c.width, c.height).data
  let h = 0
  for (let i = 0; i < d.length; i += 4) h = (h * 31 + d[i] + d[i + 1] * 7 + d[i + 2] * 13) | 0
  return h
})

async function box (page, label) {
  const r = await page.evaluate(l => window.__8bit.find(l), label)
  expect(r, `no element "${label}"`).toBeTruthy()
  return r
}

async function click (page, label) {
  const r = await box(page, label)
  await page.mouse.click(r[0] + r[2] / 2, r[1] + r[3] / 2)
}

// touch drags a finger, holding first when hold > 0.
async function touch (page, from, to, hold = 0) {
  const cdp = await page.context().newCDPSession(page)
  await cdp.send('Input.dispatchTouchEvent', { type: 'touchStart', touchPoints: [{ x: from[0], y: from[1] }] })
  if (hold) await page.waitForTimeout(hold)
  for (let i = 1; i <= 4; i++) {
    const p = { x: from[0] + (to[0] - from[0]) * i / 4, y: from[1] + (to[1] - from[1]) * i / 4 }
    await cdp.send('Input.dispatchTouchEvent', { type: 'touchMove', touchPoints: [p] })
  }
  await cdp.send('Input.dispatchTouchEvent', { type: 'touchEnd', touchPoints: [] })
}

test('boots, starts on Enter, shows home', async ({ page }) => {
  const errors = await open(page)
  expect(await texts(page)).toContain('Press start')
  await page.keyboard.press('Enter')
  expect(await texts(page)).toContain('Foundations')
  expect(errors).toEqual([])
})

test('deep links open pages; unknown paths end the game', async ({ page }) => {
  const errors = await open(page, '#/f/icons')
  expect(await where(page)).toBe('/f/icons')
  expect(await texts(page)).toContain('COPY CODE')
  await page.goto('/#/nowhere')
  await page.waitForFunction(() => window.__8bit?.page() === '/nowhere')
  expect(await texts(page)).toContain('No page at /nowhere')
  expect(errors).toEqual([])
})

test('slash focuses search without typing a slash', async ({ page }) => {
  await open(page, '#/')
  await page.keyboard.press('/')
  await page.keyboard.type('dith')
  expect(await searched(page, 'dith')).toBe(true)
  await page.keyboard.press('Enter')
  expect(await where(page)).toBe('/f/dither')
  expect(page.url()).toContain('#/f/dither')
})

test('Chinese text reaches search', async ({ page }) => {
  await open(page, '#/')
  await page.keyboard.press('/')
  await page.keyboard.insertText('像素')
  expect(await searched(page, '像素')).toBe(true)
})

test('clicks navigate, Back returns, P repaints', async ({ page }) => {
  await open(page, '#/')
  await click(page, 'Foundations')
  expect(await where(page)).toBe('/cat/foundations')
  await click(page, 'Palettes')
  expect(await where(page)).toBe('/f/palettes')
  await page.goBack()
  await page.waitForFunction(() => window.__8bit.page() === '/cat/foundations')
  const before = await pixels(page)
  await page.keyboard.press('p')
  expect(await pixels(page)).not.toBe(before)
})

test('typing in a field stays in the field', async ({ page }) => {
  await open(page, '#/f/editor')
  await click(page, 'Name')
  await page.keyboard.type('p]')
  expect(await where(page)).toBe('/f/editor')
  expect((await texts(page)).filter(t => t === '> Name changed')).toHaveLength(2)
})

test('keyboard copy and paste use the system clipboard', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write'])
  await open(page, '#/f/editor')
  await click(page, 'Name')
  await page.keyboard.type('pixel')
  await page.keyboard.press('ControlOrMeta+A')
  await page.keyboard.press('ControlOrMeta+C')
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe('pixel')
  await page.evaluate(() => navigator.clipboard.writeText('dith'))
  await click(page, 'Search components')
  await page.keyboard.press('ControlOrMeta+V')
  expect(await searched(page, 'dith')).toBe(true)
})

test('an input method composes, then commits', async ({ page }) => {
  await open(page, '#/')
  await page.keyboard.press('/')
  const cdp = await page.context().newCDPSession(page)
  await cdp.send('Input.imeSetComposition', { text: 'xiang', selectionStart: 5, selectionEnd: 5 })
  expect(await searched(page, 'xiang')).toBe(false)
  await cdp.send('Input.insertText', { text: '像' })
  expect(await searched(page, '像')).toBe(true)
})

test('Command or Super with Backspace erases the line', async ({ page }) => {
  await open(page, '#/')
  await page.keyboard.press('/')
  await page.keyboard.type('zz dith')
  await page.keyboard.press('Meta+Backspace')
  await page.keyboard.type('dith')
  expect(await searched(page, 'dith')).toBe(true)
})

test('wheel scrolls the sidebar', async ({ page }) => {
  await open(page, '#/')
  const y = (await box(page, 'Display'))[1]
  await page.mouse.move(100, 500)
  await page.mouse.wheel(0, 100)
  expect((await box(page, 'Display'))[1]).toBeLessThan(y)
})

test('an unmapped key keeps animations running', async ({ page }) => {
  await open(page)
  await page.keyboard.press('Shift')
  const frames = (await page.evaluate(() => window.__8bit.stats())).frames
  await page.waitForTimeout(600)
  expect((await page.evaluate(() => window.__8bit.stats())).frames).toBeGreaterThan(frames)
})

test('copying the same text again reaches the clipboard', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write'])
  await open(page, '#/f/icons')
  await click(page, 'COPY CODE')
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain('icon.')
  await page.evaluate(() => navigator.clipboard.writeText('external'))
  await click(page, 'COPY CODE')
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain('icon.')
})

test('zoom must be a whole number', async ({ page }) => {
  await open(page, '?zoom=1.5#/')
  const s = await page.evaluate(() => window.__8bit.stats())
  expect(Number.isInteger(s.zoom)).toBe(true)
  expect(await page.evaluate(() => window.__8bit.log.some(l => l.includes('zoom=1.5')))).toBe(true)
})

test.describe('touch', () => {
  test.use({ hasTouch: true })

  test('a tap clicks', async ({ page }) => {
    await open(page, '#/')
    const r = await box(page, 'Foundations')
    await page.touchscreen.tap(r[0] + r[2] / 2, r[1] + r[3] / 2)
    expect(await where(page)).toBe('/cat/foundations')
  })

  test('a swipe scrolls', async ({ page }) => {
    await open(page, '#/')
    const y = (await box(page, 'Display'))[1]
    await touch(page, [100, 600], [100, 400])
    expect((await box(page, 'Display'))[1]).toBeLessThan(y)
  })

  test('a canceled hold does not click', async ({ page }) => {
    await open(page, '#/')
    const r = await box(page, 'Foundations')
    const cdp = await page.context().newCDPSession(page)
    await cdp.send('Input.dispatchTouchEvent', { type: 'touchStart', touchPoints: [{ x: r[0] + 20, y: r[1] + 10 }] })
    await page.waitForTimeout(500)
    await cdp.send('Input.dispatchTouchEvent', { type: 'touchCancel', touchPoints: [] })
    expect(await where(page)).toBe('/')
  })

  test('hold, then drag, selects text', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write'])
    await open(page, '#/f/editor')
    await click(page, 'Name')
    await page.keyboard.type('abcdef')
    const r = await box(page, 'Name')
    await touch(page, [r[0] + 4, r[1] + r[3] / 2], [r[0] + 300, r[1] + r[3] / 2], 500)
    await page.keyboard.press('ControlOrMeta+C')
    await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe('abcdef')
  })
})

test('COPY CODE fills the clipboard', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write'])
  await open(page, '#/f/icons')
  await click(page, 'COPY CODE')
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain('icon.')
})

test('resize keeps whole pixels per dot', async ({ page }) => {
  await open(page, '#/')
  await page.setViewportSize({ width: 800, height: 600 })
  await page.waitForFunction(() => window.__8bit.stats().w === Math.floor(innerWidth * devicePixelRatio / window.__8bit.stats().zoom))
  const s = await page.evaluate(() => window.__8bit.stats())
  expect(Number.isInteger(s.zoom)).toBe(true)
})

test('a Go panic shows GAME OVER in pixels; CONTINUE reloads', async ({ page }) => {
  test.info().annotations.push({ type: 'console errors', description: 'the panic is on purpose' })
  await open(page, '#/f/icons')
  await page.evaluate(() => window.__8bitGo.pointer('bogus', 0, 0))
  expect(await texts(page)).toContain('web: unknown pointer kind bogus')
  await expect(page.locator('#over')).toBeHidden()
  await click(page, 'CONTINUE')
  await page.waitForFunction(() => window.__8bit?.ready && window.__8bit.page() === '/f/icons')
  expect(await texts(page)).toContain('COPY CODE')
})

test('resizing while wasm loads is safe', async ({ page }) => {
  const errors = []
  page.on('pageerror', e => errors.push(e.message))
  await page.route('**/a/app.*.wasm', async route => {
    await new Promise(resolve => setTimeout(resolve, 800))
    await route.continue()
  })
  await page.goto('/')
  await page.setViewportSize({ width: 900, height: 700 })
  await page.waitForFunction(() => window.__8bit?.ready)
  expect(errors).toEqual([])
})

test('scripts and wasm have content-addressed, immutable URLs', async ({ page, request }) => {
  await open(page, '#/')
  const urls = await page.evaluate(() => performance.getEntriesByType('resource').map(r => r.name).filter(u => u.includes('/a/')))
  expect(urls.map(u => new URL(u).pathname).sort()).toEqual([
    expect.stringMatching(/^\/a\/app\.[0-9a-f]{12}\.wasm$/),
    expect.stringMatching(/^\/a\/main\.[0-9a-f]{12}\.js$/),
    expect.stringMatching(/^\/a\/wasm_exec\.[0-9a-f]{12}\.js$/)
  ])
  // Cache headers come from Cloudflare, not the local server.
  if (!process.env.BASE_URL) return
  for (const u of urls) expect((await request.head(u)).headers()['cache-control']).toContain('immutable')
  expect((await request.head('/')).headers()['cache-control']).toContain('must-revalidate')
})

test('a load failure shows the page GAME OVER', async ({ page }) => {
  test.info().annotations.push({ type: 'console errors', description: 'the 404 is on purpose' })
  await page.route('**/a/app.*.wasm', route => route.fulfill({ status: 404 }))
  await page.goto('/')
  await expect(page.locator('#over')).toBeVisible()
  await expect(page.locator('#why')).toContainText(/a\/app\.[0-9a-f]{12}\.wasm: HTTP 404/)
})

test('components take clicks and typing on their pages', async ({ page }) => {
  const errors = await open(page, '#/c/066-button')
  // The demo's log panel shows each event as "> event".
  const logged = s => page.waitForFunction(l => window.__8bit.texts().includes(l), '> ' + s)
  await click(page, 'PRIMARY')
  await logged('PRIMARY')
  await page.goto('/#/c/082-textinput')
  await page.waitForFunction(() => window.__8bit.page() === '/c/082-textinput')
  await click(page, 'Hero name')
  await page.keyboard.type('Zed')
  await logged('name Zed')
  await page.keyboard.press('Enter')
  await logged('submitted Zed')
  expect(errors).toEqual([])
})

test('the budget: first frame, p95 frame, wasm with brotli', async ({ page, request }) => {
  const errors = await open(page, '#/c/295-starfield')
  // Stars rush at 12 fps; sixty frames fill the window.
  await page.waitForFunction(() => window.__8bit.stats().frames >= 60)
  const s = await page.evaluate(() => window.__8bit.stats())
  const [wasm] = await page.evaluate(() => performance.getEntriesByType('resource')
    .filter(e => e.name.endsWith('.wasm')).map(e => ({ url: e.name, sent: e.encodedBodySize, size: e.decodedBodySize })))
  let br = wasm.sent
  if (br >= wasm.size) {
    // Served plain: compress as Cloudflare's brotli does.
    const body = await (await request.get(wasm.url)).body()
    br = brotliCompressSync(body, { params: { [constants.BROTLI_PARAM_QUALITY]: 4 } }).length
  }
  console.log(`budget: first frame ${s.first.toFixed(0)} ms, p50 ${s.p50.toFixed(1)} ms, p95 ${s.p95.toFixed(1)} ms, wasm ${(br / 1e6).toFixed(2)} MB brotli`)
  expect(s.first).toBeLessThanOrEqual(300)
  expect(s.p95).toBeLessThanOrEqual(16)
  expect(br).toBeLessThanOrEqual(8e6)
  expect(errors).toEqual([])
})

test('frame times count the drawing', async ({ page }) => {
  // A slow canvas must show in p95, not hide.
  await page.addInitScript(() => {
    const put = CanvasRenderingContext2D.prototype.putImageData
    CanvasRenderingContext2D.prototype.putImageData = function (...a) {
      const until = performance.now() + 25
      while (performance.now() < until) {}
      return put.apply(this, a)
    }
  })
  await open(page, '#/c/295-starfield')
  await page.waitForFunction(() => window.__8bit.stats().frames >= 20)
  expect((await page.evaluate(() => window.__8bit.stats())).p95).toBeGreaterThanOrEqual(25)
})
