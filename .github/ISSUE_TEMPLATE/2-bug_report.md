---
name: Bug report
about: Create a report to help us improve
title: 'Bug Report'
labels: 'bug'
assignees: ''
body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report!
  - type: input
    id: cvwonder_version
    attributes:
      label: CVWonder Version
      description: What version of CVWonder are you using?
      placeholder: ex. 0.4.0
    validations:
      required: true
  - type: input
    id: os_version
    attributes:
      label: OS Version
      description: What operating system and version are you using?
      placeholder: ex. Windows 11, macOS 15, Ubuntu 24.04
    validations:
      required: true
  - type: textarea
    id: context
    attributes:
      label: Context
      description: What were you doing when the bug occurred?
      placeholder: Provide any relevant context...
    validations:
      required: false
  - type: textarea
    id: bug_description
    attributes:
      label: Bug Description
      description: Please provide as much detail as possible about the bug.
      placeholder: Describe the bug you encountered...
    validations:
      required: true
  - type: textarea
    id: steps_to_reproduce
    attributes:
      label: Steps to Reproduce
      description: Please list the steps to reproduce the bug.
      placeholder: 1. Step one\n2. Step two\n3. Step three...
    validations:
      required: true
  - type: textarea
    id: expected_behavior
    attributes:
      label: Expected Behavior
      description: What did you expect to happen?
      placeholder: Describe what you expected to happen...
    validations:
      required: false
  - type: checkboxes
    id: terms
    attributes:
      label: Terms and Conditions
      description: Please confirm that you agree to the following terms:
      options:
        - label: I have read and agree to the project's [Code of Conduct](https://github.com/germainlefebvre4/cvwonder/blob/main/CODE_OF_CONDUCT.md).
          value: agree_code_of_conduct
        - label: I have read and agree to the project's [Contributing Guidelines](https://github.com/germainlefebvre4/cvwonder/blob/main/CONTRIBUTING.md).
          value: agree_contributing_guidelines
    validations:
        required: true
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Desktop (please complete the following information):**
 - OS: [e.g. iOS]
 - Browser [e.g. chrome, safari]
 - Version [e.g. 22]

**Smartphone (please complete the following information):**
 - Device: [e.g. iPhone6]
 - OS: [e.g. iOS8.1]
 - Browser [e.g. stock browser, safari]
 - Version [e.g. 22]

**Additional context**
Add any other context about the problem here.
