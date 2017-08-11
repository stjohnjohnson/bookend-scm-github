# bookend-scm-github

> Clones repositories and locally merges PRs for a Screwdriver build

## Usage

### Habitat

```bash
$ hab pkg exec stjohn/bookend-scm-github bookend-scm-github --help
```

### Normal
```bash
GIT_PATH=/usr/bin/git ./bookend-scm-github --host github.com --repo screwdriver-cd/screwdriver --sha 580712fb634ec01ae43246cacf186a8ecdac0d55 --target-dir /tmp/foo
Bookend:	v1.0.0
Git Client:	v2.13.3

☛ Cloning github.com/screwdriver-cd/screwdriver, on branch master
remote: Counting objects: 4605, done.
remote: Compressing objects: 100% (54/54), done.
remote: Total 4605 (delta 26), reused 24 (delta 11), pack-reused 4540
Receiving objects: 100% (4605/4605), 1.33 MiB | 1.01 MiB/s, done.
Resolving deltas: 100% (2846/2846), done.

☛ Saving local git config

☛ Resetting to 580712fb634ec01ae43246cacf186a8ecdac0d55
HEAD is now at 580712f remove delay trigger section

✓ Done
```

### Pull Requests
```bash
GIT_PATH=/usr/bin/git ./bookend-scm-github --host github.com --repo screwdriver-cd/screwdriver --sha fd3d3ac2fc765356cb230e96100293ffa33c4c98 --pull-request 692 --target-dir /tmp/foo
☛ Cloning github.com/screwdriver-cd/screwdriver, on branch master
remote: Counting objects: 4605, done.
remote: Compressing objects: 100% (54/54), done.
remote: Total 4605 (delta 26), reused 24 (delta 11), pack-reused 4540
Receiving objects: 100% (4605/4605), 1.33 MiB | 1.01 MiB/s, done.
Resolving deltas: 100% (2846/2846), done.

☛ Saving local git config

☛ Fetching PR 692
remote: Counting objects: 11, done.
remote: Compressing objects: 100% (3/3), done.
remote: Total 11 (delta 8), reused 11 (delta 8), pack-reused 0
Unpacking objects: 100% (11/11), done.
From https://github.com/screwdriver-cd/screwdriver
 * [new ref]         refs/pull/692/head -> pr

☛ Merging with master
Merge made by the 'recursive' strategy.
 package.json              | 2 +-
 plugins/auth/contexts.js  | 7 ++++---
 test/plugins/auth.test.js | 4 ++--
 3 files changed, 7 insertions(+), 6 deletions(-)

☛ Checked out 6f677d49d00217080a67409989dba37981f43e1d

✓ Done
```

## Testing

```bash
go test -cover ./...
```

## License

Code licensed under the BSD 3-Clause license. See LICENSE file for terms.
