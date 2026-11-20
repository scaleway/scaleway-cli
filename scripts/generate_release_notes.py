#!/usr/bin/env python3
"""Generate release notes for one release range, following RELEASE_NOTES_GUIDELINES.md.

Usage:
  python3 scripts/generate_release_notes.py <repo> <version> [--prev <prev-tag>] [--out <dir>]

The tag range is explicit: notes are built from `git log <prev>..<version>`.
Omit --prev for the initial release (all commits reachable from <version>).

Prerelease tags (rc, beta, -dbg) should never be passed: only final releases
get notes.
"""
import argparse
import os
import re
import subprocess

BOT_RE = re.compile(r"(scaleway[- ]bot|dependabot|renovate|github-actions|github\.com)", re.I)


def git(repo, *args):
    return subprocess.run(["git", "-C", repo, *args], capture_output=True, text=True, check=True).stdout


def tag_date(repo, tag):
    return git(repo, "log", "-1", "--format=%ci", tag).strip()[:10]


def commits(repo, a, b):
    """List of dict(hash, author, email, date, subject) for range a..b."""
    out = git(repo, "log", "--format=%H|%aI|%aN|%aE|%s", f"{a}..{b}").splitlines()
    return [dict(hash=h, date=d, name=n, email=e, subject=s)
            for h, d, n, e, s in [l.split("|", 4) for l in out]]


def commits_from(repo, ref):
    out = git(repo, "log", "--format=%H|%aI|%aN|%aE|%s", ref).splitlines()
    return [dict(hash=h, date=d, name=n, email=e, subject=s)
            for h, d, n, e, s in [l.split("|", 4) for l in out]]


def pr_link(subject):
    m = re.search(r"\(#(\d+)\)\s*[\"']*$", subject)
    if not m:
        return subject, None
    clean = subject[: m.start()].rstrip(" ,;\"'")
    return clean, m.group(1)


def scope(subject):
    m = re.match(r"(?:feat|fix)\(([^,)]+)", subject, re.I)
    return m.group(1).strip() if m else ""


SCOPE_NAMES = {
    "k8s": "k8s", "rdb": "RDB", "iam": "IAM", "lb": "LB", "vpc": "VPC",
    "vpcgw": "VPC Gateway", "dns": "DNS", "ipam": "IPAM", "s3": "S3",
    "mnq": "MNQ", "object": "Object Storage", "sbs": "SBS", "tem": "TEM",
    "wbh": "Webhosting", "sm": "Secret Manager", "kms": "KMS",
    "audit_trail": "Audit trail", "audit-trail": "Audit trail",
    "edge_services": "Edge services", "apple_silicon": "Apple Silicon",
    "key_manager": "Key Manager", "secret_manager": "Secret Manager",
    "mongodb": "MongoDB", "kafka": "Kafka", "redis": "Redis",
    "data_orchestrator": "Data Orchestrator", "multiple_namespaces": "Core",
    "core": "Core",
}


def pretty_scope(scp):
    if scp.lower() in SCOPE_NAMES:
        return SCOPE_NAMES[scp.lower()]
    words = re.split(r"[_-]", scp)
    return " ".join(w.capitalize() for w in words) if words else scp


def strip_prefix(subject):
    return re.sub(r"^(?:feat|fix)(?:\([^)]*\))?:\s*", "", subject, flags=re.I)


def is_revert(subject):
    return bool(re.search(r"\brevert\b", subject, re.I))


def is_breaking(subject):
    if re.search(r"\b(golden|test|tests|typo|fixture)\b", subject, re.I):
        return False
    return bool(re.search(r"\b(migrations?|migrating|rename|renames?|removed|removes?|deprecat|breaking|drop support|remove)\b", subject, re.I))


def is_bot(name):
    return bool(BOT_RE.search(name or ""))


def new_contributor_map(repo):
    """name -> earliest commit timestamp across all history (any email)."""
    out = [l.split("|", 1) for l in git(repo, "log", "--all", "--format=%aN|%aI").splitlines() if "|" in l]
    first = {}
    for name, ts in out:
        first[name] = min(first.get(name, ts), ts)
    return first


def anchor(heading):
    a = heading.lower().strip("# ")
    a = re.sub(r"[`*_]", "", a)
    a = re.sub(r"[^\w\s-]", "", a)
    a = re.sub(r"\s+", "-", a)
    return a


def fmt_links(entries):
    out = []
    for subject in entries:
        clean, pr = pr_link(subject)
        clean = strip_prefix(clean)
        if pr:
            out.append(f"- {clean} ([#{pr}](https://github.com/scaleway/scaleway-cli/pull/{pr}))")
        else:
            out.append(f"- {clean}")
    return out


