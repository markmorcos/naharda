# Re-publish the official Umami image under our namespace (add-analytics).
#
# Umami is a prebuilt third-party image, but the infra build pipeline builds
# every ghcr.io/* service from a Dockerfile (and pushes to our namespace). This
# one-line passthrough lets analytics flow through the same build/merge/deploy
# path as api/web — built as ghcr.io/markmorcos/naharda-analytics:naharda-latest
# — instead of failing on a missing Dockerfile.
#
# buildx --platform pulls the matching arch from Umami's multi-arch manifest.
# Pin to a digest here if you want to freeze the Umami version.
FROM ghcr.io/umami-software/umami:postgresql-latest
