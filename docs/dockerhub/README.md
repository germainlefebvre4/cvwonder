# Quick reference

* **Maintained by**:<br>
  [Germain LEFEBVRE](https://github.com/germainlefebvre4)

* **Where to get help**:<br>
  [Github Discussions](https://github.com/germainlefebvre4/cvwonder/discussions)

# Supported tags and respective Dockerfile links

* `latest`, `v0`, `v0.2`, `v0.2.0` ([Dockerfile](https://github.com/germainlefebvre4/cvwonder/blob/v0.2.0/Dockerfile)

# Quick reference (cont.)

* **Where to file issues**:<br>
  https://github.com/germainlefebvre4/cvwonder/issues⁠

* Supported architectures: ([more info⁠]())<br>
  amd64, arm64

* **Source of this description**:<br>
  [cvwonder repo's `docs/dockerhub/` directory](https://github.com/germainlefebvre4/cvwonder/tree/main/docs/dockerhub/) ([history](https://github.com/docker-library/docs/commits/master/nginx))

# What is CV Wonder?

CV Wonder is a tool that allows you to create a CV in a few minutes.
It allows you to massively generate CVs, base on a theme, for thousands of people in a few seconds without friction.
The Theme system allows you to use community themes and create your own for your purposes.

Don't waste any more time formatting your CV, let CV Wonder do it for you and just **focus** on the content.

# How to use this image

Generate your CV in HTML format:

```bash
docker run --rm -v $(pwd):/app germainlefebvre4/cvwonder generate
