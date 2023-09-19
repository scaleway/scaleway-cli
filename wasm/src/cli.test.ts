// Tests have to run with go installed
// It will load go misc files in
// /usr/local/go/misc/wasm

import { describe, it, expect, afterAll, beforeAll } from 'vitest'

import { CLI, RunConfig } from '../cli'
import { loadWasmBinary } from './utils'
import {wasmURL} from "../index";


const emptyConfig = (override?: {
  jwt?: string
  defaultProjectID?: string
  defaultOrganizationID?: string
  apiUrl?: string
}): RunConfig => ({
  jwt: override?.jwt || '',
  defaultProjectID: override?.defaultProjectID || '',
  defaultOrganizationID: override?.defaultOrganizationID || '',
  apiUrl: override?.apiUrl || '',
})

describe('With wasm CLI', async () => {
  let cli: CLI

  beforeAll(async () => {
    // @ts-ignore
    cli = (await loadWasmBinary(wasmURL)) as CLI
  })

  const run = async (
    expected: string | RegExp,
    command: string[],
    runCfg: RunConfig | null = null,
  ) => {
    if (runCfg === null) {
      runCfg = emptyConfig()
    }

    const resp = await cli.run(runCfg, command)
    expect(resp.exitCode).toBe(0)
    expect(resp.stdout, `out: '${resp.stdout}'`).toMatch(expected)
  }

  const runWithError = async (
    expected: string | RegExp,
    command: string[],
    runCfg: RunConfig | null = null,
  ) => {
    if (runCfg === null) {
      runCfg = emptyConfig()
    }
    const resp = await cli.run(runCfg, command)
    expect(resp.exitCode).toBeGreaterThan(0)
    expect(resp.stderr, `error: '${resp.stderr}'`).toMatch(expected)
  }

  const complete = async (
    expected: string[],
    command: string[],
    runCfg: RunConfig | null = null,
  ) => {
    if (runCfg === null) {
      runCfg = emptyConfig()
    }
    let toComplete = command.pop() || ''

    const suggestions = await cli.complete({
      jwt: runCfg.jwt,
      leftWords: command,
      rightWords: [],
      selectedWord: toComplete,
    })
    expected.forEach(suggestion => expect(suggestions).toContain(suggestion))
  }

  it('can run cli commands', async () => run(/profile.*default/, ['info']))

  it('can run help', async () => runWithError(/USAGE:\n.*scw <command>.*/, []))

  it('can use jwt', async () =>
    runWithError(/.*denied authentication.*invalid JWT.*/, [
      'instance',
      'server',
      'list',
    ]))

  it('can complete', async () =>
    complete(['server', 'image', 'volume'], ['instance', '']))

  it('can configure terminal size', async () => {
    const runCfg = emptyConfig()

    await cli.configureOutput({ width: 100, color: false })
    const resp = await cli.run(runCfg, ['marketplace', 'image', 'list'])
    expect(resp.exitCode, `error: '${resp.stderr}'`).toBe(0)

    const lines = resp.stdout.split('\n')
    expect(lines.length).toBeGreaterThan(1)
    expect(lines[2].length).toBeLessThan(100)
  })

  it('can enable colors', async () => {
    const runCfg = emptyConfig()

    await cli.configureOutput({ width: 100, color: false })
    const resp = await cli.run(runCfg, ['invalid'])
    await cli.configureOutput({ width: 100, color: true })
    const coloredResp = await cli.run(runCfg, ['invalid'])

    expect(coloredResp.stderr.length).toBeGreaterThan(resp.stderr.length)
    expect(coloredResp.stderr).not.toEqual(resp.stderr)
  })

  it('print version', async () => run(/Version +2\.\d+\.\d+.*/, ['version']))

  afterAll(async () => {
    try {
      await cli.stop()
      go._resume()
    } catch (e) {
      console.log(e)
    }
  })
})
