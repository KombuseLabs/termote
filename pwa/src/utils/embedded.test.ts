import { beforeEach, describe, expect, it } from 'vitest'
import { readSettings, writeSettingsPatch } from '../hooks/use-settings'
import { applyEmbedded, resetEmbedded } from './embedded'
import { WELCOME_VERSION } from './welcome'

describe('applyEmbedded', () => {
  beforeEach(() => {
    localStorage.clear()
    resetEmbedded()
  })

  it('fits panes to the window and skips the welcome, once per page', () => {
    expect(readSettings().driveTerminalSize).toBe(false)
    applyEmbedded()
    expect(readSettings().driveTerminalSize).toBe(true)
    expect(readSettings().welcomeSeenVersion).toBe(WELCOME_VERSION)
    // Turned off in Settings afterwards: it stays off on this page
    writeSettingsPatch({ driveTerminalSize: false })
    applyEmbedded()
    expect(readSettings().driveTerminalSize).toBe(false)
  })
})
