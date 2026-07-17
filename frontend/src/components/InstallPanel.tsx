import { type ReactNode } from "react";
import { motion } from "framer-motion";
import type {
  GameStatus,
  InstallProgress,
  InstallState,
  Mod,
  Release,
} from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";
import { ProgressBar } from "./ProgressBar";
import { VersionSelect } from "./VersionSelect";

type Props = {
  mods: Mod[];
  selectedMod: Mod | null;
  onSelectMod: (mod: Mod) => void;
  game: GameStatus | null;
  releases: Release[];
  selectedTag: string;
  onSelectTag: (tag: string) => void;
  installedModIds: Set<string>;
  installState: InstallState | null;
  progress: InstallProgress | null;
  busy: boolean;
  error: string | null;
  onInstall: () => void;
  onPlay: () => void;
  onUninstall: () => void;
};

export function InstallPanel(props: Props) {
  const {
    mods,
    selectedMod,
    onSelectMod,
    game,
    releases,
    selectedTag,
    onSelectTag,
    installedModIds,
    installState,
    progress,
    busy,
    error,
    onInstall,
    onPlay,
    onUninstall,
  } = props;

  const installed = Boolean(installState?.installed);
  const selectedIsInstalled =
    installed && (!installState?.version || installState.version === selectedTag);
  const selectedRelease = releases.find((r) => r.tagName === selectedTag) ?? null;
  const releaseInstallable = selectedRelease ? selectedRelease.installable : false;
  const installedReleaseIndex = installState?.version
    ? releases.findIndex((r) => r.tagName === installState.version)
    : -1;
  const updateRelease =
    installed && installedReleaseIndex > 0 ? releases[0] : null;

  const canInstall =
    Boolean(selectedMod?.enabled) &&
    Boolean(game?.found) &&
    !game?.running &&
    Boolean(selectedTag) &&
    releaseInstallable &&
    !busy;

  const installTarget = installState?.path || "Next to the detected Among Us installation";

  // Status is driven by the selected release's compatibility with the game.
  const statusText = error
    ? error
    : busy && progress
      ? progress.message
      : game?.running
        ? "Among Us is running. Close it to manage or launch mods."
        : selectedIsInstalled
          ? `Ready to play · ${installState?.version}`
          : installed && installState?.version
            ? `Installed ${installState.version}. Select Install to switch to ${selectedTag || "this version"}.`
            : selectedRelease?.compatReason
              ? selectedRelease.compatReason
              : game?.message ?? "Checking...";

  const statusTone: "ok" | "warn" | "bad" = error
    ? "bad"
    : !game?.found || game.running
      ? "bad"
      : selectedIsInstalled
        ? "ok"
        : selectedRelease?.compatLevel === "unsupported"
          ? "bad"
          : selectedRelease?.compatLevel === "caution"
            ? "warn"
            : releaseInstallable
              ? "ok"
              : "warn";

  const versionLabel =
    game?.version && game.version !== "unknown" ? `v${game.version.replace(/^v/i, "")}` : "version unknown";

  return (
    <section className="flex h-full min-h-0 flex-col overflow-hidden rounded-2xl border border-white/10 bg-panel/90 p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
      {/* Identity */}
      <header className="shrink-0">
        <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-purple">
          Mod Hub
        </p>
        <h2 className="mt-1 text-xl font-bold tracking-tight text-white">
          {selectedIsInstalled ? "Play" : "Install"} {selectedMod?.shortName ?? "Mod"}
        </h2>
        <p className="mt-1.5 line-clamp-2 text-[13px] leading-snug text-muted">
          {selectedMod?.description ?? "Select a mod to get started."}
        </p>
      </header>

      {/* Mod switcher */}
      <div className="mt-3 flex shrink-0 gap-2">
        {mods.map((mod) => {
          const active = selectedMod?.id === mod.id;
          return (
            <button
              key={mod.id}
              type="button"
              disabled={mod.comingSoon}
              onClick={() => onSelectMod(mod)}
              className={[
                "min-w-0 flex-1 rounded-xl border px-3 py-2 text-left transition",
                active
                  ? "border-purple/60 bg-purple/15 text-white glow-purple"
                  : "border-line bg-panel-2 text-muted hover:border-white/20 hover:text-white",
                mod.comingSoon ? "cursor-not-allowed opacity-50" : "",
              ].join(" ")}
            >
              <div className="flex items-center justify-between gap-2">
                <div className="truncate text-sm font-medium">{mod.shortName}</div>
                {installedModIds.has(mod.id) && (
                  <span
                    className="grid h-5 w-5 shrink-0 place-items-center rounded-full bg-emerald-500/20 text-xs font-bold text-emerald-300"
                    title="Installed"
                    aria-label="Installed"
                  >
                    ✓
                  </span>
                )}
              </div>
              <div className="truncate text-[11px] opacity-70">
                {mod.comingSoon ? "Coming soon" : mod.author}
              </div>
            </button>
          );
        })}
      </div>

      {/* Form — no scrolling; denser single-page layout */}
      <div className="mt-3 flex min-h-0 flex-1 flex-col gap-2.5 overflow-hidden">
        <Field label="Version">
          <VersionSelect
            releases={releases}
            value={selectedTag}
            installedTag={installState?.version || ""}
            disabled={busy}
            onChange={onSelectTag}
          />
        </Field>

        <Field label="Your Among Us install">
          <div
            className={[
              "rounded-xl border bg-panel-2 px-3 py-2",
              game?.found ? (game.canInstall ? "border-line" : "border-amber-500/40") : "border-red/40",
            ].join(" ")}
          >
            <div className="flex flex-wrap items-center gap-1.5">
              <MetaChip
                text={game?.platformLabel || "Not detected"}
                tone={game?.found ? "purple" : "muted"}
              />
              {game?.found && <MetaChip text={versionLabel} tone={game.canInstall ? "green" : "amber"} />}
              {game?.branch && game.branch !== "public" && (
                <MetaChip text={`branch: ${game.branch}`} tone="muted" />
              )}
              {game?.assetHint && game.found && (
                <MetaChip text={game.assetHint} tone="muted" />
              )}
            </div>
            <div
              className={[
                "mt-1.5 truncate text-[13px]",
                game?.found ? "text-white/90" : "text-red",
              ].join(" ")}
              title={game?.path || "Not found"}
            >
              {game?.path || "Not found"}
            </div>
          </div>
        </Field>

        <Field label="Install location">
          <PathBox value={installTarget} tone="neutral" />
          <p className="mt-1 text-[11px] text-muted">
            Managed automatically beside your original game. The original is never modified.
          </p>
        </Field>
      </div>

      {/* CTA footer — always visible */}
      <div className="mt-3 shrink-0 space-y-2.5 border-t border-white/5 pt-3">
        <div className="rounded-xl border border-line/80 bg-black/25 px-3 py-2.5">
          <div className="flex items-center justify-between gap-3">
            <div className="min-w-0">
              <div className="text-[10px] font-medium uppercase tracking-wide text-muted">Status</div>
              <div
                className={[
                  "text-sm leading-snug",
                  statusTone === "bad"
                    ? "text-red"
                    : statusTone === "warn"
                      ? "text-amber-300"
                      : "text-white/90",
                ].join(" ")}
                title={statusText}
              >
                {statusText}
              </div>
            </div>
            {!busy && updateRelease && selectedIsInstalled ? (
              <motion.button
                type="button"
                onClick={() => onSelectTag(updateRelease.tagName)}
                initial={false}
                animate={{
                  boxShadow: [
                    "0 0 6px rgba(168,85,247,0.18), inset 0 0 0 rgba(168,85,247,0)",
                    "0 0 14px rgba(168,85,247,0.5), inset 0 0 8px rgba(168,85,247,0.12)",
                  ],
                  borderColor: [
                    "rgba(168,85,247,0.4)",
                    "rgba(168,85,247,0.85)",
                  ],
                  backgroundColor: [
                    "rgba(168,85,247,0.1)",
                    "rgba(168,85,247,0.22)",
                  ],
                }}
                transition={{
                  duration: 2.4,
                  repeat: Infinity,
                  repeatType: "mirror",
                  ease: "easeInOut",
                }}
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.97 }}
                className="relative shrink-0 cursor-pointer overflow-hidden inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wide text-purple"
                title={`Switch to ${updateRelease.tagName}`}
                aria-label={`Update available. Select ${updateRelease.tagName}`}
              >
                <span className="update-stripes" aria-hidden="true" />
                <motion.span
                  className="relative z-10 h-1.5 w-1.5 rounded-[2px] bg-purple"
                  animate={{
                    boxShadow: [
                      "0 0 2px rgba(168,85,247,0.3)",
                      "0 0 8px rgba(168,85,247,0.95)",
                    ],
                    opacity: [0.55, 1],
                  }}
                  transition={{
                    duration: 2.4,
                    repeat: Infinity,
                    repeatType: "mirror",
                    ease: "easeInOut",
                  }}
                />
                <span className="relative z-10">Update · {updateRelease.tagName}</span>
              </motion.button>
            ) : !busy && !selectedIsInstalled && statusTone === "ok" ? (
              <span className="shrink-0 rounded-md bg-emerald-500/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-emerald-300">
                Ready
              </span>
            ) : !busy && !selectedIsInstalled && statusTone === "warn" ? (
              <span className="shrink-0 rounded-md bg-amber-500/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-amber-300">
                Check
              </span>
            ) : !busy && !selectedIsInstalled && statusTone === "bad" && game?.found ? (
              <span className="shrink-0 rounded-md bg-red/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-red">
                Blocked
              </span>
            ) : null}
          </div>
          {busy && progress && (
            <div className="mt-2.5">
              <ProgressBar percent={progress.percent} message="" stage={progress.stage} />
            </div>
          )}
        </div>

        {selectedIsInstalled ? (
          <>
            <button
              type="button"
              onClick={onUninstall}
              disabled={busy || game?.running || !installState?.managed}
              title={
                installState?.managed
                  ? "Remove this modded copy"
                  : "Reinstall once before uninstalling this legacy copy"
              }
              className="w-full rounded-xl border border-red/30 bg-red/10 px-4 py-2.5 text-sm font-medium text-red transition hover:border-red/50 hover:bg-red/15 disabled:cursor-not-allowed disabled:opacity-50"
            >
              Uninstall
            </button>
            <ActionButton
              onClick={onPlay}
              disabled={busy || game?.running}
              variant={game?.running ? "running" : "play"}
            >
              {game?.running ? "Running..." : "Play"}
            </ActionButton>
          </>
        ) : (
          <ActionButton onClick={onInstall} disabled={!canInstall} variant="install">
            {busy ? "Installing..." : "Install"}
          </ActionButton>
        )}
      </div>
    </section>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="shrink-0 space-y-1">
      <span className="text-[10px] font-medium uppercase tracking-wide text-muted">{label}</span>
      {children}
    </div>
  );
}

