export function selectInitialSiteId(
  sites: ReadonlyArray<{ id: string }>,
  rememberedSiteId: string | null
): string | null {
  if (rememberedSiteId !== null && sites.some(({ id }) => id === rememberedSiteId)) {
    return rememberedSiteId;
  }
  return sites.length === 1 ? (sites[0]?.id ?? null) : null;
}
