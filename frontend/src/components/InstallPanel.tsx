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
  installLocation: string;
  createShortcut: boolean;
  launchAfter: boolean;
  onToggleShortcut: (v: boolean) => void;
  onToggleLaunch: (v: boolean) => void;
  onBrowseInstallLocation: () => void;
  installState: InstallState | null;
  progress: InstallProgress | null;
  busy: boolean;
  error: string | null;
  onInstall: () => void;
  onPlay: () => void;
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
    installLocation,
    createShortcut,
    launchAfter,
    onToggleShortcut,
    onToggleLaunch,
    onBrowseInstallLocation,
    installState,
    progress,
    busy,
    error,
    onInstall,
    onPlay,
  } = props;

  const installed = Boolean(installState?.installed);
  const selectedRelease = releases.find((r) => r.tagName === selectedTag) ?? null;
  const releaseInstallable = selectedRelease ? selectedRelease.installable : false;

  const canInstall =
    Boolean(selectedMod?.enabled) &&
    Boolean(game?.found) &&
    !game?.running &&
    Boolean(selectedTag) &&
    releaseInstallable &&
    Boolean(installLocation) &&
    !busy;

  const installTarget = installLocation
    ? `${installLocation}\\${selectedMod?.folderName ?? ""}`
    : "Next to launcher";

  // Status is driven by the selected release's compatibility with the game.
  const statusText = error
    ? error
    : busy && progress
      ? progress.message
      : installed
        ? "Ready to play"
        : selectedRelease?.compatReason
          ? selectedRelease.compatReason
          : game?.message ?? "Checking...";

  const statusTone: "ok" | "warn" | "bad" = error
    ? "bad"
    : !game?.found || game.running
      ? "bad"
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
          Community Installer
        </p>
        <h2 className="mt-1 text-xl font-bold tracking-tight text-white">
          Install {selectedMod?.shortName ?? "Mod"}
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
              <div className="truncate text-sm font-medium">{mod.shortName}</div>
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
            disabled={busy}
            onChange={onSelectTag}
          />
        </Field>

        <Field label="Among Us install">
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
          <div className="flex items-stretch gap-2">
            <div className="min-w-0 flex-1">
              <PathBox value={installTarget} tone="neutral" />
            </div>
            <button
              type="button"
              disabled={busy}
              onClick={onBrowseInstallLocation}
              className="shrink-0 rounded-xl border border-line bg-panel-2 px-3.5 text-sm font-medium text-white/90 transition hover:border-purple/50 hover:bg-purple/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
              title="Choose a different parent folder"
            >
              Browse
            </button>
          </div>
          <p className="mt-1 text-[11px] text-muted">
            Defaults next to the launcher. Browse to change the parent folder.
          </p>
        </Field>

        <div className="grid grid-cols-2 gap-2">
          <ToggleChip
            checked={createShortcut}
            disabled={busy}
            onChange={onToggleShortcut}
            label="Desktop shortcut"
          />
          <ToggleChip
            checked={launchAfter}
            disabled={busy}
            onChange={onToggleLaunch}
            label="Launch after install"
          />
        </div>
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
            {!busy && !installed && statusTone === "ok" && (
              <span className="shrink-0 rounded-md bg-emerald-500/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-emerald-300">
                Ready
              </span>
            )}
            {!busy && !installed && statusTone === "warn" && (
              <span className="shrink-0 rounded-md bg-amber-500/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-amber-300">
                Check
              </span>
            )}
            {!busy && !installed && statusTone === "bad" && game?.found && (
              <span className="shrink-0 rounded-md bg-red/15 px-2 py-1 text-[10px] font-semibold uppercase tracking-wide text-red">
                Blocked
              </span>
            )}
          </div>
          {busy && progress && (
            <div className="mt-2.5">
              <ProgressBar percent={progress.percent} message="" stage={progress.stage} />
            </div>
          )}
        </div>

        <div className="flex gap-2">
          {installed ? (
            <ActionButton onClick={onPlay} disabled={busy} variant="play">
              Play
            </ActionButton>
          ) : (
            <ActionButton onClick={onInstall} disabled={!canInstall} variant="install">
              {busy ? "Installing..." : "Install"}
            </ActionButton>
          )}
          {installed && (
            <ActionButton onClick={onInstall} disabled={busy} variant="secondary">
              Reinstall
            </ActionButton>
          )}
        </div>
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

function ToggleChip({
  checked,
  disabled,
  onChange,
  label,
}: {
  checked: boolean;
  disabled?: boolean;
  onChange: (v: boolean) => void;
  label: string;
}) {
  return (
    <label
      className={[
        "flex cursor-pointer items-center gap-2.5 rounded-xl border px-3 py-2.5 text-[13px] transition",
        checked
          ? "border-purple/40 bg-purple/10 text-white"
          : "border-line bg-panel-2 text-muted hover:border-white/15 hover:text-white",
        disabled ? "cursor-not-allowed opacity-50" : "",
      ].join(" ")}
    >
      <input
        type="checkbox"
        className="accent-purple"
        checked={checked}
        disabled={disabled}
        onChange={(e) => onChange(e.target.checked)}
      />
      <span className="leading-tight">{label}</span>
    </label>
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
  variant: "install" | "play" | "secondary";
}) {
  const styles =
    variant === "secondary"
      ? "border border-line bg-panel-2 text-white hover:border-white/20"
      : variant === "play"
        ? "bg-gradient-to-r from-emerald-500 to-cyan-500 text-white glow-purple"
        : "bg-gradient-to-r from-purple to-red text-white glow-purple";

  return (
    <motion.button
      type="button"
      whileHover={disabled ? undefined : { scale: 1.015 }}
      whileTap={disabled ? undefined : { scale: 0.985 }}
      disabled={disabled}
      onClick={onClick}
      className={[
        "flex-1 rounded-xl px-4 py-3 text-sm font-semibold transition",
        styles,
        disabled ? "cursor-not-allowed opacity-50" : "hover:brightness-110",
      ].join(" ")}
    >
      {children}
    </motion.button>
  );
}
