import { useMemo, type ReactNode } from "react";
import { motion } from "framer-motion";
import type {
  GameStatus,
  InstallProgress,
  InstallState,
  Mod,
  Release,
} from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";
import { ProgressBar } from "./ProgressBar";

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
    installState,
    progress,
    busy,
    error,
    onInstall,
    onPlay,
  } = props;

  const channels = useMemo(() => {
    const latest = releases[0];
    const stable = releases.find((r) => !r.prerelease);
    const beta = releases.find((r) => r.prerelease);
    return { latest, stable, beta };
  }, [releases]);

  const installed = Boolean(installState?.installed);
  const canInstall =
    Boolean(selectedMod?.enabled) &&
    Boolean(game?.found) &&
    !game?.running &&
    Boolean(selectedTag) &&
    !busy;

  return (
    <section className="flex h-full flex-col rounded-2xl border border-white/10 bg-panel/90 p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
      <header className="mb-4">
        <p className="text-xs font-semibold uppercase tracking-[0.18em] text-purple">Community Installer</p>
        <h2 className="mt-1 text-2xl font-bold text-white">
          Install {selectedMod?.shortName ?? "Mod"}
        </h2>
        <p className="mt-2 text-sm leading-relaxed text-muted">
          {selectedMod?.description ?? "Select a mod to get started."}
        </p>
      </header>

      <div className="mb-4 flex flex-wrap gap-2">
        {mods.map((mod) => {
          const active = selectedMod?.id === mod.id;
          return (
            <button
              key={mod.id}
              type="button"
              disabled={mod.comingSoon}
              onClick={() => onSelectMod(mod)}
              className={[
                "rounded-xl border px-3 py-2 text-left text-sm transition",
                active
                  ? "border-purple/60 bg-purple/15 text-white glow-purple"
                  : "border-line bg-panel-2 text-muted hover:border-white/20 hover:text-white",
                mod.comingSoon ? "cursor-not-allowed opacity-50" : "",
              ].join(" ")}
            >
              <div className="font-medium">{mod.shortName}</div>
              <div className="text-[11px] opacity-70">{mod.comingSoon ? "Coming soon" : mod.author}</div>
            </button>
          );
        })}
      </div>

      <div className="scrollbar-thin flex-1 space-y-3 overflow-y-auto pr-1">
        <Field label="Version">
          <select
            className="w-full rounded-xl border border-line bg-panel-2 px-3 py-2.5 text-sm text-white outline-none transition focus:border-purple/60"
            value={selectedTag}
            disabled={!releases.length || busy}
            onChange={(e) => onSelectTag(e.target.value)}
          >
            {!releases.length && <option value="">Loading releases...</option>}
            {releases.map((r) => (
              <option key={r.tagName} value={r.tagName}>
                {r.tagName}
                {channels.latest?.tagName === r.tagName ? " · Latest" : ""}
                {!r.prerelease && channels.stable?.tagName === r.tagName ? " · Stable" : ""}
                {r.prerelease ? " · Beta" : ""}
              </option>
            ))}
          </select>
        </Field>

        <Field label="Among Us folder">
          <PathBox
            value={game?.path || "Not found"}
            ok={Boolean(game?.found)}
            hint={game?.message}
          />
        </Field>

        <Field label="Install location">
          <PathBox
            value={installLocation ? `${installLocation}\\${selectedMod?.folderName ?? ""}` : "Next to launcher"}
            ok
            hint="Default: same folder as the launcher executable"
          />
        </Field>

        <label className="flex cursor-pointer items-center gap-3 rounded-xl border border-line bg-panel-2 px-3 py-3 text-sm transition hover:border-white/15">
          <input
            type="checkbox"
            className="accent-purple"
            checked={createShortcut}
            disabled={busy}
            onChange={(e) => onToggleShortcut(e.target.checked)}
          />
          <span>Create desktop shortcut</span>
        </label>

        <label className="flex cursor-pointer items-center gap-3 rounded-xl border border-line bg-panel-2 px-3 py-3 text-sm transition hover:border-white/15">
          <input
            type="checkbox"
            className="accent-purple"
            checked={launchAfter}
            disabled={busy}
            onChange={(e) => onToggleLaunch(e.target.checked)}
          />
          <span>Launch after install</span>
        </label>

        <div className="rounded-xl border border-line bg-black/20 px-3 py-3">
          <div className="mb-1 text-xs uppercase tracking-wide text-muted">Status</div>
          <div className="text-sm text-white/90">
            {error ? (
              <span className="text-red">{error}</span>
            ) : busy && progress ? (
              progress.message
            ) : installed ? (
              "Ready to play"
            ) : (
              game?.message ?? "Checking..."
            )}
          </div>
          {(busy || (progress && !progress.done)) && progress && (
            <div className="mt-3">
              <ProgressBar percent={progress.percent} message={progress.message} stage={progress.stage} />
            </div>
          )}
        </div>
      </div>

      <div className="mt-4 flex gap-2">
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
          <ActionButton onClick={onInstall} disabled={!canInstall && !installed} variant="secondary">
            Reinstall
          </ActionButton>
        )}
      </div>
    </section>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="block space-y-1.5">
      <span className="text-xs font-medium uppercase tracking-wide text-muted">{label}</span>
      {children}
    </label>
  );
}

function PathBox({ value, ok, hint }: { value: string; ok?: boolean; hint?: string }) {
  return (
    <div className="rounded-xl border border-line bg-panel-2 px-3 py-2.5">
      <div className="truncate text-sm text-white/90" title={value}>
        {value}
      </div>
      {hint && (
        <div className={`mt-1 text-[11px] ${ok ? "text-muted" : "text-red"}`}>{hint}</div>
      )}
    </div>
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
