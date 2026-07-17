import { useEffect, useId, useMemo, useRef, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import type { Release } from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";

type Props = {
  releases: Release[];
  value: string;
  installedTag?: string;
  disabled?: boolean;
  onChange: (tag: string) => void;
};

export function VersionSelect({ releases, value, installedTag, disabled, onChange }: Props) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement | null>(null);
  const listId = useId();

  const channels = useMemo(() => {
    const latest = releases[0];
    const stable = releases.find((r) => !r.prerelease);
    const beta = releases.find((r) => r.prerelease);
    return { latest, stable, beta };
  }, [releases]);

  const selected = releases.find((r) => r.tagName === value) ?? null;

  useEffect(() => {
    if (!open) return;
    const onPointerDown = (e: MouseEvent) => {
      if (!rootRef.current?.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const labelFor = (r: Release) => {
    const tags: string[] = [];
    if (channels.latest?.tagName === r.tagName) tags.push("Latest");
    if (!r.prerelease && channels.stable?.tagName === r.tagName) tags.push("Stable");
    if (r.prerelease) tags.push("Beta");
    return tags;
  };

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        disabled={disabled || releases.length === 0}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-controls={listId}
        onClick={() => setOpen((v) => !v)}
        className={[
          "flex w-full items-center gap-3 rounded-xl border bg-panel-2 px-3 py-2.5 text-left text-sm outline-none transition",
          open ? "border-purple/60" : "border-line hover:border-white/20",
          disabled || releases.length === 0 ? "cursor-not-allowed opacity-50" : "cursor-pointer",
        ].join(" ")}
      >
        <span className="min-w-0 flex-1 truncate text-white/90">
          {selected ? (
            <span className="inline-flex flex-wrap items-center gap-2">
              <CompatDot level={selected.compatLevel} />
              <span>{selected.tagName}</span>
              {selected.recommended && <Badge text="Recommended" />}
              {installedTag && selected.tagName === installedTag && <Badge text="Installed" />}
              {labelFor(selected).map((tag) => (
                <Badge key={tag} text={tag} />
              ))}
            </span>
          ) : (
            <span className="text-muted">Loading releases...</span>
          )}
        </span>
        <svg
          width="14"
          height="14"
          viewBox="0 0 14 14"
          fill="none"
          aria-hidden="true"
          className={[
            "shrink-0 text-muted transition-transform duration-200",
            open ? "rotate-180 text-purple" : "",
          ].join(" ")}
        >
          <path
            d="M3.5 5.25L7 8.75L10.5 5.25"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </button>

      <AnimatePresence>
        {open && (
          <motion.ul
            id={listId}
            role="listbox"
            initial={{ opacity: 0, y: -4, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -4, scale: 0.98 }}
            transition={{ duration: 0.14 }}
            className="scrollbar-thin absolute left-0 right-0 z-30 mt-2 max-h-52 overflow-y-auto rounded-xl border border-white/10 bg-[#14141e] p-1.5 shadow-[0_18px_40px_rgba(0,0,0,0.55)]"
          >
            {releases.map((r) => {
              const active = r.tagName === value;
              return (
                <li key={r.tagName} role="option" aria-selected={active}>
                  <button
                    type="button"
                    className={[
                      "flex w-full items-center justify-between gap-3 rounded-lg px-3 py-2.5 text-left text-sm transition",
                      active
                        ? "bg-purple/20 text-white"
                        : "text-white/85 hover:bg-white/5 hover:text-white",
                      r.compatLevel === "unsupported" ? "opacity-70" : "",
                    ].join(" ")}
                    onClick={() => {
                      onChange(r.tagName);
                      setOpen(false);
                    }}
                    title={r.compatReason}
                  >
                    <span className="flex min-w-0 items-center gap-2">
                      <CompatDot level={r.compatLevel} />
                      <span className="truncate font-medium">{r.tagName}</span>
                    </span>
                    <span className="flex shrink-0 items-center gap-1">
                      {installedTag && r.tagName === installedTag && <Badge text="Installed" />}
                      {r.recommended && <Badge text="Rec" />}
                      {labelFor(r).map((tag) => (
                        <Badge key={tag} text={tag} />
                      ))}
                    </span>
                  </button>
                </li>
              );
            })}
          </motion.ul>
        )}
      </AnimatePresence>
    </div>
  );
}

function CompatDot({ level }: { level?: string }) {
  const color =
    level === "compatible"
      ? "bg-emerald-400"
      : level === "caution"
        ? "bg-amber-400"
        : level === "unsupported"
          ? "bg-red"
          : "bg-white/30";
  const title =
    level === "compatible"
      ? "Compatible"
      : level === "caution"
        ? "Check compatibility"
        : level === "unsupported"
          ? "Not compatible"
          : "Unknown";
  return <span className={`h-2 w-2 shrink-0 rounded-full ${color}`} title={title} aria-label={title} />;
}

function Badge({ text }: { text: string }) {
  const tone =
    text === "Latest"
      ? "bg-purple/20 text-purple"
      : text === "Stable" || text === "Installed"
        ? "bg-emerald-500/15 text-emerald-300"
        : text === "Recommended" || text === "Rec"
          ? "bg-purple/20 text-purple"
          : "bg-amber-500/15 text-amber-300";

  return (
    <span className={`rounded-md px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide ${tone}`}>
      {text}
    </span>
  );
}
