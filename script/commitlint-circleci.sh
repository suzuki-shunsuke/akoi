if [ -n "$CIRCLE_PULL_REQUEST" ]; then
  npx commitlint --from master --to $CIRCLE_BRANCH || exit 1
else
  npx commitlint --from HEAD~3 --to HEAD || exit 1
fi