def build(repo, version, prev):
    date = tag_date(repo, version)
    if prev:
        log = commits(repo, prev, version)
    else:
        log = commits_from(repo, version)

    feats, fixes, breaking, other = [], [], [], []
    for c in log:
        if is_revert(c["subject"]):
            fixes.append(c["subject"])  # reverts undo a change, surfaced as fixes
            continue
        if not re.match(r"feat|fix", c["subject"], re.I):
            other.append(c)
            continue
        if re.search(r"^dep\w*:|deps|dependabot|renovate", c["subject"], re.I):
            other.append(c)
            continue
        if is_breaking(c["subject"]):
            breaking.append(c["subject"])
        if c["subject"].lower().startswith("feat"):
            feats.append(c["subject"])
        else:
            fixes.append(c["subject"])

    feat_scopes = {}
    for s in feats:
        feat_scopes.setdefault(scope(s), []).append(s)

    seen = {}
    for c in log:
        if not is_bot(c["name"]):
            seen.setdefault(c["name"], c["date"])
    first_map = new_contributor_map(repo)
    window_start = tag_date(repo, prev) if prev else "0000-00-00"
    newbies = [n for n, d in seen.items() if first_map.get(n, d)[:10] > window_start]
    names = list(seen)

    lines = [f"# Release Notes — {version}", ""]
    lines.append(f"**Initial release** · {date}" if not prev else f"**{prev} → {version}** · {date} · {len(log)} commits")
    lines.append("")

    if not log:
        lines.append("No changes — this release is a re-tag (packaging/release process only).")
        return "\n".join(lines) + "\n"

    sections = []

    hi = [s for s in feats if re.search(r"(activate cli|new namespace|initial support|new product|first release of)", s, re.I)]
    if hi:
        sections.append(("## Highlights", fmt_links(hi[:5]), "highlights"))

    br_body = fmt_links(breaking)
    if br_body:
        sections.append(("## Breaking changes and migrations", br_body, "breaking-changes-and-migrations"))

    feat_lines = []
    # Stable section order: Highlights -> Breaking changes -> Core -> namespaces
    # (alphabetical) -> Bug fixes -> Contributors -> Internal and tooling.
    # Core hosts core-cross-cutting features (and scopes the commit subject
    # doesn't declare, a safe default); every other scope becomes its own
    # alphabetically-sorted section, so output order never depends on history.
    CORE_SCOPES = {"", "core", "multiple_namespaces"}
    core_feats = []
    for scp in CORE_SCOPES:
        core_feats += feat_scopes.pop(scp, [])
    if core_feats:
        sections.append(("## Core", fmt_links(core_feats), "core"))
    # Group by resolved section title so aliased scopes (e.g. audit_trail /
    # audit-trail) land in one section, then sort those titles alphabetically.
    ns_feats = {}
    for scp, items in feat_scopes.items():
        ns_feats.setdefault(pretty_scope(scp), []).extend(items)
    for title in sorted(ns_feats, key=str.lower):
        sections.append((f"## {title}", fmt_links(ns_feats[title]), anchor(title)))

    fix_body = fmt_links(fixes)
    if fix_body:
        sections.append(("## Bug fixes", fix_body, "bug-fixes"))

    if names:
        contrib_lines = (["Thanks to everyone who contributed to this release:" if len(names) > 1 else "Contributed by:"] + [""]
                         + [f"- {n}" for n in names])
        if newbies:
            contrib_lines += ["", "### New contributors", ""] + [f"- {n}" for n in newbies]
        sections.append(("## Contributors", contrib_lines, "contributors"))

    if other:
        dep_n = sum(1 for c in other if re.search(r"(deps|dependabot|renovate|bump)", c["subject"], re.I))
        if dep_n:
            tail = [f"- Dependency bumps ({dep_n}), CI/tooling and internal changes are not listed individually."]
        else:
            tail = [f"- Internal changes ({len(other)} commits) are not listed individually."]
        sections.append(("## Internal and tooling", tail, "internal-and-tooling"))

    if sections:
        lines += ["## Table of Contents", ""]
        for head, _, anc in sections:
            title = head.lstrip("# ").strip()
            lines.append(f"- [{title}](#{anchor(title)})")
        lines += ["", "---", ""]

    for head, body, _ in sections:
        lines += [head, ""] + body + [""]

    while lines and lines[-1] == "":
        lines.pop()
    return "\n".join(lines) + "\n"


def main():
    parser = argparse.ArgumentParser(description="Generate release notes for one tag range.")
    parser.add_argument("repo", help="path to the git repository")
    parser.add_argument("version", help="tag being released, e.g. v2.66.0 (final release only)")
    parser.add_argument("--prev", default="", help="previous release tag; omit for the initial release")
    parser.add_argument("--out", default="docs/release_notes", help="output directory (default: docs/release_notes)")
    args = parser.parse_args()

    if not re.fullmatch(r"v[0-9]+(\.[0-9]+)+", args.version) or re.search(r"(rc|beta|dbg)", args.version, re.I):
        parser.error("expected a final release tag like v2.66.0 (prerelease/rc/beta/dbg tags are excluded)")
    text = build(args.repo, args.version, args.prev or None)
    path = os.path.join(args.out, f"{args.version}.md")
    os.makedirs(args.out, exist_ok=True)
    with open(path, "w") as f:
        f.write(text)
    print(f"wrote {path}")


if __name__ == "__main__":
    main()
