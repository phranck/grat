#!/bin/sh
set -eu

# Render the icon into the raster sizes browsers and platforms still ask for.
#
# The SVG is the source and everything here is derived from it, so a change to
# the mark reaches every size by running this again. Only the sizes that are
# genuinely needed are produced: a browser that understands an SVG favicon uses
# that one, and the rest exist for the places that do not.
#
#   favicon.ico        16, 32 and 48 in one file, for a browser that asks for it
#                      by that name and for a bookmark bar
#   apple-touch-icon   180, for a home screen on iOS, which applies its own
#                      rounding and wants no transparency
#   icon-192, icon-512 for a web app manifest on Android
#   icon-maskable      512 without a drawn corner, since the platform cuts one
#
# It needs a headless Chromium to rasterise, and Pillow to write the .ico.

browser=$(ls -d "$HOME"/Library/Caches/ms-playwright/chromium_headless_shell-*/chrome-mac/headless_shell 2>/dev/null | tail -1 || true)
if [ -z "$browser" ] || [ ! -x "$browser" ]; then
	browser="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
fi
if [ ! -x "$browser" ]; then
	echo "no headless Chromium found; install Google Chrome or Playwright's browsers" >&2
	exit 1
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

# The icon is rendered from a page rather than from the file, because a browser
# screenshots a viewport and the page is what gives it an exact size with no
# margin around the mark.
render() {
	source_svg=$1
	size=$2
	output=$3
	cat > "$work/page.html" <<HTML
<!DOCTYPE html><html><head><meta charset="utf-8"><style>
  html,body{margin:0;padding:0;width:${size}px;height:${size}px;overflow:hidden;background:transparent}
  img{display:block;width:${size}px;height:${size}px}
</style></head><body><img src="icon.svg" alt=""></body></html>
HTML
	cp "$source_svg" "$work/icon.svg"

	port_file="$work/port"
	rm -f "$port_file"
	python3 - "$work" "$port_file" <<'PYEOF' &
import sys
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


class Quiet(SimpleHTTPRequestHandler):
    def log_message(self, *_):
        pass


handler = partial(Quiet, directory=sys.argv[1])
server = ThreadingHTTPServer(("127.0.0.1", 0), handler)
with open(sys.argv[2], "w") as handle:
    handle.write(str(server.server_address[1]))
server.serve_forever()
PYEOF
	server=$!

	attempts=0
	while [ ! -s "$port_file" ] && [ "$attempts" -lt 50 ]; do
		sleep 0.1
		attempts=$((attempts + 1))
	done
	port=$(cat "$port_file")

	"$browser" --headless --disable-gpu --hide-scrollbars --no-sandbox \
		--virtual-time-budget=3000 --window-size="$size,$size" \
		--default-background-color=00000000 \
		--screenshot="$output" "http://127.0.0.1:$port/page.html" > /dev/null 2>&1

	kill "$server" 2>/dev/null || true
	wait "$server" 2>/dev/null || true

	if [ ! -s "$output" ]; then
		echo "could not render $output" >&2
		exit 1
	fi
	echo "wrote $output"
}

# iOS composites a home screen icon on black and rounds it itself, so this one
# takes the full square rather than the plate with corners cut out of it. A
# rounded source would show black in the corners.
render docs/icon-maskable.svg 180 docs/apple-touch-icon.png
render docs/favicon.svg 192 docs/icon-192.png
render docs/favicon.svg 512 docs/icon-512.png
render docs/icon-maskable.svg 512 docs/icon-maskable.png

# The .ico carries three sizes, and each is rasterised at the size it will be
# seen at rather than shrunk from the large one. Downscaling loses the plate's
# rounded corners entirely at sixteen pixels, where three or four pixels of
# radius have to be drawn deliberately to survive.
for size in 16 32 48; do
	render docs/favicon.svg "$size" "$work/ico-$size.png"
done

python3 - "$work" <<'PYEOF'
import sys

from PIL import Image

work = sys.argv[1]
frames = [Image.open(f"{work}/ico-{size}.png").convert("RGBA") for size in (16, 32, 48)]
# The largest carries the file and the others ride along, which is how a browser
# is offered a choice rather than left to scale.
frames[-1].save("docs/favicon.ico", format="ICO", sizes=[(16, 16), (32, 32), (48, 48)], append_images=frames[:-1])
print("wrote docs/favicon.ico")
PYEOF

if command -v pngquant > /dev/null 2>&1; then
	for image in docs/apple-touch-icon.png docs/icon-192.png docs/icon-512.png docs/icon-maskable.png; do
		pngquant --quality 70-95 --speed 1 --force --output "$image" "$image"
	done
	echo "compressed the raster icons"
fi
