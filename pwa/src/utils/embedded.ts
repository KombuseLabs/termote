import { readSettings, writeSettingsPatch } from '../hooks/use-settings'
import { WELCOME_VERSION } from './welcome'

// Shown in another program's window for one person (health's embedded,
// TERMOTE_EMBEDDED): panes fit that window (a herdr-core terminal takes its
// size, as a tmux one does by itself) and the welcome is not shown. Once
// per page, so a choice made in Settings afterwards stays.
let applied = false

export function applyEmbedded() {
  if (applied) return
  applied = true
  const s = readSettings()
  if (s.driveTerminalSize && s.welcomeSeenVersion >= WELCOME_VERSION) return
  writeSettingsPatch({
    driveTerminalSize: true,
    welcomeSeenVersion: WELCOME_VERSION,
  })
}

// Test seam: a fresh page.
export function resetEmbedded() {
  applied = false
}
