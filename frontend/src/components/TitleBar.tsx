import * as Launcher from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/appservice/launcher";

export function TitleBar() {
  return (
    <div className="drag-region flex h-10 items-center justify-between border-b border-line/80 bg-ink/90 px-3">
      <div className="flex items-center gap-2 text-xs font-medium tracking-wide text-muted">
        <span className="inline-flex h-5 w-5 items-center justify-center rounded-md bg-gradient-to-br from-purple to-red text-[10px] font-bold text-white">
          AU
        </span>
        <span>Among Us Mod Launcher</span>
      </div>
      <div className="no-drag flex items-center gap-1">
        <button
          type="button"
          aria-label="Minimize"
          className="flex h-7 w-9 items-center justify-center rounded-md text-muted transition hover:bg-white/10 hover:text-text"
          onClick={() => void Launcher.MinimizeWindow()}
        >
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M2 6h8" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
          </svg>
        </button>
        <button
          type="button"
          aria-label="Close"
          className="flex h-7 w-9 items-center justify-center rounded-md text-muted transition hover:bg-red/90 hover:text-white"
          onClick={() => void Launcher.CloseWindow()}
        >
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M3 3l6 6M9 3L3 9" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
          </svg>
        </button>
      </div>
    </div>
  );
}
