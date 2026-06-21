# Pre-Tag Checklist

## User Request

Continue toward a public ARD release with real verification and small, pushed
milestones.

## Changes

- Added `docs/releases/PRE_TAG_CHECKLIST.md` as the public-tag readiness gate.
- Linked the checklist from deployment and collaboration docs.
- Updated the quality score and release notes to make the checklist the next release
  decision point.

## Design Intent

Creating a `v*` tag is externally visible, so the repository now separates local release
rehearsal from the explicit human release decision. The checklist keeps the required
version, changelog, dry-run, live E2E, CI, public surface, and artifact verification
steps in one versioned place.

## Important Files

- `docs/releases/PRE_TAG_CHECKLIST.md`
- `docs/DEPLOYMENT.md`
- `docs/REPO_COLLAB_GUIDE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
