import { readdirSync, readFileSync, statSync } from 'node:fs';
import { extname, join, relative, sep } from 'node:path';
import * as ts from 'typescript';

const MAX_FILE_LINES = 220;
const MAX_FUNCTION_LINES = 160;

const ROOTS = ['src', '../server/static'];
const SOURCE_EXTENSIONS = new Set(['.ts', '.tsx']);
const EXCLUDED_SEGMENTS = new Set(['dist', 'gql', 'node_modules']);
const EXCLUDED_FILES = new Set(['src/routeTree.gen.ts']);

interface SourceSizeIssue {
  file: string;
  message: string;
}

const toPosix = (path: string): string => path.split(sep).join('/');

const isExcludedPath = (path: string): boolean => {
  const normalized = toPosix(path);
  if (EXCLUDED_FILES.has(normalized)) return true;
  return normalized.split('/').some((segment) => EXCLUDED_SEGMENTS.has(segment));
};

const walk = (root: string, files: string[] = []): string[] => {
  for (const entry of readdirSync(root)) {
    const path = join(root, entry);
    const normalized = toPosix(relative(process.cwd(), path));
    if (isExcludedPath(normalized)) continue;
    const stat = statSync(path);
    if (stat.isDirectory()) {
      walk(path, files);
      continue;
    }
    if (SOURCE_EXTENSIONS.has(extname(path))) {
      files.push(path);
    }
  }
  return files;
};

const lineCount = (content: string): number => content.split(/\r?\n/u).length;

const functionName = (node: ts.Node): string => {
  if (ts.isFunctionDeclaration(node) && node.name) return node.name.text;
  if (ts.isMethodDeclaration(node) && ts.isIdentifier(node.name)) return node.name.text;
  if (ts.isArrowFunction(node) || ts.isFunctionExpression(node)) {
    const parent = node.parent;
    if (ts.isVariableDeclaration(parent) && ts.isIdentifier(parent.name)) {
      return parent.name.text;
    }
    if (ts.isPropertyAssignment(parent) && ts.isIdentifier(parent.name)) {
      return parent.name.text;
    }
  }
  return '<anonymous>';
};

const checkFunctionSpans = (file: string, content: string, issues: SourceSizeIssue[]): void => {
  const sourceFile = ts.createSourceFile(file, content, ts.ScriptTarget.Latest, true);
  const visit = (node: ts.Node): void => {
    if (
      ts.isFunctionDeclaration(node) ||
      ts.isFunctionExpression(node) ||
      ts.isMethodDeclaration(node) ||
      ts.isArrowFunction(node)
    ) {
      const start = sourceFile.getLineAndCharacterOfPosition(node.getStart(sourceFile)).line + 1;
      const end = sourceFile.getLineAndCharacterOfPosition(node.getEnd()).line + 1;
      const lines = end - start + 1;
      if (lines > MAX_FUNCTION_LINES) {
        issues.push({
          file,
          message: `${functionName(node)} spans ${lines} lines; max is ${MAX_FUNCTION_LINES}`,
        });
      }
    }
    ts.forEachChild(node, visit);
  };
  visit(sourceFile);
};

const checkFile = (path: string, issues: SourceSizeIssue[]): void => {
  const file = toPosix(relative(process.cwd(), path));
  const content = readFileSync(path, 'utf8');
  const lines = lineCount(content);
  if (lines > MAX_FILE_LINES) {
    issues.push({ file, message: `${lines} lines; max is ${MAX_FILE_LINES}` });
  }
  checkFunctionSpans(file, content, issues);
};

const roots = ROOTS.map((root) => join(process.cwd(), root));
const issues: SourceSizeIssue[] = [];

for (const root of roots) {
  walk(root).forEach((path) => checkFile(path, issues));
}

if (issues.length > 0) {
  console.error('Source size check failed:');
  for (const issue of issues) {
    console.error(`- ${issue.file}: ${issue.message}`);
  }
  process.exit(1);
}

console.log('Source size check passed.');
