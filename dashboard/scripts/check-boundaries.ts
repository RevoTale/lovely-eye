import { readdirSync, readFileSync } from 'node:fs';
import { extname, join, relative, sep } from 'node:path';

const sourceRoot = join(import.meta.dirname, '..', 'src');
const sourceExtensions = new Set(['.ts', '.tsx']);
const importPattern = /(?:from\s*|import\s*)['"]([^'"]+)['"]/gv;
const violations: string[] = [];

function sourceFiles(directory: string): string[] {
  return readdirSync(directory, { withFileTypes: true }).flatMap((entry) => {
    const path = join(directory, entry.name);
    if (entry.isDirectory()) {
      return sourceFiles(path);
    }
    return sourceExtensions.has(extname(entry.name)) ? [path] : [];
  });
}

function report(file: string, specifier: string, reason: string): void {
  violations.push(`${relative(sourceRoot, file)} imports ${specifier}: ${reason}`);
}

const allSourceFiles = sourceFiles(sourceRoot);
const forbiddenLegacyRoots = /^(?:components|config|gql|hooks|layouts|lib|pages|routes)\//u;

for (const file of allSourceFiles) {
  const normalizedFile = relative(sourceRoot, file).split(sep).join('/');
  const source = readFileSync(file, 'utf8');

  if (forbiddenLegacyRoots.test(normalizedFile)) {
    violations.push(`${normalizedFile}: legacy top-level ownership folders are not allowed`);
  }

  for (const match of source.matchAll(importPattern)) {
    const specifier = match[1];
    if (specifier === undefined) {
      continue;
    }

    if (
      normalizedFile.startsWith('shared/') &&
      /^@\/(?:app|features|components|hooks|layouts|pages|lib)(?:\/|$)/u.test(specifier)
    ) {
      report(file, specifier, 'shared code must not depend on app, features, or legacy modules');
    }

    if (normalizedFile.startsWith('features/') && /^@\/app(?:\/|$)/u.test(specifier)) {
      report(file, specifier, 'features must not depend on app composition or route definitions');
    }

    if (
      (specifier === 'radix-ui' ||
        specifier === 'react-day-picker' ||
        specifier.startsWith('@radix-ui/')) &&
      !normalizedFile.startsWith('shared/ui/')
    ) {
      report(file, specifier, 'UI primitives are owned by shared/ui');
    }

    const featureMatch = normalizedFile.match(/^features\/([^/]+)\//u);
    const importedFeatureMatch = specifier.match(/^@\/features\/([^/]+)(?:\/|$)/u);
    if (
      featureMatch?.[1] !== undefined &&
      importedFeatureMatch?.[1] !== undefined &&
      featureMatch[1] !== importedFeatureMatch[1] &&
      specifier.slice(`@/features/${importedFeatureMatch[1]}/`.length).includes('/')
    ) {
      report(file, specifier, "features must not import another feature's internals");
    }
  }
}

const forbiddenSharedBarrels = [
  'shared/config/index.ts',
  'shared/lib/index.ts',
  'shared/ui/index.ts',
];
for (const barrel of forbiddenSharedBarrels) {
  if (allSourceFiles.some((file) => relative(sourceRoot, file).split(sep).join('/') === barrel)) {
    violations.push(`${barrel}: compatibility barrels are not allowed`);
  }
}

if (violations.length > 0) {
  throw new Error(`Architecture boundary violations:\n${violations.join('\n')}`);
}

console.log('Frontend architecture boundaries are valid.');
