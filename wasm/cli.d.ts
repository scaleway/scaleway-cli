export type RunConfig = {
    jwt: string
}

export type RunResponse = {
    stdout: string
    stderr: string
    exitCode: string
}

export type CLI = {
    run(cfg: RunConfig, args: string[]): Promise<RunResponse>
}
