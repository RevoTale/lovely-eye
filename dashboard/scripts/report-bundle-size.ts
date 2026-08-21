import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { gzipSync } from 'node:zlib';

interface ManifestEntry {
  file: string;
  imports: string[];
  css: string[];
}

interface SizeSummary {
  rawBytes: number;
  gzipBytes: number;
  files: string[];
}

interface BundleReport {
  initial: SizeSummary;
  routes: Record<keyof typeof TARGETS, SizeSummary>;
  totalJavaScript: SizeSummary;
  totalCSS: SizeSummary;
}

const TARGETS = {
  login: 'src/features/auth/screens/login-screen.tsx',
  register: 'src/features/auth/screens/register-screen.tsx',
  sites: 'src/features/sites/screens/sites-screen.tsx',
  createSite: 'src/features/sites/screens/create-site-screen.tsx',
  analytics: 'src/features/analytics/screens/analytics-screen.tsx',
  siteSettings: 'src/features/site-settings/screens/site-settings-screen.tsx',
} as const;

const outputDirectory = process.argv.slice(2).find((argument) => !argument.startsWith('--'));
if (outputDirectory === undefined) {
  throw new Error('Usage: report-bundle-size.ts <vite-output-directory> [--check]');
}

const manifest = parseManifest(
  JSON.parse(readFileSync(join(outputDirectory, '.vite', 'manifest.json'), 'utf8'))
);
const initialFiles = collectFiles(manifest, 'index.html');
initialFiles.add('index.html');

const routes = Object.fromEntries(
  Object.entries(TARGETS).map(([label, key]) => {
    const routeFiles = collectFiles(manifest, key);
    const incrementalFiles = new Set([...routeFiles].filter((file) => !initialFiles.has(file)));
    return [label, summarizeFiles(outputDirectory, incrementalFiles)];
  })
) as BundleReport['routes'];

const assetFiles = readdirSync(join(outputDirectory, 'assets')).map((file) => `assets/${file}`);
const javascriptFiles = new Set(assetFiles.filter((file) => file.endsWith('.js')));
const cssFiles = new Set(assetFiles.filter((file) => file.endsWith('.css')));

const report: BundleReport = {
  initial: summarizeFiles(outputDirectory, initialFiles),
  routes,
  totalJavaScript: summarizeFiles(outputDirectory, javascriptFiles),
  totalCSS: summarizeFiles(outputDirectory, cssFiles),
};
process.stdout.write(`${JSON.stringify(report, null, 2)}\n`);

if (process.argv.includes('--check')) assertBundleBudgets(report);

function parseManifest(value: unknown): Record<string, ManifestEntry> {
  if (typeof value !== 'object' || value === null) throw new Error('Invalid Vite manifest');
  const result: Record<string, ManifestEntry> = {};
  for (const [key, entry] of Object.entries(value)) {
    if (typeof entry !== 'object' || entry === null || !('file' in entry)) {
      throw new Error(`Invalid Vite manifest entry: ${key}`);
    }
    const file = entry.file;
    if (typeof file !== 'string') throw new Error(`Invalid Vite manifest file: ${key}`);
    result[key] = {
      file,
      imports: readStringArray(entry, 'imports'),
      css: readStringArray(entry, 'css'),
    };
  }
  return result;
}

function readStringArray(value: object, key: string): string[] {
  if (!(key in value)) return [];
  const candidate: unknown = Reflect.get(value, key);
  if (!Array.isArray(candidate) || !candidate.every((item) => typeof item === 'string')) {
    throw new Error(`Invalid Vite manifest ${key}`);
  }
  return candidate;
}

function collectFiles(
  manifest: Record<string, ManifestEntry>,
  entryKey: string,
  visited = new Set<string>()
): Set<string> {
  if (visited.has(entryKey)) return new Set();
  visited.add(entryKey);
  const entry = manifest[entryKey];
  if (entry === undefined) throw new Error(`Missing Vite manifest entry: ${entryKey}`);
  const files = new Set([entry.file, ...entry.css]);
  for (const importedKey of entry.imports) {
    for (const file of collectFiles(manifest, importedKey, visited)) files.add(file);
  }
  return files;
}

function summarizeFiles(directory: string, files: Set<string>): SizeSummary {
  let rawBytes = 0;
  let gzipBytes = 0;
  for (const file of files) {
    const content = readFileSync(join(directory, file));
    rawBytes += content.byteLength;
    gzipBytes += gzipSync(content).byteLength;
  }
  return { rawBytes, gzipBytes, files: [...files].sort() };
}

function assertBundleBudgets(report: BundleReport): void {
  const budgets = [
    ['initial dashboard', report.initial.gzipBytes, 265_000],
    ['analytics route increment', report.routes.analytics.gzipBytes, 136_000],
    ['all dashboard JavaScript', report.totalJavaScript.gzipBytes, 425_000],
    ['all dashboard CSS', report.totalCSS.gzipBytes, 12_500],
  ] as const;
  const failures = budgets.filter(([, actual, maximum]) => actual > maximum);
  if (failures.length === 0) return;
  throw new Error(
    failures
      .map(([label, actual, maximum]) => `${label}: ${actual} gzip bytes exceeds ${maximum}`)
      .join('\n')
  );
}
