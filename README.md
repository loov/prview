# prview

Github PR analysis and viewing.

```
# To list conflicting PR-s by directory
prview conflicts

# To list conflicting PR-s by file path
prview -group file conflicts

# To list changes by directory
prview changes

# To list changes to directory example
prview changes example

# To list conflicts of PR #1204
prview conflicts-with 1204
```

## Setting up Github Token

GitHub API requires authentication. You can create one in [Github settings](https://github.com/settings/tokens).

Then either set the environment variable `GITHUB_TOKEN` or use `-token` command line argument.