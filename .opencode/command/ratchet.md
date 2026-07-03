---
description: Conformance health check — run the suite, ratchet upward, file issues for anomalies
agent: foreman
---

Conformance maintenance session, run by @arbiter (delegate to it; no code
changes this session):

1. Clean tree, pull, submodule update.
2. `make conformance` — report the numbers per lane (schema / instance /
   xpath) compared to the committed expectations.
3. Any case doing better → `make ratchet`, verify the diff is upward-only
   and explainable, commit as `conformance: ratchet <date> (<before> -> <after>)`.
4. Any case doing worse → do NOT touch expectations; bisect to the
   offending commit if quick (`git log --oneline` since last ratchet
   commit), and have @cartographer file a `kind/bug` issue with the case
   IDs, error output, and suspected commit.
5. Any unexplainable upward flip → also an issue (see arbiter rules).
6. @chronicler logs the session. Push.

$ARGUMENTS
