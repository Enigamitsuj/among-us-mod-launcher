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
  const [installLocation, setInstallLocation] = useState("");
  const [createShortcut, setCreateShortcut] = useState(true);
  const [launchAfter, setLaunchAfter] = useState(true);
  const [installState, setInstallState] = useState<InstallState | null>(null);
  const [progress, setProgress] = useState<InstallProgress | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [confirmReinstall, setConfirmReinstall] = useState(false);

  const selectedModRef = useRef<Mod | null>(null);
  const installLocationRef = useRef("");
  selectedModRef.current = selectedMod;
  installLocationRef.current = installLocation;

  const refreshInstallState = useCallback(async (modId: string, location: string) => {
    try {
      const state = await Launcher.GetInstallState(modId, location);
      setInstallState(state);
    } catch {
      setInstallState(null);
    }
  }, []);

  const loadReleases = useCallback(async (mod: Mod) => {
    if (!mod.enabled) {
      setReleases([]);
      setSelectedTag("");
      return;
    }
    try {
      const list = (await Launcher.GetReleases(mod.id)) ?? [];
      setReleases(list);
      setSelectedTag(list[0]?.tagName ?? "");
      setError(null);
    } catch {
      setReleases([]);
      setSelectedTag("");
      setError("GitHub unavailable.");
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const [allMods, defaultMod, location, gameStatus] = await Promise.all([
          Launcher.GetMods(),
          Launcher.GetDefaultMod(),
          Launcher.GetDefaultInstallLocation(),
          Launcher.DetectGame(),
        ]);
        if (cancelled) return;
        setMods(allMods ?? []);
        setSelectedMod(defaultMod);
        setInstallLocation(location);
        setGame(gameStatus);
        if (defaultMod?.id) {
          await loadReleases(defaultMod);
          await refreshInstallState(defaultMod.id, location);
        }
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
          void refreshInstallState(mod.id, installLocationRef.current);
        }
      }
    });

    return () => {
      cancelled = true;
      if (typeof off === "function") off();
    };
  }, [loadReleases, refreshInstallState]);

  const onSelectMod = async (mod: Mod) => {
    setSelectedMod(mod);
    setError(null);
    setProgress(null);
    await loadReleases(mod);
    await refreshInstallState(mod.id, installLocation);
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
        installLocation,
        createShortcut,
        launchAfterInstall: launchAfter,
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
      const exists = await Launcher.DestinationExists(selectedMod.id, installLocation);
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
    } catch {
      setError("Could not launch the game.");
    }
  };

  return (
    <div className="relative flex h-full flex-col bg-ink text-text">
      <TitleBar />
      <main className="grid flex-1 grid-cols-[1.05fr_0.95fr] gap-4 p-4">
        <HeroPanel mod={selectedMod} />
        <InstallPanel
          mods={mods}
          selectedMod={selectedMod}
          onSelectMod={onSelectMod}
          game={game}
          releases={releases}
          selectedTag={selectedTag}
          onSelectTag={setSelectedTag}
          installLocation={installLocation}
          createShortcut={createShortcut}
          launchAfter={launchAfter}
          onToggleShortcut={setCreateShortcut}
          onToggleLaunch={setLaunchAfter}
          installState={installState}
          progress={progress}
          busy={busy}
          error={error}
          onInstall={onInstall}
          onPlay={onPlay}
        />
      </main>

      <ConfirmDialog
        open={confirmReinstall}
        title="Town of Us already exists."
        message="Do you want to reinstall? This will replace the current modded copy. Your original Among Us install will not be touched."
        confirmLabel="Yes"
        cancelLabel="No"
        onCancel={() => setConfirmReinstall(false)}
        onConfirm={() => void startInstall(true)}
      />
    </div>
  );
}

export default App;
