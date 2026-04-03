# Security Policy

## Reporting a Vulnerability

If you discover a security issue, report it privately before disclosing it publicly.

Please include:

- A clear description of the issue
- Reproduction steps
- Impact assessment
- Any logs, stack traces, or proof-of-concept material

## What Happens Next

After a report is received:

1. The issue is triaged and validated.
2. A fix plan is prepared.
3. A patch is developed and tested.
4. Coordinated disclosure follows after a fix is ready.

## Scope Guidance

Relevant areas include:

- Command execution paths in `cmd/`
- Config loading and validation in `internal/config/`
- Security scan gating behavior in `internal/security/`
- Local state handling in `.aeroform/`

## Disclosure Expectations

- Do not publish exploit details before a fix is available.
- Avoid accessing or exfiltrating data you do not own.
- Provide enough technical detail to reproduce and patch safely.
