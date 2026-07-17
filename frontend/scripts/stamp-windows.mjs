import path from "node:path";
import { rcedit } from "rcedit";

const [exeArg, versionArg = "1.0.0", iconArg] = process.argv.slice(2);

if (!exeArg) {
  throw new Error("Usage: node stamp-windows.mjs <exe> [version] [icon]");
}

const exePath = path.resolve(exeArg);
const version = versionArg.replace(/^v/i, "");
const numericVersion = version.split("-")[0];

await rcedit(exePath, {
  "file-version": numericVersion,
  "product-version": numericVersion,
  "version-string": {
    CompanyName: "Enigamitsuj",
    FileDescription: "Community mod hub and launcher for Among Us",
    FileVersion: version,
    LegalCopyright: "© 2026 Enigamitsuj",
    OriginalFilename: "among-us-mod-launcher.exe",
    ProductName: "Among Us Mod Launcher",
    ProductVersion: version,
  },
  ...(iconArg ? { icon: path.resolve(iconArg) } : {}),
});

console.log(`Stamped ${exePath} as ${version}`);
