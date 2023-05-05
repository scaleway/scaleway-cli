export type CLI = {
    run(...args: string[]): Promise<string>
}
