# Project Agent Rules

## Graphify Before Push

- Before every `git push`, update the project Graphify index against the final
  working tree.
- Verify that `graphify-out/graph.json`, `graphify-out/GRAPH_REPORT.md`, and the
  generated visualization reflect the code being pushed.
- Version and push `graphify-out/`, this `AGENTS.md` rule, and every other
  Graphify-related project artifact to GitHub together with the code it
  describes.
