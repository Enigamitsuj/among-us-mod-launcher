import { motion } from "framer-motion";

type Props = {
  percent: number;
  message: string;
  stage: string;
};

export function ProgressBar({ percent, message, stage }: Props) {
  const value = Math.max(0, Math.min(100, percent));
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between text-xs">
        <span className="font-medium capitalize text-white/90">{message || stage}</span>
        <span className="tabular-nums text-muted">{Math.round(value)}%</span>
      </div>
      <div className="h-2.5 overflow-hidden rounded-full bg-black/40 ring-1 ring-white/10">
        <motion.div
          className="h-full rounded-full bg-gradient-to-r from-purple via-fuchsia-500 to-red"
          initial={{ width: 0 }}
          animate={{ width: `${value}%` }}
          transition={{ type: "spring", stiffness: 120, damping: 20 }}
        />
      </div>
    </div>
  );
}
