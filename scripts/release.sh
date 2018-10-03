if [ -n "$CIRCLE_TAG" ]; then
  git status --porcelain || exit 1
  # mainly ignore change of package-lock.json
  git checkout -- . || exit 1
  curl -sL https://git.io/goreleaser | bash || exit 1
fi
