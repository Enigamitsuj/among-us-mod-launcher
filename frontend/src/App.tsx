import { useCallback, useEffect, useRef, useState } from "react";
import { Events } from "@wailsio/runtime";
import * as Launcher from "../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/appservice/launcher";
import type {
  GameStatus,
  InstallProgress,
  InstallState,
  Mod,
  Release,
} from "../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";
import { TitleBar } from "./components/TitleBar";
import { HeroPanel } from "./components/HeroPanel";
import { InstallPanel } from "./components/InstallPanel";
import { ConfirmDialog } from "./components/ConfirmDialog";

function App() {
  const [mods, setMods] = useState<Mod[]>([]);
  const [selectedMod, setSelectedMod] = useState<Mod | null>(null);
  const [game, setGame] = useState<GameStatus | null>(null);
  const [releases, setReleases] = useState<Release[]>([]);
  const [selectedTag, setSelectedTag] = useState("");
  const [installState, setInstallState] = useState<InstallState | null>(null);
  const [installedModIds, setInstalledModIds] = useState<Set<string>>(new Set());
  const [progress, setProgress] = useState<InstallProgress | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [confirmReinstall, setConfirmReinstall] = useState(false);
  const [confirmUninstall, setConfirmUninstall] = useState(false);

  const selectedModRef = useRef<Mod | null>(null);
  const modsRef = useRef<Mod[]>([]);
  const gamePathRef = useRef("");
  selectedModRef.current = selectedMod;
  modsRef.current = mods;
  gamePathRef.current = game?.path ?? "";

  const refreshInstallState = useCallback(async (modId: string, gamePath: string) => {
    try {
      const state = await Launcher.GetInstallState(modId, gamePath);
      setInstallState(state);
    } catch {
      setInstallState(null);
    }
  }, []);

  const refreshInstalledMods = useCallback(async (allMods: Mod[], gamePath: string) => {
    if (!gamePath) {
      setInstalledModIds(new Set());
      return;
    }
    const states = await Promise.all(
      allMods.map(async (mod) => ({
        id: mod.id,
        state: await Launcher.GetInstallState(mod.id, gamePath),
      })),
    );
    setInstalledModIds(new Set(states.filter(({ state }) => state.installed).map(({ id }) => id)));
  }, []);

  const pickPreferredTag = (list: Release[], installedVersion?: string) => {
    if (installedVersion && list.some((r) => r.tagName === installedVersion)) {
      return installedVersion;
    }
    return list[0]?.tagName ?? "";
  };

  const loadReleases = useCallback(async (mod: Mod, gamePath: string) => {
    if (!mod.enabled) {
      setReleases([]);
      setSelectedTag("");
      return;
    }
    try {
      // Ensure detection runs first so the correct storefront ZIP is chosen.
      await Launcher.DetectGame();
      const list = (await Launcher.GetReleases(mod.id)) ?? [];
      setReleases(list);

      let installedVersion = "";
      if (gamePath) {
        try {
          const state = await Launcher.GetInstallState(mod.id, gamePath);
          installedVersion = state.version ?? "";
        } catch {
          // Fall back to latest release if install state is unavailable.
        }
      }
      setSelectedTag(pickPreferredTag(list, installedVersion));
      setError(null);
    } catch {
      setReleases([]);
      setSelectedTag("");
      setError("GitHub unavailable.");
    }
  }, []);

  useEffect(() => {
    let stopped = false;
    let checking = false;
    const checkRunning = async () => {
      if (checking) return;
      checking = true;
      try {
        const running = await Launcher.IsGameRunning();
        if (!stopped) {
          setGame((current) => (current ? { ...current, running } : current));
        }
      } catch {
        // Keep the last known state if the process check is temporarily unavailable.
      } finally {
        checking = false;
      }
    };
    const timer = window.setInterval(() => void checkRunning(), 3000);
    return () => {
      stopped = true;
      window.clearInterval(timer);
    };
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const [allMods, defaultMod, gameStatus] = await Promise.all([
          Launcher.GetMods(),
          Launcher.GetDefaultMod(),
          Launcher.DetectGame(),
        ]);
        if (cancelled) return;
        setMods(allMods ?? []);
        setSelectedMod(defaultMod);
        setGame(gameStatus);
        if (defaultMod?.id) {
          await loadReleases(defaultMod, gameStatus.path);
          await refreshInstallState(defaultMod.id, gameStatus.path);
        }
        await refreshInstalledMods(allMods ?? [], gameStatus.path);
      } catch {
        if (!cancelled) setError("Failed to initialize launcher.");
      }
    })();

    const off = Events.On("install:progress", (ev: { data?: InstallProgress }) => {
      const data = ev?.data;
      if (!data) return;
      setProgress(data);
      if (data.error) {
        setError(data.error);
        setBusy(false);
      }
      if (data.done && !data.error) {
        setBusy(false);
        setError(null);
        const mod = selectedModRef.current;
        if (mod) {
          void refreshInstallState(mod.id, gamePathRef.current);
          void refreshInstalledMods(modsRef.current, gamePathRef.current);
        }
      }
    });

    return () => {
      cancelled = true;
      if (typeof off === "function") off();
    };
  }, [loadReleases, refreshInstalledMods, refreshInstallState]);

  const onSelectMod = async (mod: Mod) => {
    setSelectedMod(mod);
    setError(null);
    setProgress(null);
    await loadReleases(mod, game?.path ?? "");
    await refreshInstallState(mod.id, game?.path ?? "");
  };

  const startInstall = async (force: boolean) => {
    if (!selectedMod) return;
    setConfirmReinstall(false);
    setBusy(true);
    setError(null);
    setProgress({
      stage: "preparing",
      message: "Preparing...",
      percent: 1,
      done: false,
      error: "",
      installDir: "",
    });

    try {
      await Launcher.Install({
        modId: selectedMod.id,
        versionTag: selectedTag,
        amongUsPath: game?.path ?? "",
        platform: game?.platform ?? "",
        forceReinstall: force,
      });
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      if (msg.toLowerCase().includes("already exists")) {
        setBusy(false);
        setConfirmReinstall(true);
        setProgress(null);
        return;
      }
      setBusy(false);
      setError(msg || "Installation failed.");
    }
  };

  const onInstall = async () => {
    if (!selectedMod) return;
    try {
      const exists = await Launcher.DestinationExists(selectedMod.id, game?.path ?? "");
      if (exists) {
        setConfirmReinstall(true);
        return;
      }
    } catch {
      // proceed
    }
    await startInstall(false);
  };

  const onPlay = async () => {
    if (!installState?.path) return;
    try {
      await Launcher.LaunchMod(installState.path);
      setGame((current) => (current ? { ...current, running: true } : current));
    } catch {
      setError("Could not launch the game.");
    }
  };

  const onUninstall = async () => {
    if (!selectedMod) return;
    setConfirmUninstall(false);
    setBusy(true);
    setError(null);
    try {
      await Launcher.Uninstall(selectedMod.id, game?.path ?? "");
      await refreshInstallState(selectedMod.id, game?.path ?? "");
      await refreshInstalledMods(mods, game?.path ?? "");
    } catch (e) {
      const message = e instanceof Error ? e.message : String(e);
      setError(message || "Could not uninstall the mod.");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="relative flex h-full min-h-0 flex-col bg-ink text-text">
      <TitleBar />
      <main className="grid min-h-0 flex-1 grid-cols-[1.05fr_0.95fr] gap-4 overflow-hidden p-4">
        <HeroPanel
          mod={selectedMod}
          release={releases.find((r) => r.tagName === selectedTag) ?? null}
        />
        <InstallPanel
          mods={mods}
          selectedMod={selectedMod}
          onSelectMod={onSelectMod}
          game={game}
          releases={releases}
          selectedTag={selectedTag}
          onSelectTag={setSelectedTag}
          installedModIds={installedModIds}
          installState={installState}
          progress={progress}
          busy={busy}
          error={error}
          onInstall={onInstall}
          onPlay={onPlay}
          onUninstall={() => setConfirmUninstall(true)}
        />
      </main>

      <ConfirmDialog
        open={confirmReinstall}
        title={
          installState?.version && installState.version !== selectedTag
            ? `Switch to ${selectedTag}?`
            : `${selectedMod?.shortName ?? "Mod"} is already installed.`
        }
        message={
          installState?.version && installState.version !== selectedTag
            ? `Replace ${installState.version} with ${selectedTag}. Your current copy is kept until the new one is ready. The original Among Us install is never touched.`
            : "This builds the selected version first, then safely replaces your current modded copy. Your original Among Us install will not be touched."
        }
        confirmLabel={
          installState?.version && installState.version !== selectedTag ? "Install" : "Reinstall"
        }
        cancelLabel="Cancel"
        onCancel={() => setConfirmReinstall(false)}
        onConfirm={() => void startInstall(true)}
      />
      <ConfirmDialog
        open={confirmUninstall}
        title={`Uninstall ${selectedMod?.shortName ?? "mod"}?`}
        message="This removes only the launcher-managed modded copy. Your original Among Us installation will not be touched."
        confirmLabel="Uninstall"
        cancelLabel="Cancel"
        onCancel={() => setConfirmUninstall(false)}
        onConfirm={() => void onUninstall()}
      />
    </div>
  );
}

export default App;
