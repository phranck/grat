#!/bin/sh
set -eu

# Render the share card to docs/og-image.png.
#
# The card is HTML taking the site's own stylesheet, so the colours, the terminal
# and the wordmark are the page's rather than a second set drawn by hand. This
# script is what turns it into the PNG a social preview needs, and running it
# again after a design change is what keeps the two the same.
#
# It needs a headless Chromium. Playwright's is used where it is present, since
# the repository already has it for browser work, and Google Chrome otherwise.

width=1200
height=630
output=docs/og-image.png

browser=$(ls -d "$HOME"/Library/Caches/ms-playwright/chromium_headless_shell-*/chrome-mac/headless_shell 2>/dev/null | tail -1 || true)
if [ -z "$browser" ] || [ ! -x "$browser" ]; then
	browser="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
fi
if [ ! -x "$browser" ]; then
	echo "no headless Chromium found; install Google Chrome or Playwright's browsers" >&2
	exit 1
fi

# The card loads the stylesheet and the wordmark by relative path, so it is
# served rather than opened from the filesystem, where a browser refuses those.
port_file=$(mktemp)
python3 - "$port_file" <<'PYEOF' &
import sys
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


class Quiet(SimpleHTTPRequestHandler):
    """A request log would be the loudest thing this script prints."""

    def log_message(self, *_):
        pass


server = ThreadingHTTPServer(("127.0.0.1", 0), Quiet)
with open(sys.argv[1], "w") as handle:
    handle.write(str(server.server_address[1]))
server.serve_forever()
PYEOF
server=$!
trap 'kill "$server" 2>/dev/null || true; wait "$server" 2>/dev/null || true; rm -f "$port_file"' EXIT

# Wait for the port rather than sleeping for a guess.
attempts=0
while [ ! -s "$port_file" ] && [ "$attempts" -lt 50 ]; do
	sleep 0.1
	attempts=$((attempts + 1))
done
port=$(cat "$port_file")
if [ -z "$port" ]; then
	echo "the preview server did not start" >&2
	exit 1
fi

"$browser" --headless --disable-gpu --hide-scrollbars --no-sandbox \
	--virtual-time-budget=4000 \
	--window-size="$width,$height" \
	--screenshot="$output" \
	"http://127.0.0.1:$port/docs/og/card.html" > /dev/null 2>&1

if [ ! -s "$output" ]; then
	echo "the card did not render" >&2
	exit 1
fi

# A share card is fetched by whoever renders the preview, so its weight is worth
# the one pass. pngquant is optional: without it the card is simply larger.
if command -v pngquant > /dev/null 2>&1; then
	before=$(wc -c < "$output")
	pngquant --quality 70-95 --speed 1 --force --output "$output" "$output"
	after=$(wc -c < "$output")
	echo "wrote $output ($before to $after bytes)"
else
	echo "wrote $output"
fi
