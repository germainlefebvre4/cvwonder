---
sidebar_position: 2
---
# Docker

---

The Docker image is available on [DockerHub germainlefebvre4/cvwonder](https://hub.docker.com/r/germainlefebvre4/cvwonder).

You can run it using the following command:

```bash
docker run --rm -v $(pwd):/cv germainlefebvre4/cvwonder:latest generate
```

This command will mount the current directory to the `/app` directory in the container. You can then run CV Wonder commands inside the container.
