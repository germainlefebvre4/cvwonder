---
sidebar_position: 3
---
# Theme

---

Theme system is a way to customize the look and feel of your CV.

The theme must be located in the `themes` directory in the current working directory.

```tree
themes
└── <my-theme-name>
    ├── theme.yaml     # Theme metadata (required)
    └── index.html     # Theme template (required)
```

Themes offer these abilities:

* Easily switch between themes
* Write your own theme
* Publish yours to share them with the community
* Download themes from the community

## Branch Management

Themes installed from Git repositories support branch management. You can:
- Install themes from specific branches or tags
- Switch between branches without reinstalling
- The current branch is displayed when generating your CV

Example:
```bash
# Install from main branch (default)
cvwonder theme install github.com/user/my-theme

# Switch to develop branch
cvwonder theme install github.com/user/my-theme@develop
```

The theme directory structure remains simple (`themes/my-theme/`) with git managing the branch internally.
