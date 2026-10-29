export class Go {
  argv: string[]
  env: { [key: string]: string }
  importObject: {
    go: { [key: string]: (sp: number) => void }
  }
  run(instance: WebAssembly.Instance): Promise<void>
  _resume(): void
}
