// Shared NPC UI helpers aligned with NPCGen's core behavior.

export function textOrDash(value) {
  const normalized = String(value ?? "").trim();
  return normalized || "—";
}

export function isPresent(value) {
  return String(value ?? "").trim().length > 0;
}

export function validatePayload(payload) {
  const checks = [
    ["name", payload.name],
    ["type", payload.type],
    ["subtype", payload.subtype],
    ["species", payload.species],
    ["faction", payload.faction],
    ["trait", payload.trait],
  ];

  const missing = checks
    .filter(([, value]) => !isPresent(value))
    .map(([label]) => label);

  return {
    ok: missing.length === 0,
    missing,
  };
}
