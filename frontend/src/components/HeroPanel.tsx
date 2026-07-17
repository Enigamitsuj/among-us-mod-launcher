import { AnimatePresence, motion } from "framer-motion";
import { Browser } from "@wailsio/runtime";
import type { Mod, Release } from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";
import { formatReleaseNotes, type NoteBlock } from "../lib/releaseNotes";

const LAUNCHER_GITHUB_URL = "https://github.com/Enigamitsuj/among-us-mod-launcher";

type Props = {
  mod: Mod | null;
  release: Release | null;
};

export function HeroPanel({ mod, release }: Props) {
  const from = mod?.accentFrom ?? "#a855f7";
  const to = mod?.accentTo ?? "#ef4444";
  const notes = formatReleaseNotes(release?.body ?? "");
  const published = formatDate(release?.publishedAt);

  const openGitHub = () => {
    void Browser.OpenURL(LAUNCHER_GITHUB_URL);
  };

  return (
    <section className="relative flex h-full min-h-0 flex-col overflow-hidden rounded-2xl border border-white/10 bg-panel">
      <AnimatedBackdrop from={from} to={to} />

      <div className="relative z-10 flex h-full min-h-0 flex-col p-5">
        <header className="relative shrink-0 overflow-hidden rounded-2xl border border-white/10">
          {/* Crew art as top-section background */}
          <div className="absolute inset-0">
            <img
              src="/crew-banner.png"
              alt=""
              className="h-full w-full object-cover object-[center_55%]"
              draggable={false}
            />
            <div
              className="absolute inset-0"
              style={{
                background: `
                  linear-gradient(105deg, rgba(8,8,14,0.88) 0%, rgba(8,8,14,0.72) 42%, rgba(8,8,14,0.28) 72%, rgba(8,8,14,0.45) 100%),
                  linear-gradient(180deg, rgba(10,10,16,0.15) 0%, transparent 40%, rgba(10,10,16,0.55) 100%),
                  linear-gradient(90deg, ${from}33 0%, transparent 35%, ${to}22 100%)
                `,
              }}
            />
          </div>

          <div className="relative px-4 py-4">
            <motion.p
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              className="text-[11px] font-semibold uppercase tracking-[0.2em] text-purple/90"
            >
              Community Installer
            </motion.p>
            <motion.h1
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.04 }}
              className="mt-1 max-w-[58%] text-3xl font-bold tracking-tight text-white drop-shadow-[0_2px_12px_rgba(0,0,0,0.65)]"
            >
              {mod?.name ?? "Among Us Mod Launcher"}
            </motion.h1>
            <p className="mt-1.5 max-w-[55%] text-[13px] text-white/75 drop-shadow-[0_1px_8px_rgba(0,0,0,0.55)]">
              MIRA HQ atmosphere · launch, install, play
            </p>
          </div>
        </header>

        <AnimatePresence mode="wait">
          <motion.div
            key={release?.tagName ?? "empty"}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -6 }}
            transition={{ duration: 0.22 }}
            className="mt-4 flex min-h-0 flex-1 flex-col overflow-hidden rounded-2xl border border-white/10 bg-black/45 shadow-[0_12px_40px_rgba(0,0,0,0.35)] backdrop-blur-md"
          >
            <div className="flex shrink-0 items-start justify-between gap-3 border-b border-white/10 px-4 py-3">
              <div className="min-w-0">
                <div className="text-[10px] font-semibold uppercase tracking-[0.16em] text-muted">
                  Patch notes
                </div>
                <div className="mt-1 truncate text-sm font-semibold text-white">
                  {release ? release.name || release.tagName : "Select a version"}
                </div>
              </div>
              <div className="flex shrink-0 flex-col items-end gap-1">
                {release && (
                  <span className="rounded-md bg-purple/20 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-purple">
                    {release.tagName}
                  </span>
                )}
                {published && (
                  <span className="text-[10px] text-muted">{published}</span>
                )}
              </div>
            </div>

            <div className="scrollbar-thin min-h-0 flex-1 overflow-y-auto px-4 py-3">
              {release ? (
                notes.length > 0 ? (
                  <div className="space-y-2.5 text-[13px] leading-relaxed text-white/85">
                    {notes.map((block, i) => (
                      <NoteBlockView key={i} block={block} />
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-muted">No patch notes were published for this release.</p>
                )
              ) : (
                <p className="text-sm text-muted">
                  Choose a version on the right to preview its changelog here.
                </p>
              )}
            </div>
          </motion.div>
        </AnimatePresence>

        <footer className="relative z-20 mt-3 shrink-0 rounded-xl border border-white/5 bg-black/35 px-3 py-2 backdrop-blur-sm">
          <p className="text-xs text-white/75">
            Installer created by{" "}
            <button
              type="button"
              onClick={openGitHub}
              className="font-medium text-purple underline decoration-purple/40 underline-offset-2 transition hover:text-purple/90 hover:decoration-purple"
              title={LAUNCHER_GITHUB_URL}
            >
              Enigamitsuj
            </button>
          </p>
          <p className="text-[11px] text-muted">
            Community-made launcher for Town of Us: Mira and friends.
          </p>
        </footer>
      </div>
    </section>
  );
}

function AnimatedBackdrop({ from, to }: { from: string; to: string }) {
  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
      <div
        className="absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse at 18% 12%, ${from}66, transparent 42%),
            radial-gradient(ellipse at 88% 22%, ${to}44, transparent 38%),
            radial-gradient(ellipse at 50% 100%, #1e1b4bcc 0%, #0a0a10 68%),
            linear-gradient(165deg, #0b1020, #160b1d 55%, #0a0a10)
          `,
        }}
      />

      <div
        className="absolute inset-0 opacity-[0.12]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.08) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.08) 1px, transparent 1px)",
          backgroundSize: "42px 42px",
          maskImage: "radial-gradient(ellipse at center, black 35%, transparent 78%)",
        }}
      />

      <motion.div
        className="absolute -left-16 top-10 h-56 w-56 rounded-full blur-3xl"
        style={{ background: from }}
        animate={{ x: [0, 28, 0], y: [0, 18, 0], opacity: [0.22, 0.34, 0.22] }}
        transition={{ duration: 14, repeat: Infinity, ease: "easeInOut" }}
      />
      <motion.div
        className="absolute -right-10 bottom-16 h-64 w-64 rounded-full blur-3xl"
        style={{ background: to }}
        animate={{ x: [0, -22, 0], y: [0, -16, 0], opacity: [0.18, 0.3, 0.18] }}
        transition={{ duration: 16, repeat: Infinity, ease: "easeInOut" }}
      />
      <motion.div
        className="absolute left-1/3 top-1/2 h-40 w-40 -translate-y-1/2 rounded-full bg-cyan-400/20 blur-3xl"
        animate={{ scale: [1, 1.18, 1], opacity: [0.12, 0.22, 0.12] }}
        transition={{ duration: 10, repeat: Infinity, ease: "easeInOut" }}
      />

      {STARS.map((s) => (
        <motion.span
          key={s.id}
          className="absolute rounded-full bg-white"
          style={{
            left: s.x,
            top: s.y,
            width: s.size,
            height: s.size,
            boxShadow: `0 0 ${s.size * 3}px rgba(255,255,255,0.7)`,
          }}
          animate={{ opacity: [0.25, 0.95, 0.25] }}
          transition={{ duration: s.duration, delay: s.delay, repeat: Infinity, ease: "easeInOut" }}
        />
      ))}

      <svg
        className="absolute bottom-14 left-1/2 w-[130%] -translate-x-1/2 opacity-30"
        viewBox="0 0 800 180"
        fill="none"
      >
        <path
          d="M0 140 L80 120 L140 130 L220 90 L300 110 L380 70 L460 100 L560 60 L640 95 L720 80 L800 110 L800 180 L0 180 Z"
          fill="#1f2937"
        />
        <rect x="360" y="40" width="80" height="70" rx="6" fill="#334155" />
        <motion.rect
          x="380"
          y="55"
          width="18"
          height="18"
          rx="2"
          fill="#67e8f9"
          animate={{ opacity: [0.45, 0.95, 0.45] }}
          transition={{ duration: 2.4, repeat: Infinity }}
        />
        <motion.rect
          x="410"
          y="55"
          width="18"
          height="18"
          rx="2"
          fill="#c084fc"
          animate={{ opacity: [0.95, 0.4, 0.95] }}
          transition={{ duration: 2.8, repeat: Infinity }}
        />
      </svg>

      <div className="absolute inset-x-0 bottom-0 h-36 bg-gradient-to-t from-black/80 via-black/35 to-transparent" />
    </div>
  );
}

function NoteBlockView({ block }: { block: NoteBlock }) {
  if (block.type === "h") {
    return <h3 className="pt-1 text-[13px] font-semibold text-white">{block.text}</h3>;
  }
  if (block.type === "li") {
    return (
      <div className="flex gap-2 text-white/80">
        <span className="mt-[7px] h-1.5 w-1.5 shrink-0 rounded-full bg-purple/80" />
        <span>{block.text}</span>
      </div>
    );
  }
  return <p className="text-white/80">{block.text}</p>;
}

function formatDate(iso?: string) {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" });
}

const STARS = [
  { id: 1, x: "8%", y: "14%", size: 2, duration: 3.2, delay: 0.2 },
  { id: 2, x: "22%", y: "8%", size: 1.5, duration: 4.1, delay: 0.8 },
  { id: 3, x: "41%", y: "18%", size: 2, duration: 3.6, delay: 1.4 },
  { id: 4, x: "63%", y: "11%", size: 1.5, duration: 4.8, delay: 0.5 },
  { id: 5, x: "78%", y: "20%", size: 2.5, duration: 3.9, delay: 1.1 },
  { id: 6, x: "91%", y: "9%", size: 1.5, duration: 4.4, delay: 1.8 },
  { id: 7, x: "15%", y: "32%", size: 1.5, duration: 5.1, delay: 0.3 },
  { id: 8, x: "55%", y: "28%", size: 2, duration: 3.4, delay: 2.1 },
  { id: 9, x: "84%", y: "36%", size: 1.5, duration: 4.2, delay: 0.9 },
];
