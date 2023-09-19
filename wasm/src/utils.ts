import '../wasm_exec_node.cjs'
import '../wasm_exec.cjs'
import * as fs from 'fs'
import { Go } from '../wasm_exec'
import { expect } from 'vitest'

const CLI_PACKAGE = 'scw'
const CLI_CALLBACK = 'cliLoaded'

export const loadWasmBinary = async (binaryName: string): Promise<unknown> => {
  // @ts-ignore
  const go = new globalThis.Go() as Go

  const waitForCLI = new Promise(resolve => {
    // @ts-ignore
    globalThis[CLI_CALLBACK] = () => {
      resolve({})
    }
  })
  go.argv = [CLI_CALLBACK, CLI_PACKAGE]

  const buffer: BufferSource = await fs.promises.readFile(binaryName)

  WebAssembly.instantiate(buffer, go.importObject)
    .then(result => {
      return go.run(result.instance)
    })
    .catch(err => {
      expect(err, err).toBeNull()
      process.exit(1)
    })
  await waitForCLI

  return globalThis[CLI_PACKAGE]
}
