// Tests have to run with go installed
// It will load go misc files in
// /usr/local/go/misc/wasm

import {describe, it, expect} from 'vitest'

import '../wasm_exec_node.js'
import '../wasm_exec.js'
import {CLI} from '../cli'
import * as fs from 'fs'

const CLI_PACKAGE = 'scw'
const CLI_CALLBACK = 'cliLoaded'

const runWithError = async (cli: CLI, expected: string | RegExp, ...command: string[]) => {
    await expect((async () => await cli.run(...command))).rejects.toThrowError(expected)
}

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

    it('can run cli commands', async () => {
        const res = await cli.run('info')
        expect(res).toMatch(/profile.*default/)
    })

    it('can run help', async () => runWithError(cli, /USAGE:\n.*scw <command>.*/))
})
