/******************************************************************************
 * Scaleway CLI release script.
 *
 * This script will trigger a release process for the scaleway Go SDK.
 *
 * The script will proceed as follow:
 *   - Create a new remote `scaleway-release` that point to main repo
 *   - Prompt the new version number
 *   - Create release commit
 *     - Generate a changelog
 *     - Update version in cmd/scw/main.go
 *     - Update version in the README
 *     - Ask to review changes
 *     - Create a release PR on github
 *   - Create a release
 *     - Build binary that should be upload on the github release with there SHA256SUM
 *     - Create a github release
 *     - Attach compiled binary to the github release
 *   - Update S3 version file
 *   - Create a post release commit
 *     - Update cmd/scw/main.go to add +dev at the end
 *     - Ask to merge this changes to master via PR
 *   - Delete temporary remote `scaleway-release` that was created earlier
 ******************************************************************************/

// Import standard node library
const { spawnSync } = require("child_process"),
    readline = require("readline"),
    fs = require("fs"),
    util = require("util")
;

// Import third party library
const gitRawCommits = require("git-raw-commits"),
    semver = require("semver"),
    getStream = require("get-stream"),
    externalEditor = require("external-editor"),
    _ = require("colors"),
    { Octokit } = require("@octokit/rest"),
    AWS = require('aws-sdk')
;

/*
 * Required parameters
 */
// A valid github personal token that should have write access to github repo.
const GITHUB_TOKEN=process.env.GITHUB_TOKEN;
// A scaleway access key that should have write access to the devtool bucket.
const SCW_ACCESS_KEY=process.env.SCW_ACCESS_KEY;
// A scaleway secet key that should have write access to the devtool bucket.
const SCW_SECRET_KEY=process.env.SCW_SECRET_KEY;


/*
 * Global configuration
 */
// The root git directory. All future path should be relative to this one.
const ROOT_DIR = "../..";
// The script used to build release binary
const BUILD_SCRIPT = "./scripts/build.sh";
// The directory that contains release binary created by BUILD_SCRIPT
const BIN_DIR_PATH = "./bin/";
// Path to README file
const README_PATH = "./README.md";
// Path to CHANGELOG file
const CHANGELOG_PATH = "./CHANGELOG.md";
// Go file that contain the version string to replace during release process.
const GO_VERSION_PATH = "./cmd/scw/main.go";
// Name of the temporary branch that will be used during the release process.
const TMP_BRANCH = "new-release";
// Name of the temporary remote that will be used during the release process.
const TMP_REMOTE = "scaleway-release";
// Name of the github repo namespace (user or orga).
const GITHUB_OWNER = "jerome-quere";
// Name of the github repo.
const GITHUB_REPO = "scaleway-cli";
// Name of the devtool bucket.
const S3_DEVTOOL_BUCKET="scw-devtools";
// Region of the devtool bucket .
const S3_DEVTOOL_BUCKET_REGION="nl-ams";
// S3 object name of the version file that should be updated during release.
const S3_VERSION_OBJECT_NAME="scw-cli-v2-version";

/*
 * Usefull constant
 */