function PathBox({
  value,
  tone = "neutral",
}: {
  value: string;
  tone?: "neutral" | "danger";
}) {
  return (
    <div
      className={[
        "rounded-xl border bg-panel-2 px-3 py-2",
        tone === "danger" ? "border-red/40" : "border-line",
      ].join(" ")}
    >
      <div
        className={[
          "truncate text-[13px]",
          tone === "danger" ? "text-red" : "text-white/90",
        ].join(" ")}
        title={value}
      >
        {value}
      </div>
    </div>
  );
}

function MetaChip({
  text,
  tone,
}: {
  text: string;
  tone: "purple" | "green" | "amber" | "muted";
}) {
  const styles =
    tone === "purple"
      ? "bg-purple/20 text-purple"
      : tone === "green"
        ? "bg-emerald-500/15 text-emerald-300"
        : tone === "amber"
          ? "bg-amber-500/15 text-amber-300"
          : "bg-white/5 text-muted";
  return (
    <span className={`rounded-md px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide ${styles}`}>
      {text}
    </span>
  );
}

function ActionButton({
  children,
  onClick,
  disabled,
  variant,
}: {
  children: ReactNode;
  onClick: () => void;
  disabled?: boolean;
  variant: "install" | "play" | "running";
}) {
  const styles =
    variant === "running"
      ? "bg-gradient-to-r from-emerald-600/80 to-cyan-600/80 text-white cursor-default"
      : variant === "play"
        ? "bg-gradient-to-r from-emerald-500 to-cyan-500 text-white glow-purple"
        : "bg-gradient-to-r from-purple to-red text-white glow-purple";

  return (
    <motion.button
      type="button"
      whileHover={disabled || variant === "running" ? undefined : { scale: 1.015 }}
      whileTap={disabled || variant === "running" ? undefined : { scale: 0.985 }}
      disabled={disabled}
      onClick={onClick}
      animate={
        variant === "running"
          ? {
              boxShadow: [
                "0 0 10px rgba(16,185,129,0.2)",
                "0 0 22px rgba(34,211,238,0.45)",
                "0 0 10px rgba(16,185,129,0.2)",
              ],
              opacity: [0.82, 1, 0.82],
            }
          : undefined
      }
      transition={
        variant === "running"
          ? { duration: 2.1, repeat: Infinity, ease: "easeInOut" }
          : undefined
      }
      className={[
        "relative w-full overflow-hidden rounded-xl px-4 py-3 text-sm font-semibold transition",
        styles,
        disabled && variant !== "running" ? "cursor-not-allowed opacity-50" : "",
        disabled && variant === "running" ? "opacity-100" : "",
        !disabled ? "hover:brightness-110" : "",
      ].join(" ")}
    >
      {variant === "running" && <span className="update-stripes opacity-25" aria-hidden="true" />}
      <span className="relative z-10 inline-flex items-center justify-center gap-2">
        {variant === "running" && (
          <motion.span
            className="h-1.5 w-1.5 rounded-full bg-white"
            animate={{ opacity: [0.35, 1, 0.35], scale: [0.85, 1.15, 0.85] }}
            transition={{ duration: 1.2, repeat: Infinity, ease: "easeInOut" }}
          />
        )}
        {children}
      </span>
    </motion.button>
  );
}
