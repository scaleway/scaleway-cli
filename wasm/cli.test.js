// Tests have to run with go installed
// It will load go misc files in
// /usr/local/go/misc/wasm

import {describe, it, expect} from "vitest";

import "./wasm_exec_node"
import "/usr/local/go/misc/wasm/wasm_exec"

describe('With wasm CLI', async () => {
    const go = new Go();

    const cliLoaded = new Promise((resolve) => {
        globalThis.cliLoaded = () => {
            resolve()
        }
    })

    WebAssembly.instantiate(fs.readFileSync("./cli.wasm"), go.importObject).then((result) => {
        return go.run(result.instance);
    }).catch((err) => {
        console.error(err);
        process.exit(1);
    });
    await cliLoaded

    it('can run cli commands', async () => {
        const res = await cliRun("info")
        expect(res).toMatch(/profile.*default/)
    })
})
