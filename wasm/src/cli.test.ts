// Tests have to run with go installed
// It will load go misc files in
// /usr/local/go/misc/wasm

import {describe, it, expect} from 'vitest'

import '../wasm_exec_node.cjs'
import '../wasm_exec.cjs'
import {CLI, RunConfig} from '../cli'
import * as fs from 'fs'

const CLI_PACKAGE = 'scw'
const CLI_CALLBACK = 'cliLoaded'

describe('With wasm CLI', async () => {
    // @ts-ignore
    const go = new globalThis.Go()

    const waitForCLI = new Promise((resolve) => {
        // @ts-ignore
        globalThis[CLI_CALLBACK] = () => {
            resolve({})
        }
    })
    go.argv = [CLI_CALLBACK, CLI_PACKAGE]

    WebAssembly.instantiate(fs.readFileSync('./cli.wasm'), go.importObject).then((result) => {
        return go.run(result.instance)
    }).catch((err) => {
        console.error(err)
        console.error("webassembly error")
        process.exit(1)
    })
    await waitForCLI
    // @ts-ignore
    const cli = globalThis[CLI_PACKAGE] as CLI

    const run = async (expected: string | RegExp, command: string[], runCfg: RunConfig | null = null) => {
        if (runCfg === null) {
            runCfg = {
                jwt: "",
            }
        }

        const resp = await cli.run(runCfg, command)
        expect(resp.exitCode).toBe(0)
        expect(resp.stdout).toMatch(expected)
    }

    const runWithError = async (expected: string | RegExp, command: string[], runCfg: RunConfig | null = null) => {
        if (runCfg === null) {
            runCfg = {
                jwt: "",
            }
        }
        const resp = await cli.run(runCfg, command)
        expect(resp.exitCode).toBeGreaterThan(0)
        expect(resp.stderr).toMatch(expected)
    }

    it('can run cli commands', async () => run(/profile.*default/, ['info']))

    it('can run help', async () => runWithError(/USAGE:\n.*scw <command>.*/, []))

    it('can use jwt', async () => runWithError(/.*denied authentication.*invalid JWT.*/, ['instance', 'server', 'list']))
})
