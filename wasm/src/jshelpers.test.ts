import { describe, it, expect, afterAll, beforeAll } from 'vitest'

import '../wasm_exec_node.cjs'
import '../wasm_exec.cjs'
import { RunConfig } from '../cli'
import * as fs from 'fs'
import { loadWasmBinary } from './utils'

type CLITester = {
  stop: () => Promise<void>
  FromSlice: () => string[]
}

describe('With test environment', async () => {
  let cli: CLITester

  beforeAll(async () => {
    // @ts-ignore
    cli = (await loadWasmBinary('./cliTester.wasm')) as CLITester
  })

  it('can return array', async () => {
    const array = cli.FromSlice()

    expect(array).toHaveLength(3)
    expect(array).toContain('1')
    expect(array).toContain('2')
    expect(array).toContain('3')
  })

  afterAll(async () => {
    try {
      await cli.stop()
      go._resume()
    } catch (e) {
      console.log(e)
    }
  })
})
