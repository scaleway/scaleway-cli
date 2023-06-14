export type RunConfig = {
    jwt: string
    defaultProjectID: string
}

export type RunResponse = {
    stdout: string
    stderr: string
    exitCode: string
}

export type AutocompleteConfig = {
    jwt: string
    leftWords: string[]
    selectedWord: string
    rightWords: string[]
}

export type CLI = {
    run(cfg: RunConfig, args: string[]): Promise<RunResponse>
    complete(cfg: AutocompleteConfig): Promise<string[]>
    stop(): Promise<void>
}
