import { describe, it, expect, afterAll, beforeAll } from 'vitest'

import { loadWasmBinary } from './utils'

type CLITester = {
  stop: () => Promise<void>
  FromSlice: () => string[]
  MarshalBuildInfo: () => string
}

describe('With test environment', async () => {
  let cli: CLITester

  beforeAll(async () => {
    cli = await loadWasmBinary('./cliTester.wasm') as CLITester
  })

  it('can return array', async () => {
    const array = cli.FromSlice()

    expect(array).toHaveLength(3)
    expect(array).toContain('1')
    expect(array).toContain('2')
    expect(array).toContain('3')
  })

  it('can marshal build info', async () => {
    const buildInfo = cli.MarshalBuildInfo()
    expect(buildInfo).toContain('Version')
    expect(buildInfo).toContain('2.0.0')
  })

  // Note: skipping cli.stop() cleanup as it triggers Go runtime errors
  // The WASM runtime is cleaned up when the Node.js process exits
})
