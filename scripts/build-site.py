#!/usr/bin/env python3
"""Write the shared parts of the site into the pages that are not the front one.

The header with its navigation, the footer, the icon sprite and the navigation
script are the same on every page, and they are about ten kilobytes of markup.
Keeping a copy per page means three places to change and two of them to forget,
which is how the navigation went missing from the legal pages in the first
place.

The front page is the source. Each other page carries a pair of markers per
block and nothing between them that a person maintains, so a page holds its own
head, its own main region, and no shell at all.

Run it after changing the header, the footer, the sprite or the script:

    python3 scripts/build-site.py

Passing --check writes nothing and fails where a page no longer matches what the
front page would produce, which is what stops a change to the header reaching
main whilst the other pages still carry the previous one. scripts/gates.sh runs
it that way.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

DOCS = Path(__file__).resolve().parent.parent / "docs"
FRONT = DOCS / "index.html"

# Every page below DOCS that takes the shell, with the footer link that should
# mark itself as the current page.
PAGES = {
    DOCS / "imprint" / "index.html": "../imprint/",
    DOCS / "privacy" / "index.html": "../privacy/",
}

BLOCKS = {
    "icons": r"<!-- icons:start -->.*?<!-- icons:end -->",
    "header": r"<header class=\"band\">.*?</header>",
    "footer": r"<footer class=\"footer\">.*?</footer>",
    "script": r"<script>\n.*?</script>",
}


def read_blocks(front: str) -> dict[str, str]:
    """Cut the shared blocks out of the front page."""
    blocks = {}
    for name, pattern in BLOCKS.items():
        match = re.search(pattern, front, re.DOTALL)
        if match is None:
            raise SystemExit(f"docs/index.html has no {name} block")
        blocks[name] = match.group(0)
    return blocks


def descend(markup: str, current: str) -> str:
    """Rewrite the front page's links for a page one directory further down.

    Only anchors are touched. A `use` element's href points into the sprite that
    travels with the markup, so it stays as it is, and an absolute URL is
    already right from anywhere.
    """

    def rewrite(match: re.Match[str]) -> str:
        attributes, target = match.group(1), match.group(2)
        if target == "#top":
            target = "../"
        elif target.startswith("#"):
            target = "../" + target
        elif not target.startswith(("http://", "https://", "mailto:", "../")):
            target = "../" + target
        return f"<a{attributes}href=\"{target}\""

    markup = re.sub(r"<a([^>]*?)href=\"([^\"]+)\"", rewrite, markup)

    # No section of the front page is the current page from here.
    markup = markup.replace(' aria-current="page"', "")

    # The footer link to the page being written is the one that is current.
    markup = markup.replace(
        f'<a href="{current}">', f'<a href="{current}" aria-current="page">'
    )
    return markup


def apply(page: Path, blocks: dict[str, str], current: str) -> str:
    """Return the page with every marked block replaced."""
    markup = page.read_text()
    for name, block in blocks.items():
        pattern = re.compile(
            rf"(<!-- shell:{name}:start -->\n).*?(<!-- shell:{name}:end -->)",
            re.DOTALL,
        )
        if not pattern.search(markup):
            raise SystemExit(f"{page} has no shell:{name} markers")
        replacement = descend(block, current)
        markup = pattern.sub(
            lambda match: match.group(1) + replacement + "\n" + match.group(2),
            markup,
            count=1,
        )
    return markup


def main() -> int:
    check_only = "--check" in sys.argv[1:]
    blocks = read_blocks(FRONT.read_text())
    stale = []
    for page, current in PAGES.items():
        written = apply(page, blocks, current)
        if written == page.read_text():
            continue
        stale.append(page.relative_to(DOCS.parent))
        if not check_only:
            page.write_text(written)

    if check_only:
        if stale:
            for path in stale:
                print(f"{path} does not carry the current shell", file=sys.stderr)
            print(
                "Rebuild them with: python3 scripts/build-site.py", file=sys.stderr
            )
            return 1
        print("every page carries the current shell")
        return 0

    for path in stale:
        print(f"updated {path}")
    if not stale:
        print("every page already carries the current shell")
    return 0


if __name__ == "__main__":
    sys.exit(main())
