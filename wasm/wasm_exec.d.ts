export class Go {
    argv: string[]
    env: {[key: string]: string}
    importObject: {
        go: {[key: string]: (sp: number) => void},
    }
    async run(instance: WebAssembly.Instance)
    _resume()
}
