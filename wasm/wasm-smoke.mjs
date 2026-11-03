// Smoke test for the wasm wrapper: boots the published package exactly like a
// JS consumer would, then runs `scw version` and checks the output.
import { readFile, rm } from 'node:fs/promises'

import './wasm_exec_node.cjs'
import './wasm_exec.cjs'
import { wasmURL } from './index.js'

const CALLBACK = 'cliLoaded'
const TARGET = 'scw'

const go = new globalThis.Go()
go.argv = [CALLBACK, TARGET]

const loaded = new Promise((resolve, reject) => {
  const timeout = setTimeout(
    () => reject(new Error('timed out waiting for the wasm CLI to boot')),
    30_000,
  )
  globalThis[CALLBACK] = () => {
    clearTimeout(timeout)
    resolve()
  }
})

const instance = await WebAssembly.instantiate(
  await readFile(wasmURL),
  go.importObject,
)
go.run(instance.instance)
await loaded

const resp = await globalThis[TARGET].run(
  { jwt: '', defaultProjectID: '', defaultOrganizationID: '', apiUrl: '' },
  ['version'],
)
// WASM has no HOME/SCW_CACHE_DIR env, so scw's after-command checks write
// their 24h throttle marker to the current directory. Do not leave it in the repo.
await rm('last-cli-checks', { force: true })
if (resp.exitCode !== 0) {
  console.error(`smoke test failed: exit code ${resp.exitCode}\n${resp.stderr}`)
  process.exit(1)
}
if (!/Version +2\.\d+\.\d+/.test(resp.stdout)) {
  console.error(`smoke test failed: unexpected version output:\n${resp.stdout}`)
  process.exit(1)
}
console.log(resp.stdout.trim())
