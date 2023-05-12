export type RunConfig = {
    jwt: string
}

export type CLI = {
    run(cfg: RunConfig, args: string[]): Promise<string>
}
