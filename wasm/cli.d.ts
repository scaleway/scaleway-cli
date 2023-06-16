export type RunConfig = {
    jwt: string
    defaultProjectID: string
    defaultOrganizationID: string
    apiUrl: string
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

export type Int = number

export type ConfigureOutputConfig = {
    width: Int
    color: boolean
}


export type CLI = {
    run(cfg: RunConfig, args: string[]): Promise<RunResponse>
    complete(cfg: AutocompleteConfig): Promise<string[]>
    configureOutput(cfg: ConfigureOutputConfig): Promise<{}>
    stop(): Promise<void>
}
