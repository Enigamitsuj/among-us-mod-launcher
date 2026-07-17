export type NoteBlock =
  | { type: "h"; text: string }
  | { type: "p"; text: string }
  | { type: "li"; text: string };

/** Lightweight GitHub release-body parser (markdown-ish → display blocks). */
export function formatReleaseNotes(raw: string): NoteBlock[] {
  const text = raw.replace(/\r\n/g, "\n").trim();
  if (!text) return [];

  const lines = text.split("\n");
  const blocks: NoteBlock[] = [];
  let paragraph: string[] = [];

  const flushParagraph = () => {
    if (paragraph.length === 0) return;
    const joined = cleanInline(paragraph.join(" ").trim());
    if (joined) blocks.push({ type: "p", text: joined });
    paragraph = [];
  };

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed) {
      flushParagraph();
      continue;
    }

    const heading = trimmed.match(/^#{1,3}\s+(.+)$/);
    if (heading) {
      flushParagraph();
      blocks.push({ type: "h", text: cleanInline(heading[1]) });
      continue;
    }

    const bullet = trimmed.match(/^[-*+]\s+(.+)$/) || trimmed.match(/^\d+\.\s+(.+)$/);
    if (bullet) {
      flushParagraph();
      blocks.push({ type: "li", text: cleanInline(bullet[1]) });
      continue;
    }

    paragraph.push(trimmed);
  }

  flushParagraph();
  return blocks.slice(0, 80);
}

function cleanInline(input: string): string {
  return input
    .replace(/!\[[^\]]*]\([^)]+\)/g, "")
    .replace(/\[([^\]]+)]\([^)]+\)/g, "$1")
    .replace(/`([^`]+)`/g, "$1")
    .replace(/\*\*([^*]+)\*\*/g, "$1")
    .replace(/__([^_]+)__/g, "$1")
    .replace(/\*([^*]+)\*/g, "$1")
    .replace(/_([^_]+)_/g, "$1")
    .replace(/^>\s?/gm, "")
    .replace(/\s+/g, " ")
    .trim();
}
