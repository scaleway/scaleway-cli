export class Go {
    argv: string[]
    env: {[key: string]: string}
    async run(instance: WebAssembly.Instance)
}