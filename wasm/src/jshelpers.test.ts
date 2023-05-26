
import {describe, it, expect, afterAll} from 'vitest'

import '../wasm_exec_node.cjs'
import '../wasm_exec.cjs'
import {RunConfig} from '../cli'
import * as fs from 'fs'

const CLI_PACKAGE = 'scw'
const CLI_CALLBACK = 'cliLoaded'

type CLITester = {
    stop: () => Promise<void>
    FromSlice: () => string[]
}

describe('With test environment', async () => {
    // @ts-ignore
    const go = new globalThis.Go()

    const waitForCLI = new Promise((resolve) => {
        // @ts-ignore
        globalThis[CLI_CALLBACK] = () => {
            resolve({})
        }
    })
    go.argv = [CLI_CALLBACK, CLI_PACKAGE]

    WebAssembly.instantiate(fs.readFileSync('./cliTester.wasm'), go.importObject).then((result) => {
        return go.run(result.instance)
    }).catch((err) => {
        console.error(err)
        console.error("webassembly error")
        process.exit(1)
    })
    await waitForCLI
    // @ts-ignore
    const cli = globalThis[CLI_PACKAGE] as CLITester

    it('can return array', async () => {
        const array = cli.FromSlice()

        expect(array).toHaveLength(3)
        expect(array).toContain("1")
        expect(array).toContain("2")
        expect(array).toContain("3")
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