const GITHUB_CLONE_URL = `git@github.com:${GITHUB_OWNER}/${GITHUB_REPO}.git"`;
const GITHUB_REPO_URL = `https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}`;
const COMMIT_REGEX = new RegExp(`${_typeReg.source}${_scopeReg.source}: *${_messageReg.source} *${_mrReg.source}`);
const _typeReg = /(?<type>[a-zA-Z]+)/;
const _scopeReg = /(\((?<scope>.*)\))?/;
const _messageReg = /(?<message>[^(]*)/;
const _mrReg = /(\(#(?<mr>[0-9]+)\))?/;


async function main() {

    /*
     * Initialization
     */

    // Chdir to root dir
    process.chdir(ROOT_DIR);

    // Initialize github client
    if (!GITHUB_TOKEN) {
        throw new Error(`You must provide a valid GITHUB_TOKEN`)
    }
    octokit = new Octokit({ auth: GITHUB_TOKEN });

    // Initialize s3 client
    if (!SCW_ACCESS_KEY || !SCW_SECRET_KEY) {
        throw new Error(`You must provide a valid access and secret key`)
    }
    const s3 = new AWS.S3({
        credentials: new AWS.Credentials(SCW_ACCESS_KEY, SCW_SECRET_KEY),
        endpoint: `s3.${S3_DEVTOOL_BUCKET_REGION}.scw.cloud`,
        region: `${S3_DEVTOOL_BUCKET_REGION}`,
    });

    /*
     * Initialize TMP_REMOTE
     */
    console.log("Adding temporary remote on local repo".blue);
    git( "remote", "add", TMP_REMOTE, GITHUB_CLONE_URL);
    console.log(`   Successfully created ${TMP_REMOTE} remote`.green);

    console.log("Make sure we are working on an up to date v2 branch".blue);
    git( "fetch", TMP_REMOTE);
    git( "checkout", `${TMP_REMOTE}/v2`);
    console.log(`   Successfully created ${TMP_REMOTE} remote`.green);

    /*
     * Trying to find the lastest tag to generate changelog
     */
    console.log("Trying to find last release tag".blue);
    const lastSemverTag = git("tag")
        .trim()
        .split("\n")
        .filter(semver.valid)
        .sort((a, b) => semver.rcompare(semver.clean(a), semver.clean(b)))[0];
    const lastVersion =  semver.clean(lastSemverTag);
    const lastVersionWithDash = lastVersion.replace(/\./g, "-");
    console.log(`    Last found release tag was ${lastSemverTag}`.green);

    console.log("Listing commit since last release".blue);
    const commits = (await getStream.array(gitRawCommits({ from: lastSemverTag, format: "%s" }))).map(c => c.toString().trim());
    commits.forEach(c => console.log(`    ${c}`.grey));
    console.log(`    We found ${commits.length} commits since last release`.green);

    const newVersion = semver.clean(await prompt("Enter new version: ".magenta));
    if (!newVersion) {
        throw new Error(`invalid version`);
    }
    const newVersionWithDash = newVersion.replace(/\./g, "-");

    //
    // Creating release commit
    //

    console.log(`Updating ${CHANGELOG_PATH} and ${GO_VERSION_PATH}`.blue);
    const changelog = buildChangelog(newVersion, commits);
    changelog.body = externalEditor.edit(changelog.body);

    replaceInFile(README_PATH, lastVersion, newVersion);
    replaceInFile(README_PATH, lastVersionWithDash, newVersionWithDash);
    replaceInFile(CHANGELOG_PATH, "# Changelog", `# Changelog\n\n${changelog.header}\n\n${changelog.body}\n`);
    replaceInFile(GO_VERSION_PATH, /Version = "[^"]*"/, `Version = "v${newVersion}"`);
    console.log(`    Update success`.green);

    await prompt(`Please review ${README_PATH}, ${CHANGELOG_PATH} and ${GO_VERSION_PATH}. When everything is fine hit enter to continue ...`.magenta);

    console.log(`Creating release commit`.blue);
    git("checkout", "-b", TMP_BRANCH);
    git("add", README_PATH, CHANGELOG_PATH, GO_VERSION_PATH);
    git("commit", "-m", `chore: release ${newVersion}`);
    git("push", "-f", "--set-upstream", TMP_REMOTE, TMP_BRANCH);

    const prResp = await octokit.pulls.create({
        owner: GITHUB_OWNER,
        repo: GITHUB_REPO,
        base: "v2",
        head: TMP_BRANCH,
        title: `chore: release ${newVersion}`
    });
    console.log(`    Successfully create pull request: ${prResp.data.html_url}`.green);
    await prompt(`Hit enter when its merged .....`.magenta);

    //
    // Creating release
    //

    console.log("Compiling release binary".blue);
    exec(BUILD_SCRIPT);

    console.log("Create Github release".blue);
    console.log("    create github release".gray);
    let releaseResp = await octokit.repos.createRelease({
        owner: GITHUB_OWNER,
        repo: GITHUB_REPO,
        tag_name: newVersion,
        target_commitish: "v2",
        name: newVersion,
        body: changelog.body,
        prerelease: true,
    });

    console.log("    attach assets to the release".gray);
    const releaseAssets = [
        `scw-${newVersionWithDash}-darwin-x86_64`,
        `scw-${newVersionWithDash}-linux-x86_64`,
        `scw-${newVersionWithDash}-windows-x86_64.exe`,
        `SHA256SUMS`,
    ];
    await Promise.all(releaseAssets.map((assetName) => {
        return octokit.repos.uploadReleaseAsset({
            owner: GITHUB_OWNER,
            repo: GITHUB_REPO,
            release_id: releaseResp.data.id,
            name: assetName,
            data: fs.readFileSync(`${BIN_DIR_PATH}/${assetName}`),
        })
    }));

    console.log(`    Successfully create release: ${releaseResp.data.html_url}`.green);
    await prompt(`Hit enter when if everything is fine .....`.magenta);

    //
    // Update version file on s3
    //
    console.log("Compiling release binary".blue);
    await util.promisify(s3.putObject.bind(s3))({
        Body: newVersion,
        Bucket: S3_DEVTOOL_BUCKET,
        Key: S3_VERSION_OBJECT_NAME
    });
    console.log(`    Successfully update s3 version file: https://${S3_DEVTOOL_BUCKET}.s3.${S3_DEVTOOL_BUCKET_REGION}.scw.cloud/${S3_VERSION_OBJECT_NAME}`.green);
    await prompt(`Hit enter to continue .....`.magenta);

    //
    // Creating post release commit
    //
    console.log("Make sure we pull the latest commit from master".blue);
    git("fetch", TMP_REMOTE);
    git("checkout", `${TMP_REMOTE}/master`);
    console.log("    Successfully checkout upstream/master".green);

    console.log(`Creating post release commit`.blue);
    git("branch", "-D", TMP_BRANCH);
    git("checkout", "-b", TMP_BRANCH);
    replaceInFile(GO_VERSION_PATH, /Version = "[^"]"/, `Version = "v${newVersion}+dev"`);
    git("add", GO_VERSION_PATH);
    git("commit", "-m", `chore: cleanup after v${newVersion} release`);
    git("push", "-f", "--set-upstream", TMP_REMOTE, TMP_BRANCH);
    git("checkout", "master");
    git("branch", "-D", TMP_BRANCH);
    const postPrResp = await octokit.pulls.create({
        owner: GITHUB_OWNER,
        repo: GITHUB_REPO,
        base: "v2",
        head: TMP_BRANCH,
    });
    console.log(`    Successfully create pull request: ${postPrResp.data.html_url}`.green);
    await prompt(`Hit enter when its merged .....`.magenta);

    console.log("Make sure we pull the latest commit from v2".blue);
    git("pull", TMP_REMOTE, "v2");
    console.log("    Successfully pull master".green);

    console.log("Remove temporary remote".blue);
    git("remote", "remove", TMP_REMOTE);
    console.log("    Successfully remove temporary remote".green);

    console.log(`🚀 Release Success `.green);
}

function git(...args) {
    return exec("git", ...args)
}

function exec(cmd, ...args) {
    console.log(`    ${cmd} ${args.join(" ")}`.grey);
    const { stdout, status, stderr } = spawnSync(cmd, args, { encoding: "utf8" });
    if (status !== 0) {
        throw new Error(`return status ${status}\n${stderr}\n`);
    }
    return stdout;
}


function replaceInFile(path, oldStr, newStr) {
    console.log(`    Editing ${path}`.grey);
    let content = fs.readFileSync(path, { encoding: "utf8" });
    if (oldStr instanceof RegExp) {
        content = content.replace(oldStr, newStr);
    } else {
        content = content.split(oldStr).join(newStr);
    }

    fs.writeFileSync(path, content);
}

function prompt(prompt) {
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    return new Promise(resolve => {
        rl.question(prompt, answer => {
            resolve(answer);
            rl.close();
        });
    });
}

function buildChangelog(newVersion, commits) {
    const changelogLines = { feat: [], fix: [], others: [] };
    commits.forEach(commit => {
        const result = COMMIT_REGEX.exec(commit);

        // If commit do not match a valid commit regex we add it in others section without formatting
        if (!result) {
            console.warn(`WARNING: Malformed commit ${commit}`.yellow);
            changelogLines.others.push(commit);
            return;
        }
        const stdCommit = result.groups;

        // If commit type is not one of [feat, fix] we add it in the other group. This will probably need further human edition.
        if (!(stdCommit.type in changelogLines)) {
            stdCommit.scope = [ stdCommit.type, stdCommit.scope ].filter(str => str).join(" - ");
            stdCommit.type = "others";
        }

        const line = [
            `*`,
            stdCommit.scope ? `**${stdCommit.scope}**:` : "",
            stdCommit.message,
            stdCommit.mr ? `([#${stdCommit.mr}](${GITHUB_REPO_URL}/pull/${stdCommit.mr}))` : ""
        ]
            .map(s => s.trim())
            .filter(v => v)
            .join(" ");
        changelogLines[stdCommit.type].push(line);
    });

    const changelogSections = [];
    if (changelogLines.feat) {
        changelogSections.push("### Features\n\n" + changelogLines.feat.sort().join("\n"));
    }
    if (changelogLines.fix) {
        changelogSections.push("### Fixes\n\n" + changelogLines.fix.sort().join("\n"));
    }
    if (changelogLines.others) {
        changelogSections.push("### Others\n\n" + changelogLines.others.sort().join("\n"));
    }
    return {
        header: `## v${newVersion} (${new Date().toISOString().substring(0, 10)})`,
        body: changelogSections.join("\n\n"),
    }
}

main().catch(console.error);
