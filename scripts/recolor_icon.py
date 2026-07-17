"""Rebuild the launcher icon assets from a checked-in source sprite.

Usage (from repo root):
  python scripts/recolor_icon.py
  python scripts/recolor_icon.py path/to/source.png
"""

from __future__ import annotations

import sys
from pathlib import Path

from PIL import Image

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_SRC = ROOT / "scripts" / "assets" / "dead-crewmate-source.png"
OUT = ROOT / "build" / "appicon.png"
PUBLIC = ROOT / "frontend" / "public" / "app-icon.png"

BG = (10, 10, 16, 255)
BODY = (255, 214, 53, 255)
BODY_SHADE = (230, 176, 20, 255)
BONE = (255, 240, 190, 255)
BONE_SHADE = (232, 208, 140, 255)
BLOOD = (200, 12, 24, 255)
BLOOD_DARK = (150, 6, 14, 255)
OUTLINE = (0, 0, 0, 255)


def main() -> int:
    src_path = Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else DEFAULT_SRC
    if not src_path.is_file():
        print(
            f"Source image not found: {src_path}\n"
            "Pass a PNG path, or place one at scripts/assets/dead-crewmate-source.png.\n"
            "Committed icons already live at build/appicon.png and frontend/public/app-icon.png.",
            file=sys.stderr,
        )
        return 1

    src = Image.open(src_path).convert("RGBA")
    w, h = src.size
    px = src.load()

    mask = [[False] * w for _ in range(h)]
    for y in range(h):
        for x in range(w):
            _, _, _, a = px[x, y]
            if a > 20:
                mask[y][x] = True

    ys = [y for y in range(h) if any(mask[y])]
    xs_all = [x for y in ys for x in range(w) if mask[y][x]]
    min_y, max_y = min(ys), max(ys)
    min_x, max_x = min(xs_all), max(xs_all)
    full_w = max_x - min_x + 1
    center_x = (min_x + max_x) / 2

    row_span = {}
    for y in range(min_y, max_y + 1):
        xs = [x for x in range(w) if mask[y][x]]
        if xs:
            row_span[y] = (min(xs), max(xs), max(xs) - min(xs) + 1)

    body_start = None
    for y in range(min_y, max_y + 1):
        if y in row_span and row_span[y][2] >= full_w * 0.48:
            body_start = y
            break
    if body_start is None:
        body_start = min_y + int((max_y - min_y) * 0.35)

    blood_depth = max(10, int((max_y - min_y) * 0.055))

    colored = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    cp = colored.load()

    for y in range(h):
        for x in range(w):
            if not mask[y][x]:
                continue

            edge = False
            for dx, dy in ((-1, 0), (1, 0), (0, -1), (0, 1)):
                nx, ny = x + dx, y + dy
                if nx < 0 or ny < 0 or nx >= w or ny >= h or not mask[ny][nx]:
                    edge = True
                    break

            if y < body_start:
                cp[x, y] = OUTLINE if edge else (BONE_SHADE if x >= center_x else BONE)
                continue

            if y <= body_start + blood_depth:
                near_stump = abs(x - center_x) <= full_w * 0.2
                on_top_edge = y <= body_start + max(4, blood_depth // 3)
                if near_stump or on_top_edge:
                    if edge and not near_stump:
                        cp[x, y] = OUTLINE
                    else:
                        cp[x, y] = BLOOD if y <= body_start + blood_depth // 2 else BLOOD_DARK
                    continue

            if edge:
                cp[x, y] = OUTLINE
            elif x > center_x + full_w * 0.1:
                cp[x, y] = BODY_SHADE
            else:
                cp[x, y] = BODY

    bbox = colored.getbbox()
    content = colored.crop(bbox)
    pad = 28
    padded = Image.new("RGBA", (content.width + pad * 2, content.height + pad * 2), (0, 0, 0, 0))
    padded.paste(content, (pad, pad), content)

    canvas = Image.new("RGBA", (512, 512), BG)
    scale = min(450 / padded.width, 450 / padded.height)
    nw = max(1, int(padded.width * scale))
    nh = max(1, int(padded.height * scale))
    scaled = padded.resize((nw, nh), Image.Resampling.NEAREST)
    canvas.paste(scaled, ((512 - nw) // 2, (512 - nh) // 2), scaled)

    OUT.parent.mkdir(parents=True, exist_ok=True)
    PUBLIC.parent.mkdir(parents=True, exist_ok=True)
    canvas.save(OUT)
    canvas.save(PUBLIC)
    print(f"body_start={body_start} blood_depth={blood_depth} wrote {OUT} and {PUBLIC}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
