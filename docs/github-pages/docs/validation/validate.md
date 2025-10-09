---
sidebar_position: 1
---
# Validate

---

:::info
Consider validating your CV file before generating to ensure quality and correctness.
:::

## Getting Started

### Basic Validation

Validate a CV file and see all errors and warnings:

```bash
cvwonder validate -i cv.yml
```

### Validate Before Generating

Ensure your CV is valid before generating output:

```bash
cvwonder generate -i cv.yml --validate -o output/
```

This will:
1. Validate the YAML file
2. Display any errors or warnings
3. Exit if validation fails
4. Generate the CV if validation passes

### Validate in Watch Mode

Combine validation with watch mode for live feedback:

```bash
cvwonder serve -i cv.yml --validate --watch
```

### CI/CD Integration

Use validation in your CI/CD pipeline. Run basically the same command in your pipeline while the return code will indicate success or failure:

```bash
cvwonder validate -i cv.yml
```
