// Tests have to run with go installed
// It will load go misc files in
// /usr/local/go/misc/wasm

import {describe, it, expect} from "vitest";

import "./wasm_exec_node"
import "/usr/local/go/misc/wasm/wasm_exec"

const CLI_PACKAGE = 'scw'
const CLI_CALLBACK = 'cliLoaded'

describe('With wasm CLI', async () => {
    const go = new Go();
    const waitForCLI = new Promise((resolve) => {
        globalThis[CLI_CALLBACK] = () => {
            resolve()
        }
    })
    go.argv = [CLI_CALLBACK, CLI_PACKAGE]

    WebAssembly.instantiate(fs.readFileSync("./cli.wasm"), go.importObject).then((result) => {
        return go.run(result.instance);
    }).catch((err) => {
        console.error(err);
        process.exit(1);
    });
    await waitForCLI
    const cli = globalThis[CLI_PACKAGE]

    it('can run cli commands', async () => {
        const res = await cli.run("info")
        expect(res).toMatch(/profile.*default/)
    })
})
