import { motion } from "framer-motion";
import type { Mod } from "../../bindings/github.com/Enigamitsuj/among-us-mod-launcher/internal/mods/models";

type Props = {
  mod: Mod | null;
};

export function HeroPanel({ mod }: Props) {
  const from = mod?.accentFrom ?? "#a855f7";
  const to = mod?.accentTo ?? "#ef4444";

  return (
    <section className="relative flex h-full flex-col overflow-hidden rounded-2xl border border-white/10 bg-panel">
      <div
        className="absolute inset-0 opacity-90"
        style={{
          background: `
            radial-gradient(ellipse at 20% 20%, ${from}55, transparent 45%),
            radial-gradient(ellipse at 80% 30%, ${to}40, transparent 40%),
            radial-gradient(ellipse at 50% 100%, #1e1b4b 0%, #0a0a10 70%),
            linear-gradient(160deg, #0b1020, #160b1d 60%, #0a0a10)
          `,
        }}
      />

      {/* Stylized space / Mira HQ silhouette */}
      <div className="pointer-events-none absolute inset-0">
        <div className="absolute left-[12%] top-[18%] h-2 w-2 rounded-full bg-white/80 shadow-[0_0_12px_#fff]" />
        <div className="absolute left-[70%] top-[22%] h-1.5 w-1.5 rounded-full bg-cyan-200/80" />
        <div className="absolute left-[40%] top-[12%] h-1 w-1 rounded-full bg-purple/80" />
        <div className="absolute bottom-0 left-0 right-0 h-40 bg-gradient-to-t from-black/70 to-transparent" />
        <svg className="absolute bottom-16 left-1/2 w-[120%] -translate-x-1/2 opacity-40" viewBox="0 0 800 180" fill="none">
          <path d="M0 140 L80 120 L140 130 L220 90 L300 110 L380 70 L460 100 L560 60 L640 95 L720 80 L800 110 L800 180 L0 180 Z" fill="#1f2937" />
          <rect x="360" y="40" width="80" height="70" rx="6" fill="#334155" />
          <rect x="380" y="55" width="18" height="18" rx="2" fill="#67e8f9" opacity="0.7" />
          <rect x="410" y="55" width="18" height="18" rx="2" fill="#c084fc" opacity="0.7" />
        </svg>
      </div>

      {/* Crewmates row */}
      <div className="pointer-events-none absolute bottom-24 left-8 flex items-end gap-3">
        <Crewmate color="#38bdf8" />
        <Crewmate color="#a3e635" small />
        <Crewmate color="#f472b6" />
        <Impostor />
      </div>

      <div className="relative z-10 flex flex-1 flex-col justify-between p-7">
        <div>
          <motion.p
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            className="mb-3 text-xs font-semibold uppercase tracking-[0.22em] text-purple/90"
          >
            Community Installer
          </motion.p>
          <motion.h1
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.05 }}
            className="max-w-md text-4xl font-bold leading-tight tracking-tight text-white"
          >
            {mod?.name ?? "Among Us Mod Launcher"}
          </motion.h1>
          <motion.p
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="mt-3 max-w-sm text-sm leading-relaxed text-muted"
          >
            MIRA HQ · neon nights · zero config. Launch, install, play.
          </motion.p>
        </div>

        <div className="space-y-1">
          <p className="text-xs text-white/70">Installer created by FBI OpenUp</p>
          <p className="text-[11px] text-muted">Community-made launcher for Town of Us: Mira and friends.</p>
        </div>
      </div>
    </section>
  );
}

function Crewmate({ color, small }: { color: string; small?: boolean }) {
  const size = small ? 34 : 42;
  return (
    <div
      className="relative rounded-[40%] border border-white/10"
      style={{
        width: size,
        height: size * 1.25,
        background: `linear-gradient(180deg, ${color}, ${color}cc)`,
        boxShadow: `0 0 18px ${color}55`,
      }}
    >
      <div className="absolute left-[28%] top-[28%] h-[28%] w-[48%] rounded-full bg-sky-200/90" />
    </div>
  );
}

function Impostor() {
  return (
    <div className="relative ml-2 flex h-[56px] w-[46px] items-center justify-center">
      <div
        className="absolute inset-0 rounded-[40%]"
        style={{
          background: "linear-gradient(180deg, #ef4444, #991b1b)",
          boxShadow: "0 0 22px rgba(239,68,68,0.55)",
        }}
      />
      <div className="absolute left-[26%] top-[26%] h-[26%] w-[50%] rounded-full bg-slate-900/80" />
      <div className="absolute -right-1 top-3 rounded bg-black/70 px-1.5 py-0.5 text-[9px] font-bold tracking-wide text-white">
        SHHH
      </div>
    </div>
  );
}
