# Release Notes Guidelines

Guidelines for writing and formatting release notes for this repository. Release
notes are published as GitHub Releases, not as entries in a changelog file.

## Process

When a **specific version's** notes are asked for, draft that version only.

1. Scope the release to its two nearest tags: `git log --oneline FROM..TO`,
   where `FROM` is the previous released tag (eg: v1.2.3) and `TO` the
   version being drafted (or HEAD if the tag doesn't exists yet). The note must
   cover every user-facing change in this range and nothing outside it.
2. **Bootstrap** the skeleton from the commit list — it gives commit subjects
   grouped by type and namespace, not the final text:

   - Commit list, one line per commit with its short and full SHA:
     ```sh
     git log --format="%s [%h](%H)" FROM..TO | sort
     ```
   - The `feat`/`fix` subjects that drive the sections, grouped by their `scope(...)` — the
     namespace list comes from the actual range, never a maintained list:
     ```sh
     git log --format="%s" FROM..TO | grep -oE '^(feat|fix)\([^)]*' | sort | uniq -c
     ```
   - Author lines for the Contributors section, deduped and sorted:
     ```sh
     git log --format="%aN <%aE>" FROM..TO | sort -u
     ```
3. **Redact** the result by hand (see [Redaction style](#redaction-style)): drop chore/ci/build/`test`
   and dependency-bot items, rewrite the rest as prose with "why" and usage examples.
4. Keep the **end-result** notes in `docs/release_notes/v<version>.md` — one file per
   released version, named after the tag, in the `docs/release_notes/` directory.
5. **Prerelease tags (`rc`, `beta`, debug/`-dbg`) are never given notes** — only final releases
   get a `v<version>.md` file.

## What to include and what to skip

- Summarize dependency bumps, CI and internal migration work in **one short
  "Internal and tooling" section** — never list them individually.
- Do not describe internal details a user cannot observe (`go-vcr` migration,
  lint config, workflow reshuffles, ...).
- Mention PR IDs in the body only if the reader can act on them (bug fixes,
  breaking changes, notable features).

### Include

- `feat(scope)` — new or improved behavior                        
- `fix(scope)` — bugs users can hit, with the symptom             
- migrations, renames, deprecations (as breaking changes) 
- New namespaces/products

### Exclude

- `chore`, `ci`, `build`, `docs`, `test(s)` (unless user-facing) 
- `fix(deps)` / `chore(deps)` / dependabot & renovate bumps      
- internal refactors with no behavioral change                   
- codegen/test-tooling churn                                     

## Redaction style

"Redaction" is the manual pass that turns a generated or commit-subject draft into prose people
read in a changelog.

### Structure

Sections appear in a **fixed order** so output is stable across releases. Use a
**two-level table of contents**: every `##`
section is listed, and each `###` subsection appears indented two spaces under
its section, all matching their headings' GitHub slugs. Then:

1. **Highlights** — the 3-6 biggest user-facing wins, one bullet each.
2. **Breaking changes and migrations** — always near the top, even if empty.
   Call out renames, deprecated flags/endpoints, and API migrations users must
   act on.
3. **Core** — cross-cutting behavior not tied to a specific product namespace.
   Only include it when there is material to report.
4. **Product namespaces** — one section per namespace (the `scope(...)` of the
   changes), **sorted alphabetically, case-insensitive, by heading** — e.g.
   `Audit trail`, `Autoscaling`, `k8s`, `Kafka`, `MongoDB`, `RDB`,
   `Secret Manager`. `Core` is not part of this run: it lives at step 3.
   The section list is the set of scopes emitted by the step-2 one-liner, and each
   scope key is resolved to the product's official display name from the source
   (`internal/namespaces/<name>/` directory, `docs/commands/<name>.md`), in the same
   style as past notes — no maintained namespace list. Every change in a namespace
   lives under its own `###` subsection — never a bare bullet directly under a `##`
   namespace heading.
5. **Bug fixes** — general fixes; keep them here even when a fix targets a
   specific product.
6. **Contributors**
7. **Internal and tooling** (brief)

Omit steps with nothing to report (e.g. no `Core` section when there are no
core features, smaller releases may have no `Highlights`); keep the listed
order for the steps that remain.

### Per-section content rules

#### Highlights

- Restate each win in the release's own words, not the commit subject ("k8s connect opens k9s",
  not `feat(k8s): add support for connect using k9s`).
- One bullet per win, **bold lead-in**, PR and commit reference.
- Keep to the 3-6 biggest user-facing wins.

#### Features — for each one

- Say **why** it exists: the problem it solves for the user (one sentence).
- Give a **barebones usage example** in a fenced code block. One or two commands
  is enough; resist pasting whole help output.
- Reference the PR and the commit that introduced the change, with links.

#### Bug fixes — symptom first

- State **what was wrong** in terms a user recognizes (error message, crash,
  wrong output), not just the internal cause.
  - e.g. "completing `engine=` suggested disabled V1 engine names" beats
    "improved engine autocompletion".
- Keep the internal explanation if useful, after the symptom.
- Always link the PR and the commit so users can find details.

#### Breaking changes

- Lead with the action the user must take ("update scripts that pass X").
- Link the migration PR(s) and commit(s).
- Group related migrations under the product scope in one bullet.

#### Contributors

- List authors of non-bot commits in the release range, sorted alphabetically by login.
- **Link every contributor to their GitHub profile**: `[@login](https://github.com/login)`.
  Resolve the login from the commit author (via the commits API, the PR author,
  or a verified `api.github.com/users/<login>`) — never guess from emails. When the
  git author name isn't obviously the login, fetch it from the commit's GitHub author
  with the gh CLI: `gh api repos/{owner}/{repo}/commits/<sha> --jq '.author.login'`.
- Tag first-time committers inline (e.g. a `(new)` marker beside the handle) with their first
  contribution and its PR and commit links.
- Exclude bots and their automated commits.
- **No emoji** here or anywhere else in the note (see [Formatting rules](#formatting-rules)).

### Formatting rules

- **PR links:** use `[#NNNN](https://github.com/{owner}/{repo}/pull/NNNN)`.
- **Commit links:** every referred change (highlight, feature, bug fix, breaking
  change, first contribution) also links the commit that introduced it:
  `[commit 1a2b3c4](https://github.com/{owner}/{repo}/commit/<full-sha>)`.
  GitHub resolves shortened SHAs, but use the full SHA in the URL.
- **Contributor links:** every contributor is a GitHub-profile link,
  `[@login](https://github.com/login)`, listed in the Contributors section (see
  [Contributors](#contributors)).
- **Range link:** the header links the tag range it covers:
  `[**v1.1.0 → v1.2.0**](https://github.com/{owner}/{repo}/compare/v1.1.0...v1.2.0)`.
- **No emoji:** avoid emojis entirely — in headings, bullets, examples and contributor markers.
  Hand-added decorations (e.g. `⚠️` in "Breaking changes and
  migrations", `🎉` for first-time contributors) make output inconsistent; use plain text
  instead (e.g. a `(new)` marker for first-time contributors).
- **Anchor-safe headings:** keep headings free of characters that make GitHub
  slugs ambiguous — no `&` (write "and"), no emoji, no `--` flags, no slashes in
  heading text. Keep backticked flags *inside the body*, not in the heading.
- **Two-level TOC:** list every `##` section in file order, with each `###`
  subsection indented two spaces under its section. Every entry matches its
  heading's GitHub slug. Check by stripping punctuation/emoji, lowercasing,
  and replacing whitespace runs with `-`.
- **Namespace subsections:** in `##` product namespace sections, every change is
  a `###` subsection (prose + example + its PR/commit links). No bare bullets
  directly under a namespace heading — e.g. `## Kafka` must contain
  `### Volume resize`, not a lone bullet. Bullets are reserved for the meta
  sections (Highlights, Breaking changes, Bug fixes, Contributors, Internal).
- Use fenced code blocks (```sh) for usage examples.
- Prefer dashes/backticks from the real command (`scw instance server list`)
  over paraphrasing.

## Verification checklist

Before publishing, verify:

- [ ] Every user-facing commit in the version's tag range is represented, and no
      item from a neighbouring version's note leaked in.
- [ ] The note was redacted in the style above: prose with "why" and usage
      examples, not raw commit subjects.
- [ ] Every PR link resolves (HTTP 200) and every commit hash exists in the
      release's tag range.
- [ ] TOC is two-level and complete: every `##` section and every `###`
      subsection is listed, indented correctly, and its anchor matches.
- [ ] Every namespace change sits under a `###` subsection (no bare bullets
      under a `##` namespace heading).
- [ ] Every referred change (highlights, features, fixes, breaking changes)
      links both its PR and its commit.
- [ ] All TOC anchors match their headings (see rule above).
- [ ] No dependency-bump / bot commits listed individually.
- [ ] All bug fixes state the observable symptom.
- [ ] All highlighted features have a "why" and a usage example.
- [ ] Contributor handles verified (e.g. against `api.github.com/users/<login>`
      or the PR author), not guessed from emails.
- [ ] Every contributor is a GitHub-profile link
      (`[@login](https://github.com/login)`), sorted by login.
- [ ] No emoji anywhere in the note (headings, bullets, contributor markers).
- [ ] Breaking changes are listed near the top and tell the user what to change.

## Maintaining this file

- This file lives at the repo root; the end-result notes live in `docs/release_notes/`.
- Keep this guide in sync with how release notes are actually produced: update
  it when the structure, tooling or conventions change.
